package mapo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo/abi"
	"math/big"
)

func (b *Bridge) GetObservationsStdTx(txIn *types.TxIn) ([]byte, error) {
	//  check
	if len(txIn.TxArray) == 0 {
		return nil, nil
	}
	// Here we construct tx according to method， and return hex tx2bytes
	var (
		err   error
		input []byte
	)
	ele := txIn.TxArray[0]
	switch ele.Method {
	case constants.VoteTxIn:
		input, err = b.mainAbi.Pack(constants.VoteTxIn, &abi.VoteTxIn{
			TxInType:  ele.TxInType,
			ToChain:   big.NewInt(0),
			Height:    ele.Height,
			FromChain: ele.FromChain,
			Amount:    ele.Amount,
			OrderId:   ele.OrderId,
			Vault:     ele.Vault,
			Token:     ele.Token,
			From:      nil,
			To:        ele.To,
			Payload:   nil,
		})

	case constants.VoteTxOut:
	}

	if err != nil {
		return nil, fmt.Errorf("fail to method(%s) pack input: %w",
			txIn.TxArray[0].Method, err)
	}

	return b.assemblyInternalTx(context.Background(), input, 0)
}

func (b *Bridge) assemblyInternalTx(ctx context.Context, input []byte, recommendLimit uint64) ([]byte, error) {
	// estimate gas
	gasFeeCap := b.gasPrice
	fromAddr, _ := b.keys.GetEthAddress()
	to := ecommon.HexToAddress(b.cfg.Maintainer)
	gasLimit, err := b.ethClient.EstimateGas(ctx, ethereum.CallMsg{
		From:     fromAddr,
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to estimate gas limit: %w", err)
	}
	if gasFeeCap.Cmp(big.NewInt(0)) == 0 {
		head, err := b.ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch head number: %w", err)
		}
		gasFeeCap = head.BaseFee
	}

	// assemble tx
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
