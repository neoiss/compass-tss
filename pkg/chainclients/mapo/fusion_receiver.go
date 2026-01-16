package mapo

import (
	"context"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
	"math/big"
)

func (b *Bridge) GetChainIDFromFusionReceiver(name string) (*big.Int, error) {
	if name == "" {
		return nil, errors.New("chain name is empty")
	}
	method := "chainIdByName"
	input, err := b.fusionReceiverAbi.Pack(method, name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.FusionReceiver)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.fusionReceiverAbi.Methods[method].Outputs
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
