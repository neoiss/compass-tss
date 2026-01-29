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

type VaultStatus int32

const (
	VaultStatus_InactiveVault VaultStatus = 0
	VaultStatus_ActiveVault   VaultStatus = 1
	VaultStatus_RetiringVault VaultStatus = 2
	VaultStatus_InitVault     VaultStatus = 3
)
