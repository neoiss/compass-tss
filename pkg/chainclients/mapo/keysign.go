package mapo

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
)

var ErrNotFound = fmt.Errorf("not found")

type QueryKeysign struct {
	Keysign   types.TxOut `json:"keysign"`
	Signature string      `json:"signature"`
}

func (b *Bridge) getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func (b *Bridge) getFilterLogs(query ethereum.FilterQuery) ([]etypes.Log, error) {
	ctx, cancel := b.getContextWithTimeout()
	defer cancel()
	return b.ethClient.FilterLogs(ctx, query)
}

// GetTxByBlockNumber retrieves txout from this block height from mapBridge
func (b *Bridge) GetTxByBlockNumber(blockHeight int64, mos string) (types.TxOut, error) {
	// get block
	if blockHeight%100 == 0 {
		b.logger.Info().Int64("height", blockHeight).Msg("fetching txs for height")
	}

	block, err := b.ethRpc.GetBlock(blockHeight)
	if err != nil {
		return types.TxOut{}, err
	}
	err = b.processBlock(block)
	if err != nil {
		b.logger.Error().Err(err).Int64("height", blockHeight).Msg("Failed to search tx in block")
		return types.TxOut{}, fmt.Errorf("failed to process block: %d, err:%w", blockHeight, err)
	}
	// done
	logs, err := b.getFilterLogs(ethereum.FilterQuery{
		FromBlock: big.NewInt(blockHeight),
		ToBlock:   big.NewInt(blockHeight),
		Addresses: []ecommon.Address{ecommon.HexToAddress(mos)}, // done
		Topics: [][]ecommon.Hash{{
			constants.RelayEventOfMigration.GetTopic(),
			constants.RelayEventOfTransferCall.GetTopic(),
			constants.RelayEventOfTransferOut.GetTopic(),
		}},
	})
	if len(logs) == 0 {
		return types.TxOut{}, err
	}
	b.logger.Info().Msgf("Find tx blockHeight=%v, logs=%d", blockHeight, len(logs))
	b.logger.Info().Msgf("Find tx blockHeight=%v, logs=%d", blockHeight, len(logs))
	b.logger.Info().Msgf("Find tx blockHeight=%v, logs=%d", blockHeight, len(logs))
	time.Sleep(time.Minute)

	ret := types.TxOut{
		Height:  blockHeight,
		TxArray: make([]types.TxArrayItem, 0, len(logs)),
	}

	// todo handler parse coins & gas
	for _, ele := range logs {
		tmp := ele
		item := &types.TxArrayItem{}
		p := evm.NewSmartContractLogParser(nil,
			nil,
			nil,
			nil,
			b.mainAbi,
			common.MAPAsset,
			0)
		p.GetTxOutItem(&tmp, item)
		if item.Chain == nil {
			continue
		}

		ret.TxArray = append(ret.TxArray, *item)
	}

	return ret, nil
}

func (b *Bridge) processBlock(block *etypes.Block) error {
	// collect gas prices of txs in current block
	var txsGas []*big.Int
	for _, tx := range block.Transactions() {
		txsGas = append(txsGas, tx.GasPrice())
	}
	b.updateGasPrice(txsGas)

	return nil
}

// updateGasPrice calculates and stores the current gas price to reported to thornode
func (b *Bridge) updateGasPrice(prices []*big.Int) {
	// skip empty blocks
	if len(prices) == 0 {
		return
	}

	// find the median gas price in the block
	sort.Slice(prices, func(i, j int) bool { return prices[i].Cmp(prices[j]) == -1 })
	gasPrice := prices[len(prices)/2]

	// add to the cache
	b.gasCache = append(b.gasCache, gasPrice)
	if len(b.gasCache) > 20 {
		b.gasCache = b.gasCache[(len(b.gasCache) - 20):]
	}

	// skip update unless cache is full
	if len(b.gasCache) < 20 { // b.cfg.GasCacheBlocks todo handler add cfg
		return
	}

	// compute the median of the median prices in the cache
	medians := []*big.Int{}
	medians = append(medians, b.gasCache...)
	sort.Slice(medians, func(i, j int) bool { return medians[i].Cmp(medians[j]) == -1 })
	median := medians[len(medians)/2]

	// round the price up to nearest configured resolution
	resolution := big.NewInt(100000000) // todo handler add cfg
	median.Add(median, new(big.Int).Sub(resolution, big.NewInt(1)))
	median = median.Div(median, resolution)
	median = median.Mul(median, resolution)
	b.gasPrice = median

	// // record metrics
	// gasPriceFloat, _ := new(big.Float).SetInt64(b.gasPrice.Int64()).Float64()
	// if b.m == nil {
	// 	return
	// }
	// b.m.GetGauge(metrics.GasPrice(b.cfg.ChainID)).Set(gasPriceFloat)
	// b.m.GetCounter(metrics.GasPriceChange(b.cfg.ChainID)).Inc()
}
