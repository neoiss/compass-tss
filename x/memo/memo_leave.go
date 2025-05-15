package thorchain

import (
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

type LeaveMemo struct {
	MemoBase
	NodeAddress cosmos.AccAddress
}

func (m LeaveMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }

func NewLeaveMemo(addr cosmos.AccAddress) LeaveMemo {
	return LeaveMemo{
		MemoBase:    MemoBase{TxType: TxLeave},
		NodeAddress: addr,
	}
}

func (p *parser) ParseLeaveMemo() (LeaveMemo, error) {
	addr := p.getAccAddress(1, true, nil)
	return NewLeaveMemo(addr), p.Error()
}
