//go:build !mocknet
// +build !mocknet

package utxo

import (
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/constants"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
)

func GetConfMulBasisPoint(chain string, bridge shareTypes.Bridge) (cosmos.Uint, error) {
	confMultiplier, err := bridge.GetMimirWithRef(constants.MimirTemplateConfMultiplierBasisPoints, chain)
	// should never be negative
	if err != nil || confMultiplier <= 0 {
		return cosmos.NewUint(constants.MaxBasisPts), err
	}
	return cosmos.NewUint(uint64(confMultiplier)), nil
}

func MaxConfAdjustment(confirm uint64, chain string, bridge shareTypes.Bridge) (uint64, error) {
	maxConfirmations, err := bridge.GetMimirWithRef(constants.MimirTemplateMaxConfirmations, chain)
	if err != nil {
		return confirm, err
	}
	if maxConfirmations > 0 && confirm > uint64(maxConfirmations) {
		confirm = uint64(maxConfirmations)
	}
	return confirm, nil
}
