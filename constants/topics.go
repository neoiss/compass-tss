package constants

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

func (es EventSig) String() string {
	return string(es)
}

const (
	EventOfMessageOut = "0x469059a9fd182ad3741bdd67b925e15056d35262609ea83393db7e8fb5a05ab1"
)

// relay event
const (
	EventOfBridgeRelay     EventSig = "BridgeRelay(bytes32,uint256,uint8,bytes,uint256,bytes,uint256,bytes,bytes,bytes)" // -> dst gateway 合约 transferAllowance
	EventOfBridgeCompleted EventSig = "BridgeCompleted(bytes32,uint256,uint8,bytes,uint256,address,bytes)"               // -> dst gateway 合约 transferOut
	EventOfBridgeInRelay   EventSig = "BridgeCompleted(bytes32,uint256,uint8,bytes,uint256,address,bytes)"               // 目标链是map，弹出此事件
)

// src or dst chain event
const (
	EventOfBridgeOut EventSig = "BridgeOut(bytes32,uint256,uint8,bytes,address,uint256,address,bytes,bytes)"                // -> tss_manager 合约 voteTxIn
	EventOfBridgeIn  EventSig = "BridgeIn(bytes32,uint256,uint8,bytes,uint256,address,address,uint256,bytes,address,bytes)" // -> tss_manager 合约 voteTxOut
)

type TxInType uint8

const (
	DEPOSIT TxInType = iota
	TRANSFER
	MIGRATE
	REFUND
	MESSAGE
)

func (tx TxInType) String() string {
	switch tx {
	case DEPOSIT:
		return "DEPOSIT"
	case TRANSFER:
		return "TRANSFER"
	case MIGRATE:
		return "MIGRATE"
	case REFUND:
		return "REFUND"
	case MESSAGE:
		return "MESSAGE"
	}
	return "UNKNOWN"
}
