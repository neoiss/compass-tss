package mapo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/structure"
	stypes "github.com/mapprotocol/compass-tss/x/types"
	"github.com/pkg/errors"
	"math/big"
)

func (b *Bridge) FetchActiveNodes() ([]common.PubKey, error) {
	// todo handler
	na, err := b.GetNodeAccounts(b.epoch)
	if err != nil {
		return nil, fmt.Errorf("fail to get node accounts: %w", err)
	}
	active := make([]common.PubKey, 0)
	for _, item := range na {
		if stypes.NodeStatus(item.Status) == stypes.NodeStatus_Active {
			// todo handler
			active = append(active, common.PubKey(ecommon.Bytes2Hex(item.Secp256Pubkey)))
		}
	}
	return active, nil
}

// FetchNodeStatus get current node status from mapBridge
func (b *Bridge) FetchNodeStatus() (stypes.NodeStatus, error) {
	addr, err := b.keys.GetEthAddress()
	if err != nil {
		return stypes.NodeStatus_Unknown, nil
	}
	// done
	na, err := b.GetNodeAccount(addr.String())
	if err != nil {
		return stypes.NodeStatus_Unknown, fmt.Errorf("failed to get node status: %w", err)
	}

	return stypes.NodeStatus(na.Status), nil
}

// todo handler address to string

// GetNodeAccount retrieves node account for this address from mapBridge
func (b *Bridge) GetNodeAccount(addr string) (*structure.MaintainerInfo, error) {
	method := "getMaitainer"
	input, err := b.mainAbi.Pack(method, ecommon.HexToAddress(addr))
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

	type Back struct {
		Info structure.MaintainerInfo `json:"info"`
	}
	ret := Back{}

	if err = outputs.Copy(&ret, unpack); err != nil {
		return nil, errors.Wrap(err, "copy output")
	}
	b.logger.Info().Interface("ret", ret).Msg("GetNodeAccount back")
	return &ret.Info, nil
}

// GetNodeAccounts retrieves all node accounts from mapBridge
func (b *Bridge) GetNodeAccounts(epoch *big.Int) ([]structure.MaintainerInfo, error) {
	method := "getMaitainers"
	input, err := b.mainAbi.Pack(method, epoch)
	if err != nil {
		return nil, err
	}

	to := ecommon.HexToAddress(b.cfg.Maintainer)
	outPut, err := b.ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: constants.ZeroAddress,
		To:   &to,
		Data: input,
	}, nil)

	outputs := b.mainAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(outPut)
	if err != nil {
		return nil, errors.Wrap(err, "unpack output")
	}

	type Back struct {
		Info []structure.MaintainerInfo `json:"info"`
	}
	var ret Back
	if err = outputs.Copy(&ret, unpack); err != nil {
		return nil, errors.Wrap(err, "copy output")
	}
	return ret.Info, nil
}

func (b *Bridge) GetEpochInfo(epoch *big.Int) (*structure.EpochInfo, error) {
	method := "getEpochInfo"
	input, err := b.mainAbi.Pack(method, epoch)
	if err != nil {
		return nil, err
	}

	to := ecommon.HexToAddress(b.cfg.Maintainer)
	output, err := b.ethClient.CallContract(
		context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &to,
			Data: input,
		},
		nil,
	)

	outputs := b.mainAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	ret := struct {
		Info structure.EpochInfo
	}{}

	if err = outputs.Copy(&ret, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return &ret.Info, nil
}
