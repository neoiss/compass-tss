package mapo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/structure"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/pkg/errors"
)

// GetKeygenBlock retrieves keygen request for the given block height from mapBridge
func (b *Bridge) GetKeygenBlock() (*structure.KeyGen, error) {
	method := "getElectionEpoch"
	input, err := b.mainAbi.Pack(method)
	if err != nil {
		return nil, errors.Wrap(err, "fail to pack input")
	}

	to := ecommon.HexToAddress(b.cfg.Maintainer)
	outPut, err := b.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: constants.ZeroAddress,
		To:   &to,
		Data: input,
	}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract")
	}

	outputs := b.mainAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(outPut)
	if err != nil {
		return nil, errors.Wrap(err, "unpack output")
	}

	var epoch *big.Int
	if err = outputs.Copy(&epoch, unpack); err != nil {
		return nil, errors.Wrap(err, "copy output")
	}
	b.logger.Info().Int64("epoch", epoch.Int64()).Msg("KeyGen Block")
	if epoch.Uint64() == 0 { // not in epoch
		return nil, nil
	}
	// done
	ret, err := b.GetNodeAccounts(epoch)
	if err != nil {
		return nil, err
	}

	return &structure.KeyGen{
		Epoch: epoch,
		Ms:    ret,
	}, nil
}

// SendKeyGenStdTx get keygen tx from params
func (b *Bridge) SendKeyGenStdTx(epoch *big.Int, poolPubKey common.PubKey, signature []byte, blames []ecommon.Address,
	members []ecommon.Address) (string, error) {
	idAbi, _ := newIdABi()
	id, err := idAbi.Methods["idPack"].Inputs.Pack(poolPubKey, members, epoch, blames)
	if err != nil {
		return "", errors.Wrap(err, "id pack input failed")
	}

	id32 := ecommon.BytesToHash(crypto.Keccak256(id))
	method := "voteUpdateTssPool"
	input, err := b.mainAbi.Pack(method, &structure.TssPoolParam{
		Id:        id32,
		Epoch:     epoch,
		Pubkey:    ecommon.Hex2Bytes(poolPubKey.String()),
		Members:   members,
		Blames:    blames,
		Signature: signature,
	})
	if err != nil {
		return "", errors.Wrap(err, "fail to pack input")
	}

	fromAddr := "0x2b7588165556aB2fA1d30c520491C385BAa424d8"
	nonce, err := b.ethRpc.GetNonce(fromAddr)
	if err != nil {
		return "", fmt.Errorf("fail to fetch account(%s) nonce : %w", fromAddr, err)
	}

	// abort signing if the pending nonce is too far in the future
	var finalizedNonce uint64
	finalizedNonce, err = b.ethRpc.GetNonceFinalized(fromAddr)
	if err != nil {
		return "", fmt.Errorf("fail to fetch account(%s) finalized nonce: %w", fromAddr, err)
	}
	// todo handler add cfg
	if (nonce - finalizedNonce) > 3 {
		b.logger.Warn().
			Uint64("nonce", nonce).
			Uint64("finalizedNonce", finalizedNonce).
			Msg("pending nonce too far in future")
		return "", fmt.Errorf("pending nonce too far in future")
	}

	gasFeeCap := b.gasPrice
	to := ecommon.HexToAddress(b.cfg.Maintainer)
	createdTx := ethereum.CallMsg{
		From:     ecommon.HexToAddress(fromAddr),
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	}

	gasLimit, err := b.ethClient.EstimateGas(context.Background(), createdTx)
	if err != nil {
		b.logger.Err(err).Msgf("fail to estimate gas")
		return "", nil
	}

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
	txID, err := b.Broadcast(&stypes.TxOutItem{}, sign)
	if err != nil {
		return "", err
	}
	return txID, nil
}
