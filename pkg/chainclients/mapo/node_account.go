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
)

func (b *Bridge) FetchActiveNodes() ([]common.PubKey, error) {
	na, err := b.GetNodeAccounts()
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
	signerAddr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		return stypes.NodeStatus_Unknown, fmt.Errorf("fail to get signer address: %w", err)
	}
	bepAddr := signerAddr.String()
	if len(bepAddr) == 0 {
		return stypes.NodeStatus_Unknown, errors.New("bep address is empty")
	}
	// todo handler
	na, err := b.GetNodeAccount("0x2b7588165556aB2fA1d30c520491C385BAa424d8")
	if err != nil {
		return stypes.NodeStatus_Unknown, fmt.Errorf("failed to get node status: %w", err)
	}

	return stypes.NodeStatus(na.Status), nil
}

// todo handler address to string

// GetNodeAccount retrieves node account for this address from mapBridge
func (b *Bridge) GetNodeAccount(thorAddr string) (*structure.MaintainerInfo, error) {
	method := "getMaitainer"
	input, err := b.mainAbi.Pack(method, ecommon.HexToAddress(thorAddr))
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
		Info *structure.MaintainerInfo `json:"info"`
	}
	ret := Back{}
	if err = outputs.Copy(&ret, unpack); err != nil {
		return nil, errors.Wrap(err, "copy output")
	}
	fmt.Println("GetNodeAccount ret ----------- ", ret)
	return ret.Info, nil
}

// GetNodeAccounts retrieves all node accounts from mapBridge
func (b *Bridge) GetNodeAccounts() ([]*structure.MaintainerInfo, error) {
	method := "getMaitainers"
	input, err := b.mainAbi.Pack(method)
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
		Info []*structure.MaintainerInfo `json:"info"`
	}
	var ret Back
	if err = outputs.Copy(&ret, unpack); err != nil {
		return nil, errors.Wrap(err, "copy output")
	}
	fmt.Println("GetNodeAccounts ret ----------- ", ret)
	return ret.Info, nil
}
