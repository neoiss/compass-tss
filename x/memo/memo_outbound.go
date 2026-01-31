package thorchain

import (
	"fmt"
	"github.com/mapprotocol/compass-tss/common"
	"math/big"
	"strings"
)

type Affiliates []*Affiliate

type Affiliate struct {
	Name       string
	Bps        *big.Int
	Compressed bool
}

type OutboundMemo struct {
	MemoBase
	Chain      string
	Token      string
	Receiver   string
	Amount     *big.Int
	Affiliates Affiliates
}

func (m OutboundMemo) GetChain() common.Chain {
	return common.Chain(m.Chain)
}

func (m OutboundMemo) GetToken() string {
	return strings.ToUpper(m.Token)
}

func (m OutboundMemo) GetDestination() string { return m.Receiver }

func (m OutboundMemo) GetAffiliates() Affiliates { return m.Affiliates }

func (m OutboundMemo) String() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s", m.TxType.String(), m.Chain, m.Token, m.Receiver, m.Amount, m.Affiliates)
}

func NewOutboundMemo(chain, token, receiver string, amount *big.Int, affiliates Affiliates) OutboundMemo {
	return OutboundMemo{
		MemoBase:   MemoBase{TxType: TxOutbound},
		Chain:      chain,
		Token:      token,
		Receiver:   receiver,
		Amount:     amount,
		Affiliates: affiliates,
	}
}

func (p *parser) ParseOutboundMemo() (OutboundMemo, error) {
	chainID := p.get(1)
	token := p.get(2)
	receiver := p.get(3)
	amount := p.getMinAmount(4)
	affiliates := p.getAffiliates(5)
	return NewOutboundMemo(chainID, token, receiver, amount, affiliates), p.Error()
}

func (as Affiliates) String() string {
	if len(as) == 0 {
		return ""
	}

	first := as[0]
	if first == nil {
		return ""
	}

	if first.Compressed {
		var result strings.Builder
		for _, aff := range as {
			result.WriteString(aff.Name)
			result.WriteString(aff.Bps.String())
		}
		return result.String()
	}

	var result strings.Builder
	for i, aff := range as {
		if aff != nil {
			if i > 0 {
				result.WriteString(",")
			}
			result.WriteString(aff.Name)
			result.WriteString(":")
			result.WriteString(aff.Bps.String())
		}
	}
	return result.String()
}
