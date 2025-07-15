package constants

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

const (
	EventOfMessageOut = "0x469059a9fd182ad3741bdd67b925e15056d35262609ea83393db7e8fb5a05ab1"
)

// relay event
const (
	RelayEventOfMigration    EventSig = "Migration(bytes32,uint256,bytes,bytes,TokenAllowance[],uint256,uint256)"            // -> dst gateway 合约 transferAllowance
	RelayEventOfTransferOut  EventSig = "RelayTransferOut(bytes32,bytes,uint256,uint256,bytes,bytes,uint256,uint256)"        // -> dst gateway 合约 transferOut
	RelayEventOfTransferCall EventSig = "RelayTransferCall(bytes32,bytes,uint256,uint256,bytes,bytes,bytes,uint256,uint256)" // -> dst gateway 合约 transferOutCall
)

// src or dst chain event
const (
	EventOfDeposit EventSig = "Deposit(bytes32,address,address,address,uint256,address)"          // -> maintainer 合约 voteTxIn
	EventOfSwap    EventSig = "Swap(bytes32,address,address,address,uint256,uint256,bytes,bytes)" // -> maintainer 合约 voteTxIn
)

const (
	EventOfTransferOut       EventSig = "TransferOut(bytes32,address,address,uint256,address,bool)"      // -> maintainer 合约 voteTxOut
	EventOfTransferAllowance EventSig = "TransferAllowance(bytes32,address,address,EVMTokenAllowance[])" // -> maintainer 合约 voteTxOut
)

type TxInType uint8

const (
	DEPOSIT TxInType = iota
	SWAP
)
