package keeperv1

import (
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func (k KVStore) setObservedNetworkFeeVoter(ctx cosmos.Context, key string, record ObservedNetworkFeeVoter) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getObservedNetworkFeeVoter(ctx cosmos.Context, key string, record *ObservedNetworkFeeVoter) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, dbError(ctx, fmt.Sprintf("Unmarshal kvstore: (%T) %s", record, key), err)
	}
	return true, nil
}

// SetObservedNetworkFeeVoter - save a observed network fee voter object
func (k KVStore) SetObservedNetworkFeeVoter(ctx cosmos.Context, networkFeeVoter ObservedNetworkFeeVoter) {
	key := networkFeeVoter.ID()
	k.setObservedNetworkFeeVoter(ctx, k.GetKey(prefixNetworkFeeVoter, key), networkFeeVoter)
}

// GetObservedNetworkFeeVoterIterator iterate tx in voters
func (k KVStore) GetObservedNetworkFeeVoterIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixNetworkFeeVoter)
}

// GetObservedNetworkFeeVoter - gets information of an observed network fee voter
func (k KVStore) GetObservedNetworkFeeVoter(ctx cosmos.Context, height int64, chain common.Chain, rate, size int64) (ObservedNetworkFeeVoter, error) {
	record := NewObservedNetworkFeeVoter(height, chain)
	if rate > 0 {
		record.FeeRate = rate
	}
	if size > 0 {
		record.TransactionSize = size
	}
	key := record.ID()
	_, err := k.getObservedNetworkFeeVoter(ctx, k.GetKey(prefixNetworkFeeVoter, key), &record)
	return record, err
}
