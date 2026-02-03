package memo

import (
	"fmt"
)

type AddLiquidityMemo struct {
	MemoBase
	Receiver string
}

func (m AddLiquidityMemo) GetDestination() string { return m.Receiver }

func (m AddLiquidityMemo) String() string {
	return fmt.Sprintf("%s|%s", m.TxType.String(), m.Receiver)
}

func NewAddLiquidityMemo(receiver string) AddLiquidityMemo {
	return AddLiquidityMemo{
		MemoBase: MemoBase{TxType: TxAdd},
		Receiver: receiver,
	}
}

func (p *parser) ParseAddLiquidityMemo() (AddLiquidityMemo, error) {
	receiver := p.get(1)
	return NewAddLiquidityMemo(receiver), p.Error()
}
