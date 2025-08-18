package thorchain

import (
	"fmt"
	"strings"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

// TXTYPE:STATE1:STATE2:STATE3:FINALMEMO

type TxType uint8

const (
	TxUnknown TxType = iota
	TxInbound
	TxOutbound
	TxRefund
	TxAdd
)

var txToStringMap = map[TxType]string{
	TxInbound:  "M>",
	TxOutbound: "Mx",
	TxRefund:   "M<",
	TxAdd:      "M+",
}

var stringToTxTypeMap = map[string]TxType{
	"m>": TxInbound,
	"mx": TxOutbound,
	"m<": TxRefund,
	"m+": TxAdd,
}

func StringToTxType(s string) (TxType, error) {
	sl := strings.ToLower(s)
	if t, ok := stringToTxTypeMap[sl]; ok {
		return t, nil
	}

	return TxUnknown, fmt.Errorf("invalid tx type: %s", s)
}

func (tx TxType) IsEmpty() bool {
	return tx == TxUnknown
}

func (tx TxType) Equals(tx2 TxType) bool {
	return tx == tx2
}

// Converts a txType into a string
func (tx TxType) String() string {
	return txToStringMap[tx]
}

type Memo interface {
	String() string
	IsType(tx TxType) bool
	GetType() TxType
	IsEmpty() bool
	GetChain() string
	GetToken() string
	GetAmount() cosmos.Uint
	GetDestination() string
	GetTxHash() string
	IsValid() bool
}

type MemoBase struct {
	TxType TxType
	Asset  common.Asset
}

var EmptyMemo = MemoBase{TxType: TxUnknown, Asset: common.EmptyAsset}

func (m MemoBase) String() string         { return "" }
func (m MemoBase) GetType() TxType        { return m.TxType }
func (m MemoBase) IsType(tx TxType) bool  { return m.TxType.Equals(tx) }
func (m MemoBase) GetChain() string       { return "" }
func (m MemoBase) GetToken() string       { return "" }
func (m MemoBase) GetAmount() cosmos.Uint { return cosmos.ZeroUint() }
func (m MemoBase) GetDestination() string { return "" }
func (m MemoBase) GetTxHash() string      { return "" }
func (m MemoBase) IsEmpty() bool          { return m.TxType.IsEmpty() }
func (m MemoBase) IsValid() bool {
	_, ok := txToStringMap[m.GetType()]
	return ok
}

func ParseMemo(memo string) (mem Memo, err error) {
	defer func() {
		if r := recover(); r != nil {
			mem = EmptyMemo
			err = fmt.Errorf("panicked parsing memo(%s), err: %s", memo, r)
		}
	}()

	parser, err := newParser(memo)
	if err != nil {
		return EmptyMemo, err
	}

	return parser.parse()
}
