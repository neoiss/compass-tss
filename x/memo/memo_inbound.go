package thorchain

import (
	"fmt"
)

type InboundMemo struct {
	MemoBase
	Chain  string
	TxHash string
}

func (m InboundMemo) GetTxHash() string { return m.TxHash }

func (m InboundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s", m.TxType.String(), m.Chain, m.TxHash)
}

func NewInboundMemo(chain, txHash string) InboundMemo {
	return InboundMemo{
		MemoBase: MemoBase{TxType: TxInbound},
		Chain:    chain,
		TxHash:   txHash,
	}
}

func (p *parser) ParseInboundMemo() (InboundMemo, error) {
	chain := p.get(1)
	tx := p.get(2)
	return NewInboundMemo(chain, tx), p.Error()
}
