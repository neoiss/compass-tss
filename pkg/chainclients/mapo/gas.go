package mapo

import (
	"fmt"
	"math/big"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

// HasNetworkFee checks whether the given chain has set a network fee - determined by
// whether the `outbound_tx_size` for the inbound address response is non-zero.
func (b *Bridge) HasNetworkFee(chain common.Chain) (bool, error) {
	rate, size, sizeWIthCall, err := b.GetNetworkFee(chain)
	if err != nil {
		return false, err
	}
	if rate != 0 && size != 0 && sizeWIthCall != 0 {
		return true, nil
	}

	return false, fmt.Errorf("no inbound address found for chain: %s", chain)
}

// GetNetworkFee get chain's network fee from THORNode.
func (b *Bridge) GetNetworkFee(chain common.Chain) (uint64, uint64, uint64, error) {
	method := constants.GetNetworkFeeInfo
	cId, err := chain.ChainID()
	if err != nil {
		return 0, 0, 0, err
	}
	input, err := b.gasAbi.Pack(method, cId)
	if err != nil {
		return 0, 0, 0, err
	}

	ret := struct {
		TransactionRate         *big.Int
		TransactionSize         *big.Int
		TransactionSizeWithCall *big.Int
	}{}
	err = b.callContract(&ret, b.cfg.GasService, method, input, b.gasAbi)
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "unable to call %s", method)
	}

	return 0, 0, 0, nil
}
