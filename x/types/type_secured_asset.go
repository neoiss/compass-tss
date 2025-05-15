package types

import (
	"github.com/mapprotocol/compass-tss/common"
	"gitlab.com/thorchain/thornode/v3/common/cosmos"
)

func NewSecuredAsset(asset common.Asset) SecuredAsset {
	return SecuredAsset{
		Asset: asset,
		Depth: cosmos.ZeroUint(),
	}
}

func (tu SecuredAsset) Key() string {
	return tu.Asset.String()
}
