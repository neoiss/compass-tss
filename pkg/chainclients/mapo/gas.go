package mapo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/structure"
)

// HasNetworkFee checks whether the given chain has set a network fee - determined by
// whether the `outbound_tx_size` for the inbound address response is non-zero.
func (b *Bridge) HasNetworkFee(chain common.Chain) (bool, error) {
	var gasList []structure.Gas
	err := b.viewCall.Call(constants.GetNetworkFee, &gasList, 0)
	if err != nil {
		fmt.Println("err ------", err)
		return false, err
	}
	if len(gasList) == 0 {
		return false, nil
	}

	cId, _ := chain.ChainID()
	for _, gas := range gasList {
		if gas.Chain != nil && gas.Chain.Int64() == cId.Int64() && gas.GasRate != nil && gas.TxSize != nil &&
			gas.TxSizeWithCall != nil {
			return true, nil
		}
	}

	return false, fmt.Errorf("no inbound address found for chain: %s", chain)
}

// GetNetworkFee get chain's network fee from THORNode.
func (b *Bridge) GetNetworkFee(chain common.Chain) (uint64, uint64, uint64, error) {
	var gasList []structure.Gas
	err := b.viewCall.Call(constants.GetNetworkFee, &gasList, 0)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(gasList) == 0 {
		return 0, 0, 0, nil
	}

	cId, _ := chain.ChainID()
	for _, gas := range gasList {
		if gas.Chain.Int64() != cId.Int64() {
			continue
		}

		return gas.TxSize.Uint64(), gas.TxSizeWithCall.Uint64(), gas.GasRate.Uint64(), nil
	}

	return 0, 0, 0, nil
}

// PostNetworkFee send network fee message to MAP
func (b *Bridge) PostNetworkFee(ctx context.Context, height int64, chainId *big.Int, transactionSize,
	transactionSizeWithCall, transactionRate uint64) (string, error) {
	// done next 1
	input, err := b.mainAbi.Pack(constants.VoteNetworkFee, b.epoch, chainId, big.NewInt(height),
		big.NewInt(0).SetUint64(transactionRate),
		big.NewInt(0).SetUint64(transactionSize),
		big.NewInt(0).SetUint64(transactionSizeWithCall))
	if err != nil {
		return "", fmt.Errorf("fail to pack input: %w", err)
	}

	tx, err := b.assemblyTx(ctx, input, 2000000)
	if err != nil {
		return "", fmt.Errorf("fail to assembly tx: %w", err)
	}

	return b.Broadcast(tx)
}
