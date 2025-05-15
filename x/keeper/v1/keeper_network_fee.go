package keeperv1

import (
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func (k KVStore) setNetworkFee(ctx cosmos.Context, key string, record NetworkFee) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getNetworkFee(ctx cosmos.Context, key string, record *NetworkFee) (bool, error) {
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

// GetNetworkFee get the network fee of the given chain from kv store , if it doesn't exist , it will create an empty one
func (k KVStore) GetNetworkFee(ctx cosmos.Context, chain common.Chain) (NetworkFee, error) {
	record := NetworkFee{
		Chain:              chain,
		TransactionSize:    0,
		TransactionFeeRate: 0,
	}
	_, err := k.getNetworkFee(ctx, k.GetKey(prefixNetworkFee, chain.String()), &record)
	return record, err
}

// SaveNetworkFee save the network fee to kv store
func (k KVStore) SaveNetworkFee(ctx cosmos.Context, chain common.Chain, networkFee NetworkFee) error {
	if err := networkFee.Valid(); err != nil {
		return err
	}
	k.setNetworkFee(ctx, k.GetKey(prefixNetworkFee, chain.String()), networkFee)
	return nil
}

// GetNetworkFeeIterator
func (k KVStore) GetNetworkFeeIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixNetworkFee)
}
