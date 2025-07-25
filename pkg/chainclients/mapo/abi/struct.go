package abi

import (
	"github.com/mapprotocol/compass-tss/constants"
	"math/big"
)

//type VoteTxIn struct {
//	TxInItem TxIn
//}

type VoteTxIn struct {
	TxInType  constants.TxInType
	ToChain   *big.Int
	Height    *big.Int
	FromChain *big.Int
	Amount    *big.Int
	OrderId   [32]byte
	Vault     []byte
	Token     []byte
	From      []byte
	To        []byte
	Payload   []byte
}

type Gas struct {
	Chain          *big.Int
	Pubkey         []byte
	Router         []byte
	GasRate        *big.Int
	TxSize         *big.Int
	TxSizeWithCall *big.Int
}
