package ethereum

import (
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/x/aggregators"
)

func LatestAggregatorContracts() []common.Address {
	addrs := []common.Address{}
	for _, agg := range aggregators.DexAggregators(common.LatestVersion) {
		if agg.Chain.Equals(common.ETHChain) {
			addrs = append(addrs, common.Address(agg.Address))
		}
	}
	return addrs
}
