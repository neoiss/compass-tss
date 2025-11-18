package thorchain

import (
	"fmt"
)

type InboundMemo struct {
	MemoBase
	Chain   string
	OrderID string
}

func (m InboundMemo) GetOrderID() string { return m.OrderID }

// String returns a string representation of the memo
// format: M>|from chain|order id
func (m InboundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.OrderID)
}

func NewInboundMemo(chain, orderID string) InboundMemo {
	return InboundMemo{
		MemoBase: MemoBase{TxType: TxInbound},
		Chain:    chain,
		OrderID:  orderID,
	}
}

func (p *parser) ParseInboundMemo() (InboundMemo, error) {
	chain := p.get(1)
	tx := p.get(2)
	return NewInboundMemo(chain, tx), p.Error()
}
