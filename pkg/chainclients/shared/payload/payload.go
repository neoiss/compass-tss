package payload

import (
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/constants"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/utxo"
	mem "github.com/mapprotocol/compass-tss/x/memo"
)

type Editer struct {
	bridge shareTypes.Bridge
}

func New(bridge shareTypes.Bridge) *Editer {
	return &Editer{
		bridge: bridge,
	}
}

func (e *Editer) Encode(nativeToken, destToken string, mapChainID, destChainID *big.Int, to []byte, parsedMemo mem.Memo) (payload []byte, err error) {
	var relayData []byte
	// if the dest token is not native token, we need to build the relay data
	if !strings.EqualFold(destToken, nativeToken) {
		destTokenAddress, err := e.bridge.GetTokenAddress(mapChainID, destToken)
		if err != nil {
			return nil, fmt.Errorf("fail to get token address: %w, chainID: %s, token: %s", err, mapChainID, destToken)
		}
		decimals, err := e.bridge.GetTokenDecimals(mapChainID, destTokenAddress)
		if err != nil {
			return nil, fmt.Errorf("fail to get token decimals: %w, chainID: %s, token: %s", err, mapChainID, destToken)
		}
		if decimals.Cmp(big.NewInt(0)) == 0 {
			decimals = big.NewInt(constants.DefaultTokenDecimals)
		}
		minAmount := utxo.ConvertDecimal(parsedMemo.GetAmount().BigInt(), 6, decimals.Uint64())
		relayData, err = utxo.EncodeRelayData(ethcommon.BytesToAddress(destTokenAddress), minAmount)
		if err != nil {
			return nil, fmt.Errorf("fail to encode relay data: %w, token: %s, minAmount: %d", err, ethcommon.BytesToAddress(destTokenAddress), minAmount)
		}
	}

	var targetData []byte
	// if the dest chain is not the map chain, we need to build the target data.
	if destChainID.Cmp(mapChainID) != 0 {
		targetData, err = utxo.EncodeTargetData(to, destChainID)
		if err != nil {
			return nil, fmt.Errorf("fail to encode target data: %w, to: %s, chainID: %d", err, parsedMemo.GetDestination(), destChainID)
		}
	}

	affiliateData, err := e.encodeAffiliates(parsedMemo.GetAffiliates())
	if err != nil {
		return nil, fmt.Errorf("fail to encode affiliates: %w", err)
	}

	payload, err = utxo.EncodePayload(affiliateData, relayData, targetData)
	if err != nil {
		return nil, fmt.Errorf("fail to encode payload: %w", err)
	}
	return payload, nil
}

func (e *Editer) encodeAffiliates(affiliates mem.Affiliates) ([]byte, error) {
	as := make([]*utxo.Affiliate, 0, len(affiliates))
	for _, aff := range affiliates {
		if aff.Compressed {
			id, err := e.bridge.GetAffiliateIDByAlias(aff.Name)
			if err != nil {
				return []byte{}, fmt.Errorf("fail to get affiliate id by alias: %w", err)
			}
			as = append(as, &utxo.Affiliate{
				ID:  id,
				Bps: uint16(aff.Bps.Uint64()),
			})
			continue
		}

		id, err := e.bridge.GetAffiliateIDByName(aff.Name)
		if err != nil {
			return []byte{}, fmt.Errorf("fail to get affiliate id by name: %w", err)
		}

		as = append(as, &utxo.Affiliate{
			ID:  id,
			Bps: uint16(aff.Bps.Uint64()),
		})
	}
	return utxo.EncodeAffiliateData(as)
}
