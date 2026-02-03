package memo

import (
	"fmt"
)

type MigrateMemo struct {
	MemoBase
	Chain   string
	OrderID string
}

func (m MigrateMemo) GetOrderID() string { return m.OrderID }

// String returns a string representation of the memo
// format: M>|from chain|order id
func (m MigrateMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.OrderID)
}

func NewMigrateMemo(chain, orderID string) InboundMemo {
	return InboundMemo{
		MemoBase: MemoBase{TxType: TxMigrate},
		Chain:    chain,
		OrderID:  orderID,
	}
}

func (p *parser) ParseMigrateMemo() (InboundMemo, error) {
	chain := p.get(1)
	order := p.get(2)
	return NewMigrateMemo(chain, order), p.Error()
}
