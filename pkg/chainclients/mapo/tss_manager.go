package mapo

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/structure"
	"github.com/pkg/errors"
)

// SendKeyGenStdTx get keygen tx from params
func (b *Bridge) SendKeyGenStdTx(epoch *big.Int, poolPubKey common.PubKey, signature, keyShares []byte, blames []ecommon.Address,
	members []ecommon.Address) (string, error) {
	var txID string
	if epoch.Cmp(b.epoch) == 0 {
		b.logger.Info().Any("epoch", epoch).Msg("The epoch is the same as the last one, skip sending keygen tx")
		return "", nil
	}

	pubBytes := make([]byte, 0)
	if !poolPubKey.IsEmpty() {
		ethPubKey, err := crypto.DecompressPubkey(ecommon.Hex2Bytes(poolPubKey.String()))
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal ECDSA public key: %w", err)
		}
		pubBytes = crypto.FromECDSAPub(ethPubKey)
		pubBytes = pubBytes[1:] // remove 0x04 prefix
	}

	method := constants.VoteUpdateTssPool
	input, err := b.tssAbi.Pack(method, &structure.TssPoolParam{
		Epoch:     epoch,
		Pubkey:    pubBytes,
		KeyShare:  keyShares,
		Members:   members,
		Blames:    blames,
		Signature: signature,
	})
	if err != nil {
		return "", errors.Wrap(err, "fail to pack input")
	}

	fromAddr, _ := b.keys.GetEthAddress()
	nonce, err := b.ethRpc.GetNonce(fromAddr.Hex())
	if err != nil {
		return "", fmt.Errorf("fail to fetch account(%s) nonce : %w", fromAddr, err)
	}

	// abort signing if the pending nonce is too far in the future
	var finalizedNonce uint64
	finalizedNonce, err = b.ethRpc.GetNonceFinalized(fromAddr.Hex())
	if err != nil {
		return "", fmt.Errorf("fail to fetch account(%s) finalized nonce: %w", fromAddr, err)
	}
	// todo handler add cfg
	if (nonce - finalizedNonce) > 3 {
		b.logger.Warn().Uint64("nonce", nonce).Uint64("finalizedNonce", finalizedNonce).
			Msg("pending nonce too far in future")
		return "", fmt.Errorf("pending nonce too far in future")
	}

	fmt.Println(" input ----------- ", ecommon.Bytes2Hex(input))
	gasFeeCap := b.gasPrice
	to := ecommon.HexToAddress(b.cfg.TssManager)
	createdTx := ethereum.CallMsg{
		From:     fromAddr,
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	}

	gasLimit, err := b.ethClient.EstimateGas(context.Background(), createdTx)
	if err != nil {
		b.logger.Info().Any("err", err).Msgf("Fail to estimate gas")
		return "", err
	}

	if gasFeeCap.Cmp(big.NewInt(0)) == 0 {
		head, err := b.ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return "", err
		}
		gasFeeCap = head.BaseFee
	}
	gasLimit = gasLimit*2 + 2000000
	// tip cap at configured percentage of max fee
	tipCap := new(big.Int).Mul(gasFeeCap, big.NewInt(10))
	tipCap.Div(tipCap, big.NewInt(100))
	td := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Value:     nil,
		To:        &to,
		Gas:       gasLimit,
		GasTipCap: tipCap,
		GasFeeCap: gasFeeCap,
		Data:      input,
	})

	sign, err := b.kw.LocalSign(td)
	if err != nil {
		return "", err
	}
	txID, err = b.Broadcast(sign)
	if err != nil {
		return "", err
	}
	// todo handler tx online
	b.epoch = epoch
	b.logger.Info().Msg("SendKeyGenStdTx txID is " + txID)
	time.Sleep(time.Second * 10)
	return txID, nil
}

type TSSManagerKeyShare struct {
	Pubkey   []byte
	KeyShare []byte
}

func (b *Bridge) GetKeyShare() ([]byte, []byte, error) {
	method := constants.GetKeyShare
	signerAddr, _ := b.keys.GetEthAddress()
	input, err := b.tssAbi.Pack(method, signerAddr)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to pack input")
	}
	type retStruct struct {
		KeyShare TSSManagerKeyShare `json:"keyShare"`
	}
	var ret retStruct
	err = b.callContract(&ret, b.cfg.TssManager, method, input, b.tssAbi)
	if err != nil {
		return nil, nil, err
	}
	return ret.KeyShare.KeyShare, ret.KeyShare.Pubkey, nil
}

func (b *Bridge) getTssStatus(epoch *big.Int) (constants.TssStatus, error) {
	var ret uint8
	method := constants.GetTSSStatus
	input, err := b.tssAbi.Pack(method, epoch)
	if err != nil {
		return constants.TssStatusUnknown, errors.Wrap(err, "fail to pack input")
	}
	err = b.callContract(&ret, b.cfg.TssManager, method, input, b.tssAbi)
	if err != nil {
		return constants.TssStatusUnknown, errors.Wrap(err, "fail to call contract")
	}
	return constants.TssStatus(ret), nil
}

// PostNetworkFee send network fee message to MAP
func (b *Bridge) PostNetworkFee(ctx context.Context, height int64, chainId *big.Int, transactionSize,
	transactionSizeWithCall, transactionRate uint64) (string, error) {
	input, err := b.tssAbi.Pack(constants.VoteNetworkFee,
		chainId,
		big.NewInt(height),
		big.NewInt(0).SetUint64(transactionRate),
		big.NewInt(0).SetUint64(transactionSize),
		big.NewInt(0).SetUint64(transactionSizeWithCall))
	if err != nil {
		return "", fmt.Errorf("fail to pack input: %w", err)
	}

	tx, err := b.assemblyTx(ctx, input, 2000000, b.cfg.TssManager)
	if err != nil {
		return "", fmt.Errorf("fail to assembly tx: %w", err)
	}

	return b.Broadcast(tx)
}
