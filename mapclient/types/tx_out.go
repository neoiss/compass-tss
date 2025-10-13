package types

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/common"
)

// TxOutItem represent the information of a tx bifrost need to process
type TxOutItem struct {
	Height          int64        `json:"height,omitempty"`
	Memo            string       `json:"memo,omitempty"`
	OrderId         ecommon.Hash `json:"order_id,omitempty"`
	Token           []byte       `json:"token,omitempty"`
	Vault           []byte       `json:"vault,omitempty"`
	To              []byte       `json:"to,omitempty"`
	Amount          *big.Int     `json:"amount,omitempty"`
	FromChain       *big.Int     `json:"from_chain,omitempty"`
	Chain           *big.Int     `json:"chain,omitempty"`
	TransactionRate *big.Int     `json:"transaction_rate,omitempty"`
	TransactionSize *big.Int     `json:"transaction_size,omitempty"`
	LogIndex        uint         `json:"log_index,omitempty"`
	TxHash          string       `json:"tx_hash,omitempty"`
	InTxHash        string       `json:"in_tx_hash"`
	Method          string       `json:"method,omitempty"`
	// relayTransferCall
	Payload []byte `json:"payload,omitempty"`
	// migration
	FromVault   []byte        `json:"from_vault,omitempty"`
	ToVault     []byte        `json:"to_vault,omitempty"`
	VaultPubKey common.PubKey `json:"vault_pubkey"`
	Checkpoint  []byte        `json:"-"`
	// bridgeRelay add new field
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Sequence         *big.Int
	From             []byte
	Data             []byte
}

// Hash return a sha256 hash that can uniquely represent the TxOutItem
func (tx TxOutItem) Hash() string {
	str := fmt.Sprintf("%s|%s|%s|%s|%s|%s", tx.Chain, ecommon.Bytes2Hex(tx.To),
		ecommon.Bytes2Hex(tx.Vault), tx.Amount.String(), tx.Memo, tx.TxHash)
	return fmt.Sprintf("%X", sha256.Sum256([]byte(str)))
}

// CacheHash return a hash that doesn't include VaultPubKey , thus this one can be used as cache key for txOutItem across different vaults
func (tx TxOutItem) CacheHash() string {
	str := fmt.Sprintf("%s|%s|%s|%s|%s", tx.Chain, ecommon.Bytes2Hex(tx.To),
		tx.Amount.String(), tx.Memo, tx.TxHash)
	return fmt.Sprintf("%X", sha256.Sum256([]byte(str)))
}

func (tx TxOutItem) CacheVault(chain common.Chain) string {
	return BroadcastCacheKey(ecommon.Bytes2Hex(tx.Vault), chain.String())
}

// Equals returns true when the TxOutItems are equal.
//
// NOTE: The height field should NOT be compared. This is necessary to pass through on
// the TxOutItem to the unstuck routine to determine the position within the signing
// period, but should not be used to determine equality for deduplication.
func (tx TxOutItem) Equals(tx2 TxOutItem) bool {
	if tx.Chain.Cmp(tx2.Chain) != 0 {
		return false
	}
	if tx.TransactionRate.Cmp(tx2.TransactionRate) != 0 {
		return false
	}
	if tx.TransactionSize.Cmp(tx2.TransactionSize) != 0 {
		return false
	}
	if tx.LogIndex != tx2.LogIndex {
		return false
	}
	if tx.Method != tx2.Method {
		return false
	}
	if !strings.EqualFold(tx.TxHash, tx2.TxHash) {
		return false
	}
	if tx.OrderId.String() != tx2.OrderId.String() {
		return false
	}

	return true
}

// TxArrayItem used to represent the tx out item coming from THORChain, there is little difference between TxArrayItem
// and TxOutItem defined above , only Coin <-> Coins field are different.
// TxArrayItem from THORChain has Coin , which only have a single coin
// TxOutItem used in bifrost need to support Coins for multisend
type TxArrayItem struct {
	Memo            string `json:"memo,omitempty"`
	Chain           *big.Int
	TransactionRate *big.Int
	TransactionSize *big.Int
	LogIndex        uint
	TxHash          string
	Method          string
	// bridgeRelay add new field
	OrderId          ecommon.Hash `json:"order_id"`
	ChainAndGasLimit *big.Int
	TxType           uint8
	Vault            []byte
	To               []byte
	Token            []byte
	Amount           *big.Int
	Sequence         *big.Int
	Hash             ecommon.Hash
	From             []byte
	Data             []byte // relaySigned relayData -> data
	Sender           []byte // bridgeCompleted event field
	Signature        []byte // relaySigned event field
}

// TxOutItem convert the information to TxOutItem
func (tx TxArrayItem) TxOutItem(height int64) TxOutItem {
	return TxOutItem{
		Chain:            tx.Chain,
		Memo:             tx.Memo,
		Height:           height,
		OrderId:          tx.OrderId,
		Token:            tx.Token,
		Vault:            tx.Vault,
		To:               tx.To,
		Amount:           tx.Amount,
		TransactionRate:  tx.TransactionRate,
		TransactionSize:  tx.TransactionSize,
		LogIndex:         tx.LogIndex,
		TxHash:           tx.TxHash,
		Method:           tx.Method,
		ChainAndGasLimit: tx.ChainAndGasLimit,
		TxOutType:        tx.TxType,
		Sequence:         tx.Sequence,
		From:             tx.From,
		Data:             tx.Data,
	}
}

// TxOut represent the tx out information , bifrost need to sign and process
type TxOut struct {
	Height  int64         `json:"height"`
	TxArray []TxArrayItem `json:"tx_array"`
}

func BroadcastCacheKey(vault, chain string) string {
	return fmt.Sprintf("broadcast-%s-%s", vault, chain)
}
