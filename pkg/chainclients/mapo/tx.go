package mapo

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/structure"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/pkg/errors"
)

var ErrOfOrderExist = errors.New("order already exist")

func (b *Bridge) GetObservationsStdTx(txIn *types.TxInItem) ([]byte, error) {
	//  check
	if txIn == nil {
		return nil, nil
	}
	// Here we construct tx according to method， and return tx hex bytes
	var (
		err   error
		input []byte
		ele   = txIn
	)

	switch ele.Method {
	case constants.VoteTxIn:
		input, err = b.tssAbi.Pack(constants.VoteTxIn, &structure.TxInItem{
			Height:     ele.Height,
			OrderId:    ele.OrderId,
			RefundAddr: ele.RefundAddr.Bytes(),
			BridgeItem: structure.BridgeItem{
				ChainAndGasLimit: ele.ChainAndGasLimit,
				Vault:            ele.Vault,
				TxType:           ele.TxOutType,
				Sequence:         ele.Sequence,
				Token:            ele.Token,
				Amount:           ele.Amount,
				From:             ele.From,
				To:               ele.To,
				Payload:          ele.Payload,
			},
		})
	case constants.VoteTxOut:
		input, err = b.tssAbi.Pack(constants.VoteTxOut, &structure.VoteTxOut{
			Height:  ele.Height,
			GasUsed: ele.GasUsed,
			OrderId: ele.OrderId,
			Sender:  ecommon.HexToAddress(ele.Sender),
			BridgeItem: structure.BridgeItem{
				ChainAndGasLimit: ele.ChainAndGasLimit,
				Vault:            ele.Vault,
				TxType:           ele.TxOutType,
				Sequence:         ele.Sequence,
				Token:            ele.Token,
				Amount:           ele.Amount,
				From:             ele.From,
				To:               ele.To,
				Payload:          ele.Payload,
			},
		})
	}

	if err != nil {
		return nil, fmt.Errorf("fail to method(%s) pack input: %w", txIn.Method, err)
	}

	return b.assemblyTx(context.Background(), input, 0, b.cfg.TssManager)
}

func (b *Bridge) GetOracleStdTx(txOut *types.TxOutItem) ([]byte, error) {
	//  check
	if txOut == nil {
		return nil, nil
	}
	// Here we construct tx according to method， and return tx hex bytes
	var (
		err   error
		input []byte
	)

	// todo will next 2
	input, err = b.mainAbi.Pack(constants.VoteTxOut, &structure.VoteTxOut{})
	if err != nil {
		return nil, fmt.Errorf("fail to pack oracleMethod: %w", err)
	}

	return b.assemblyTx(context.Background(), input, 0, b.cfg.TssManager) // todo this add is error
}

func (b *Bridge) assemblyTx(ctx context.Context, input []byte, recommendLimit uint64, addr string) ([]byte, error) {
	// estimate gas
	gasFeeCap := b.gasPrice
	fromAddr, _ := b.keys.GetEthAddress()
	to := ecommon.HexToAddress(addr)
	gasLimit, err := b.ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:     fromAddr,
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	})
	if err != nil {
		b.logger.Error().Any("err", err).Str("input", ecommon.Bytes2Hex(input)).Msg("Estimate failed")
		return nil, fmt.Errorf("fail to estimate gas limit: %w", err)
	}
	if gasFeeCap.Cmp(big.NewInt(0)) == 0 {
		head, err := b.ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch head number: %w", err)
		}
		gasFeeCap = head.BaseFee
	}

	nonce, err := b.ethRpc.GetNonce(fromAddr.Hex())
	if err != nil {
		return nil, fmt.Errorf("fail to fetch account(%s) nonce : %w", fromAddr, err)
	}

	var finalizedNonce uint64
	finalizedNonce, err = b.ethRpc.GetNonceFinalized(fromAddr.Hex())
	if err != nil {
		return nil, fmt.Errorf("fail to fetch account(%s) finalized nonce: %w", fromAddr, err)
	}
	// todo handler add cfg
	if (nonce - finalizedNonce) > 3 {
		b.logger.Warn().Uint64("nonce", nonce).Uint64("finalizedNonce", finalizedNonce).
			Msg("pending nonce too far in future")
		return nil, fmt.Errorf("pending nonce too far in future")
	}

	tipCap := new(big.Int).Mul(gasFeeCap, big.NewInt(10)) // todo add cfg
	tipCap.Div(tipCap, big.NewInt(100))
	if recommendLimit != 0 {
		gasLimit = recommendLimit
	}
	td := etypes.NewTx(&etypes.DynamicFeeTx{
		Nonce:     nonce,
		Value:     nil,
		To:        &to,
		Gas:       gasLimit,
		GasTipCap: tipCap,
		GasFeeCap: gasFeeCap,
		Data:      input,
	})

	ret, err := b.kw.LocalSign(td)
	if err != nil {
		return nil, fmt.Errorf("fail to sign transaction: %w", err)
	}

	return ret, nil
}

// Broadcast Broadcasts tx to mapBridge
func (b *Bridge) Broadcast(hexTx []byte) (string, error) {
	// done
	b.broadcastLock.Lock()
	defer b.broadcastLock.Unlock()

	// decode the transaction
	tx := &etypes.Transaction{}
	if err := tx.UnmarshalJSON(hexTx); err != nil {
		return "", err
	}
	txID := tx.Hash().String()

	// get context with default timeout
	ctx, cancel := b.getTimeoutContext()
	defer cancel()

	// send the transaction
	if err := b.ethClient.SendTransaction(ctx, tx); !isAcceptableError(err) {
		b.logger.Error().Str("txId", txID).Err(err).Msg("Failed to send transaction")
		return "", err
	}
	b.logger.Debug().Str("txId", txID).Msg("Broadcast tx")
	return txID, nil
}

func (b *Bridge) getTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}
