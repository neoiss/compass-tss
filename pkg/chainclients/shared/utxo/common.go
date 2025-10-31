package utxo

import (
	"encoding/hex"
	"fmt"
	"strings"

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

func GetAsgardPubKeyByAddress(chain common.Chain, bridge shareTypes.Bridge, address common.Address) ([]byte, error) {
	vaults, err := bridge.GetAsgardPubKeys()
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}

	for _, v := range vaults {
		addr, err := v.PubKey.GetAddress(chain)
		if err != nil {
			continue
		}
		if !address.Equals(addr) {
			continue
		}

		pubKey, err := hex.DecodeString(strings.TrimPrefix("04", v.PubKey.String()))
		if err != nil {
			return nil, fmt.Errorf("fail to decode pubkey(%s)", v.PubKey.String())
		}
		return pubKey, nil
	}
	return nil, fmt.Errorf("fail to get asgard pub key by address(%s)", address)
}

func GetAsgardAddress2PubKeyMapped(chain common.Chain, bridge shareTypes.Bridge) (map[common.Address][]byte, error) {
	vaults, err := bridge.GetAsgardPubKeys()
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}

	addr2pub := make(map[common.Address][]byte, 0)
	for _, v := range vaults {
		addr, err := v.PubKey.GetAddress(chain)
		if err != nil {
			continue
		}
		pubKey, err := hex.DecodeString(strings.TrimPrefix("04", v.PubKey.String()))
		if err != nil {
			continue
		}
		addr2pub[addr] = pubKey
	}
	return addr2pub, nil
}
