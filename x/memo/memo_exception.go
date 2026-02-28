package memo

import (
	"fmt"
)

type ExceptionMemo struct {
	MemoBase
	Chain   string
	OrderID string
}

func (m ExceptionMemo) GetOrderID() string { return m.OrderID }

// String implement fmt.Stringer
func (m ExceptionMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.OrderID)
}

// NewExceptionMemo create a new ExceptionMemo
func NewExceptionMemo(chain, orderID string) ExceptionMemo {
	return ExceptionMemo{
		MemoBase: MemoBase{TxType: TxException},
		Chain:    chain,
		OrderID:  orderID,
	}
}

func (p *parser) ParseExceptionMemo() (ExceptionMemo, error) {
	chain := p.get(1)
	orderID := p.get(2)
	return NewExceptionMemo(chain, orderID), p.Error()
}
