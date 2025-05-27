//go:build mocknet
// +build mocknet

package utxo

import (
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func GetConfMulBasisPoint(chain string, bridge mapo.ThorchainBridge) (cosmos.Uint, error) {
	return cosmos.NewUint(1), nil
}

func MaxConfAdjustment(confirm uint64, chain string, bridge mapo.ThorchainBridge) (uint64, error) {
	return 1, nil
}
