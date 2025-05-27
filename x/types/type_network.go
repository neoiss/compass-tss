package types

import (
	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
)

// NewNetwork create a new instance Network it is empty though
func NewNetwork() Network {
	return Network{
		BondRewardRune: cosmos.ZeroUint(),
		TotalBondUnits: cosmos.ZeroUint(),
	}
}

// CalcNodeRewards calculate node rewards
func (m *Network) CalcNodeRewards(nodeUnits cosmos.Uint) cosmos.Uint {
	return common.GetUncappedShare(nodeUnits, m.TotalBondUnits, m.BondRewardRune)
}
