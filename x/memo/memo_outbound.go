package thorchain

import (
	"fmt"
	"strings"
)

type OutboundMemo struct {
	MemoBase
	ChainID   string
	Token     string
	Receiver  string
	Amount    string
	Affiliate string
}

func (m OutboundMemo) GetChain() string {
	return m.ChainID
}

func (m OutboundMemo) GetToken() string {
	return strings.ToUpper(m.Token)
}

func (m OutboundMemo) GetDestination() string { return m.Receiver }

func (m OutboundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s", m.TxType.String(), m.ChainID, m.Token, m.Receiver, m.Amount, m.Affiliate)
}

func NewOutboundMemo(chainID, token, receiver, amount, affiliate string) OutboundMemo {
	return OutboundMemo{
		MemoBase:  MemoBase{TxType: TxOutbound},
		ChainID:   chainID,
		Token:     token,
		Receiver:  receiver,
		Amount:    amount,
		Affiliate: affiliate,
	}
}

func (p *parser) ParseOutboundMemo() (OutboundMemo, error) {
	chainID := p.get(1)
	token := p.get(2)
	receiver := p.get(3)
	amount := p.get(4)
	affiliate := p.get(5)
	return NewOutboundMemo(chainID, token, receiver, amount, affiliate), p.Error()
}
