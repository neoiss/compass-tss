package mapo

import (
	"fmt"

	"github.com/mapprotocol/compass-tss/common"
	openapi "github.com/mapprotocol/compass-tss/openapi/gen"
)

// GetLastObservedInHeight returns the lastobservedin value for the chain past in
func (b *Bridge) GetLastObservedInHeight(chain common.Chain) (int64, error) {
	// todo handler
	lastblock, err := b.getLastBlock(chain)
	if err != nil {
		return 0, fmt.Errorf("failed to GetLastObservedInHeight: %w", err)
	}
	for _, item := range lastblock {
		if item.Chain == chain.String() {
			return item.LastObservedIn, nil
		}
	}
	return 0, fmt.Errorf("fail to GetLastObservedInHeight,chain(%s)", chain)
}

// GetLastSignedOutHeight returns the lastsignedout value for mapBridge
func (b *Bridge) GetLastSignedOutHeight(chain common.Chain) (int64, error) {
	// todo handler
	lastblock, err := b.getLastBlock(chain)
	if err != nil {
		return 0, fmt.Errorf("failed to GetLastSignedOutHeight: %w", err)
	}
	for _, item := range lastblock {
		if item.Chain == chain.String() {
			return item.LastSignedOut, nil
		}
	}
	return 0, fmt.Errorf("fail to GetLastSignedOutHeight,chain(%s)", chain)
}

// GetBlockHeight returns the current height for mapBridge blocks
func (b *Bridge) GetBlockHeight() (int64, error) {
	// done
	return b.ethRpc.GetBlockHeight()
}

// getLastBlock calls the /lastblock/{chain} endpoint and Unmarshal's into the QueryResLastBlockHeights type
func (b *Bridge) getLastBlock(chain common.Chain) ([]openapi.LastBlock, error) {
	//path := LastBlockEndpoint
	//if !chain.IsEmpty() {
	//	path = fmt.Sprintf("%s/%s", path, chain.String())
	//}
	//buf, _, err := b.getWithPath(path)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get lastblock: %w", err)
	//}
	//var lastBlock []openapi.LastBlock
	//if err = json.Unmarshal(buf, &lastBlock); err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal last block: %w", err)
	//}
	// todo handler
	return []openapi.LastBlock{
		{
			Chain:          "ETH",
			LastObservedIn: 22622246,
			LastSignedOut:  22622246,
			Thorchain:      22622246,
		},
		{
			Chain:          "BASE",
			LastObservedIn: 31070793,
			LastSignedOut:  31070793,
			Thorchain:      31070793,
		},
		{
			Chain:          "MAP",
			LastObservedIn: 18229656,
			LastSignedOut:  18229656,
			Thorchain:      18229656,
		},
	}, nil
}

func (b *Bridge) GetMapBlockHeight() (int64, error) {
	return b.ethRpc.GetBlockHeight()
}
