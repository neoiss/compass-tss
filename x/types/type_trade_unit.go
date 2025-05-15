package types

import (
	"github.com/mapprotocol/compass-tss/common"
	"gitlab.com/thorchain/thornode/v3/common/cosmos"
)

func NewTradeUnit(asset common.Asset) TradeUnit {
	return TradeUnit{
		Asset: asset,
		Units: cosmos.ZeroUint(),
		Depth: cosmos.ZeroUint(),
	}
}

func (tu TradeUnit) Key() string {
	return tu.Asset.String()
}
