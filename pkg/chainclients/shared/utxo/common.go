package utxo

import (
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
)

func GetAsgardAddress(chain common.Chain, bridge shareTypes.Bridge) ([]common.Address, error) {
	vaults, err := bridge.GetAsgardPubKeys()
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}

	newAddresses := make([]common.Address, 0)
	for _, v := range vaults {
		var addr common.Address
		addr, err = v.PubKey.GetAddress(chain)
		if err != nil {
			continue
		}
		newAddresses = append(newAddresses, addr)
	}
	return newAddresses, nil
}
