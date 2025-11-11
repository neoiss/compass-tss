package mapo

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

func (b *Bridge) GetChainID(name string) (*big.Int, error) {
	method := "getChainByName"
	input, err := b.tokenRegistry.Pack(method, name)
	if err != nil {
		return nil, err
	}

	to := ecommon.HexToAddress(b.cfg.TokenRegistry)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.tokenRegistry.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	chainID := big.NewInt(0)

	if err = outputs.Copy(&chainID, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return chainID, nil
}

func (b *Bridge) GetChainName(chain *big.Int) (string, error) {
	method := "getChainName"
	input, err := b.tokenRegistry.Pack(method, chain)
	if err != nil {
		return "", err
	}

	to := ecommon.HexToAddress(b.cfg.TokenRegistry)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.tokenRegistry.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return "", errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	var name string
	if err = outputs.Copy(&name, unpack); err != nil {
		return "", errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return name, nil
}

func (b *Bridge) GetTokenAddress(chainID *big.Int, name string) ([]byte, error) {
	method := "getTokenAddressByNickname"
	input, err := b.tokenRegistry.Pack(method, chainID, name)
	if err != nil {
		return nil, err
	}

	to := ecommon.HexToAddress(b.cfg.TokenRegistry)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.tokenRegistry.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	address := make([]byte, 0)
	if err = outputs.Copy(&address, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return address, nil
}
