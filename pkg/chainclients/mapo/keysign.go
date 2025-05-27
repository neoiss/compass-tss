package mapo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"math/big"
	"time"

	"github.com/mapprotocol/compass-tss/mapclient/types"
)

var ErrNotFound = fmt.Errorf("not found")

type QueryKeysign struct {
	Keysign   types.TxOut `json:"keysign"`
	Signature string      `json:"signature"`
}

func (b *thorchainBridge) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func (b *thorchainBridge) getFilterLogs(query ethereum.FilterQuery) ([]etypes.Log, error) {
	ctx, cancel := b.getContext()
	defer cancel()
	return b.client.FilterLogs(ctx, query)
}

// GetKeysign retrieves txout from this block height from thorchain
func (b *thorchainBridge) GetKeysign(blockHeight int64, pk string) (types.TxOut, error) {
	logs, err := b.getFilterLogs(ethereum.FilterQuery{
		FromBlock: big.NewInt(blockHeight),
		ToBlock:   big.NewInt(blockHeight),
		Addresses: []ecommon.Address{ecommon.HexToAddress(constants.MosAddressOfMap)},
		Topics:    [][]ecommon.Hash{{ecommon.HexToHash(constants.EventOfMapRelay)}},
	})
	if len(logs) == 0 {
		return types.TxOut{}, err
	}

	ret := types.TxOut{
		Height:  blockHeight,
		TxArray: make([]types.TxArrayItem, 0, len(logs)),
	}
	for _, ele := range logs {
		inHash, _ := common.NewTxID(ele.TxHash.Hex())
		ret.TxArray = append(ret.TxArray, types.TxArrayItem{
			Chain:                 common.MapChain,
			ToAddress:             "",
			VaultPubKey:           "",
			Coin:                  common.Coin{},
			Memo:                  "",
			MaxGas:                nil,
			GasRate:               0,
			InHash:                inHash,
			OutHash:               "",
			Aggregator:            "",
			AggregatorTargetAsset: "",
			AggregatorTargetLimit: nil,
			CloutSpent:            "",
		})
	}

	return ret, nil
}
