package constants

import "github.com/ethereum/go-ethereum/common"

var (
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

type NodeStatus uint8

const (
	NodeStatus_Unknown NodeStatus = iota
	NodeStatus_Whitelisted
	NodeStatus_Standby
	NodeStatus_Ready
	NodeStatus_Active
	NodeStatus_Disabled
)
