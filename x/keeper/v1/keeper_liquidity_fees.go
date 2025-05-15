package keeperv1

import (
	"fmt"
	"strconv"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

// AddToLiquidityFees - measure of fees collected in each block
func (k KVStore) AddToLiquidityFees(ctx cosmos.Context, asset common.Asset, fee cosmos.Uint) error {
	currentHeight := uint64(ctx.BlockHeight())

	totalFees, err := k.GetTotalLiquidityFees(ctx, currentHeight)
	if err != nil {
		return err
	}
	poolFees, err := k.GetPoolLiquidityFees(ctx, currentHeight, asset)
	if err != nil {
		return err
	}

	totalFees = totalFees.Add(fee)
	poolFees = poolFees.Add(fee)

	// update total liquidity
	k.setUint64(ctx, k.GetKey(prefixTotalLiquidityFee, strconv.FormatUint(currentHeight, 10)), totalFees.Uint64())

	// update pool liquidity
	k.setUint64(ctx, k.GetKey(prefixPoolLiquidityFee, fmt.Sprintf("%d-%s", currentHeight, asset.String())), poolFees.Uint64())

	var currentValue uint64
	currentValue, err = k.GetRollingPoolLiquidityFee(ctx, asset)
	if err != nil {
		ctx.Logger().Error("fail to get existing rolling pool liquidity fee", "error", err)
		return nil
	}
	key := k.GetKey(prefixRollingPoolLiquidityFee, asset.String())
	k.setUint64(ctx, key, currentValue+fee.Uint64())

	return nil
}

// GetRollingPoolLiquidityFee get the given rolling liquidity fee from key value store
func (k KVStore) GetRollingPoolLiquidityFee(ctx cosmos.Context, asset common.Asset) (uint64, error) {
	key := k.GetKey(prefixRollingPoolLiquidityFee, asset.String())
	var record uint64
	_, err := k.getUint64(ctx, key, &record)
	return record, err
}

// ResetRollingPoolLiquidityFee set the given pool's rolling liquidity fee to zero
func (k KVStore) ResetRollingPoolLiquidityFee(ctx cosmos.Context, asset common.Asset) {
	key := k.GetKey(prefixRollingPoolLiquidityFee, asset.String())
	k.setUint64(ctx, key, 0)
}

func (k KVStore) getLiquidityFees(ctx cosmos.Context, key string) (cosmos.Uint, error) {
	var record uint64
	_, err := k.getUint64(ctx, key, &record)
	return cosmos.NewUint(record), err
}

// GetTotalLiquidityFees - total of all fees collected in each block
func (k KVStore) GetTotalLiquidityFees(ctx cosmos.Context, height uint64) (cosmos.Uint, error) {
	key := k.GetKey(prefixTotalLiquidityFee, strconv.FormatUint(height, 10))
	return k.getLiquidityFees(ctx, key)
}

// GetPoolLiquidityFees - total of fees collected in each block per pool
func (k KVStore) GetPoolLiquidityFees(ctx cosmos.Context, height uint64, asset common.Asset) (cosmos.Uint, error) {
	key := k.GetKey(prefixPoolLiquidityFee, fmt.Sprintf("%d-%s", height, asset.String()))
	return k.getLiquidityFees(ctx, key)
}
