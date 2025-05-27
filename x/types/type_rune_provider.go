package types

import (
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func NewRUNEProvider(addr cosmos.AccAddress) RUNEProvider {
	return RUNEProvider{
		RuneAddress: addr,
		Units:       cosmos.ZeroUint(),
	}
}

func (rp RUNEProvider) Key() string {
	return rp.RuneAddress.String()
}
