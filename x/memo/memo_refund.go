package thorchain

import (
	"fmt"
)

type RefundMemo struct {
	MemoBase
	Chain  string
	TxHash string
}

func (m RefundMemo) GetTxHash() string { return m.TxHash }

// String implement fmt.Stringer
func (m RefundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.TxHash)
}

// NewRefundMemo create a new RefundMemo
func NewRefundMemo(chain, txHash string) RefundMemo {
	return RefundMemo{
		MemoBase: MemoBase{TxType: TxRefund},
		Chain:    chain,
		TxHash:   txHash,
	}
}

func (p *parser) ParseRefundMemo() (RefundMemo, error) {
	chain := p.get(1)
	txHash := p.get(2)
	return NewRefundMemo(chain, txHash), p.Error()
}
