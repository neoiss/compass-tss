package mapo

import (
	"context"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

func (b *Bridge) GetAffiliateIDByName(name string) (uint16, error) {
	if name == "" {
		return 0, errors.New("name is empty")
	}
	method := "getInfoByNickname"
	input, err := b.affiliateFeeAbi.Pack(method, name)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to pack affiliate fee input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.AffiliateFeeManager)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.affiliateFeeAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	affiliate := struct {
		Info struct {
			Id       uint16
			BaseRate uint16
			MaxRate  uint16
			Wallet   ecommon.Address
			Nickname string
		}
	}{}

	if err = outputs.Copy(&affiliate, unpack); err != nil {
		return 0, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	if affiliate.Info.Id == 0 {
		return 0, errors.New("affiliate not found")
	}
	return affiliate.Info.Id, nil
}

func (b *Bridge) GetAffiliateIDByAlias(name string) (uint16, error) {
	if name == "" {
		return 0, errors.New("name is empty")
	}
	method := "getInfoByShortName"
	input, err := b.affiliateFeeAbi.Pack(method, name)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.AffiliateFeeManager)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.affiliateFeeAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	affiliate := struct {
		Info struct {
			Id       uint16
			BaseRate uint16
			MaxRate  uint16
			Wallet   ecommon.Address
			Nickname string
		}
	}{}

	if err = outputs.Copy(&affiliate, unpack); err != nil {
		return 0, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	if affiliate.Info.Id == 0 {
		return 0, errors.New("affiliate not found")
	}
	return affiliate.Info.Id, nil
}
