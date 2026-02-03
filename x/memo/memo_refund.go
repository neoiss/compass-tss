package memo

import (
	"fmt"
)

type RefundMemo struct {
	MemoBase
	Chain   string
	OrderID string
}

func (m RefundMemo) GetOrderID() string { return m.OrderID }

// String implement fmt.Stringer
func (m RefundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.OrderID)
}

// NewRefundMemo create a new RefundMemo
func NewRefundMemo(chain, orderID string) RefundMemo {
	return RefundMemo{
		MemoBase: MemoBase{TxType: TxRefund},
		Chain:    chain,
		OrderID:  orderID,
	}
}

func (p *parser) ParseRefundMemo() (RefundMemo, error) {
	chain := p.get(1)
	orderID := p.get(2)
	return NewRefundMemo(chain, orderID), p.Error()
}
