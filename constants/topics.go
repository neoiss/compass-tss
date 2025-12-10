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

// relay event
const (
	EventOfBridgeRelay       EventSig = "BridgeRelay(bytes32,uint256,uint8,bytes,bytes,bytes,uint256,uint256,bytes32,bytes,bytes)"
	EventOfBridgeCompleted   EventSig = "BridgeCompleted(bytes32,uint256,uint8,bytes,uint256,address,bytes)"
	EventOfBridgeRelaySigned EventSig = "BridgeRelaySigned(bytes32,uint256,bytes,bytes,bytes)"
)

// src or dst chain event
const (
	EventOfBridgeOut EventSig = "BridgeOut(bytes32,uint256,uint8,bytes,address,uint256,address,address,bytes,bytes)"  // -> tss_manager voteTxIn
	EventOfBridgeIn  EventSig = "BridgeIn(bytes32,uint256,uint8,bytes,uint256,address,address,uint256,address,bytes)" // -> tss_manager voteTxOut
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
