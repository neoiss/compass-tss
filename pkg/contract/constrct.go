package contract

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/pkg/abi"
)

type Call struct {
	abi  *abi.Abi
	toC  []common.Address
	conn *ethclient.Client
}

func New(conn *ethclient.Client, addr []common.Address, abi *abi.Abi) *Call {
	return &Call{
		conn: conn,
		toC:  addr,
		abi:  abi,
	}
}

func (c *Call) Call(method string, ret interface{}, idx int, params ...interface{}) error {
	input, err := c.abi.PackInput(method, params...)
	if err != nil {
		return err
	}

	outPut, err := c.conn.CallContract(context.Background(),
		ethereum.CallMsg{
			From: constants.ZeroAddress,
			To:   &c.toC[idx],
			Data: input,
		},
		nil,
	)
	if err != nil {
		return err
	}

	return c.abi.UnpackOutput(method, ret, outPut)
}
