package mapo

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

func (b *Bridge) GetChainID(name string) (*big.Int, error) {
	if name == "" {
		return nil, errors.New("chain name is empty")
	}
	method := "getChainByName"
	input, err := b.tokenRegistry.Pack(method, name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
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
	if chain == nil {
		return "", errors.New("chain is nil")
	}
	method := "getChainName"
	input, err := b.tokenRegistry.Pack(method, chain)
	if err != nil {
		return "", errors.Wrapf(err, "unable to pack input of %s", method)
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
	if chainID == nil {
		return nil, errors.New("chainID is nil")
	}
	if name == "" {
		return nil, errors.New("token name is empty")
	}
	method := "getTokenAddressByNickname"
	input, err := b.tokenRegistry.Pack(method, chainID, strings.ToUpper(name))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
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
	if len(address) == 0 {
		return nil, fmt.Errorf("unsupported token(%d:%s)", chainID, name)
	}
	return address, nil
}

func (b *Bridge) GetTokenDecimals(chainID *big.Int, address []byte) (*big.Int, error) {
	if chainID == nil {
		return nil, errors.New("chainID is nil")
	}
	if address == nil {
		return nil, errors.New("token address is nil")
	}
	method := "getTokenDecimals"
	input, err := b.tokenRegistry.Pack(method, chainID, address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
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

	decimals := big.NewInt(0)
	if err = outputs.Copy(&decimals, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return decimals, nil
}
