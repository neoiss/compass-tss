package types

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

const emptyHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

type TxIn struct {
	Count                string       `json:"count"`
	Chain                common.Chain `json:"chain"`
	TxArray              []*TxInItem  `json:"txArray"`
	Filtered             bool         `json:"filtered"`
	MemPool              bool         `json:"mem_pool"` // indicate whether this item is in the mempool or not
	ConfirmationRequired int64        `json:"confirmation_required"`
	// whether this originated from a "instant observation" - e.g. by a member of the signing party
	// immediately after signing, and also has incorrect gas, requiring a re-observation to correct.
	AllowFutureObservation bool   `json:"allow_future_observation"`
	Method                 string `json:"method"`
	MapRelayHash           string `json:"map_relay_hash"`
	PendingCount           int    `json:"pending_count"`
}

type TxInItem struct {
	Tx               string       `json:"tx"`
	Memo             string       `json:"memo"`
	Sender           string       `json:"sender"`
	FromChain        *big.Int     `json:"from_chain"`
	ToChain          *big.Int     `json:"to_chain"`
	Height           *big.Int     `json:"height"`
	Amount           *big.Int     `json:"amount"`
	OrderId          ecommon.Hash `json:"order_id"`
	GasUsed          *big.Int     `json:"gas_used"`
	Token            []byte       `json:"token"`
	Vault            []byte       `json:"vault"`
	From             []byte       `json:"from"`
	To               []byte       `json:"to"`
	Payload          []byte       `json:"payload"`
	Method           string       `json:"method"`
	LogIndex         uint         `json:"log_index"`
	ChainAndGasLimit *big.Int     `json:"chain_and_gas_limit"` // bridgeOut add new fields
	TxOutType        uint8        `json:"tx_out_type"`
	RefundAddr       []byte       `json:"refund_addr"`
	Sequence         *big.Int     `json:"sequence"` // bridgeIn add new fields
}

type TxInStatus byte

const (
	Processing TxInStatus = iota
	Failed
)

func (t *TxInItem) Equals(other *TxInItem) bool {
	if t.Height.Cmp(other.Height) != 0 {
		return false
	}
	if !strings.EqualFold(t.Tx, other.Tx) {
		return false
	}
	if !strings.EqualFold(t.OrderId.Hex(), other.OrderId.Hex()) {
		return false
	}
	if t.LogIndex != other.LogIndex {
		return false
	}
	if ecommon.Bytes2Hex(t.To) != ecommon.Bytes2Hex(other.To) {
		return false
	}

	return true
}

func (t *TxInItem) EqualsObservedTx(other common.ObservedTx) bool {
	// Do not compare block height, as we only keep one deck item for a tx pre/post finalization.
	// The final block height will be ConfirmationCount higher than the unfinalized tx.
	txId, err := common.NewTxID(t.Tx)
	if err != nil {
		return false
	}
	if !txId.Equals(other.Tx.ID) {
		return false
	}
	if ecommon.Bytes2Hex(t.To) != other.Tx.ToAddress.String() {
		return false
	}
	return true
}

func (t *TxInItem) Copy() *TxInItem {
	return &TxInItem{
		Tx:      t.Tx,
		ToChain: t.ToChain,
		Height:  t.Height,
		Amount:  t.Amount,
		OrderId: t.OrderId,
		GasUsed: t.GasUsed,
		Token:   t.Token,
		Vault:   t.Vault,
		To:      t.To,
		Method:  t.Method,
	}
}

// IsEmpty return true only when every field in TxInItem is empty
func (t *TxInItem) IsEmpty() bool {
	return (t.Height == nil || t.Height.Uint64() == 0) &&
		t.Tx == "" &&
		len(t.To) == 0 &&
		len(t.Token) == 0 &&
		len(t.Vault) == 0 && t.OrderId.String() == emptyHash
}

// CacheHash calculate the has used for signer cache
func (t *TxInItem) CacheHash(chain common.Chain, inboundHash string) string {
	str := fmt.Sprintf("%s|%s|%s|%s|%s", chain, t.To, t.OrderId.Hex(), t.Height.String(), inboundHash)
	return fmt.Sprintf("%X", sha256.Sum256([]byte(str)))
}

func (t *TxInItem) CacheVault(chain common.Chain) string {
	return InboundCacheKey(t.Tx+"-"+t.OrderId.String(), chain.String())
}

// GetTotalTransactionValue return the total value of the requested asset
func (t TxIn) GetTotalTransactionValue(asset common.Asset, excludeFrom []common.Address) cosmos.Uint {
	total := cosmos.ZeroUint()
	//for _, item := range t.TxArray {
	//	fromAsgard := false
	//	for _, fromAddress := range excludeFrom {
	//		if strings.EqualFold(fromAddress.String(), item.Sender) {
	//			fromAsgard = true
	//		}
	//	}
	//	if fromAsgard {
	//		continue
	//	}
	//	// skip confirmation counting if it is internal tx
	//	m, err := mem.ParseMemo(common.LatestVersion, item.Memo)
	//	if err == nil && m.IsInternal() {
	//		continue
	//	}
	//	c := item.Coins.GetCoin(asset)
	//	if c.IsEmpty() {
	//		continue
	//	}
	//	total = total.Add(c.Amount)
	//}
	return total
}

// GetTotalGas return the total gas // todo
func (t TxIn) GetTotalGas() cosmos.Uint {
	total := cosmos.ZeroUint()
	for _, item := range t.TxArray {
		if item.GasUsed == nil {
			continue
		}
		//if err := item.Gas.Valid(); err != nil {
		//	continue
		//}
		//total = total.Add(item.Gas[0].Amount)
	}
	return total
}

func InboundCacheKey(vault, chain string) string {
	return fmt.Sprintf("inbound-%s-%s", vault, chain)
}
