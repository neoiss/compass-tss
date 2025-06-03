package mapo

import (
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/x/types"
)

// GetNodeAccount retrieves node account for this address from mapBridge
func (b *Bridge) GetNodeAccount(thorAddr string) (*types.NodeAccount, error) {
	//path := fmt.Sprintf("%s/%s", NodeAccountEndpoint, thorAddr)
	//body, _, err := b.getWithPath(path)
	//if err != nil {
	//	return &types.NodeAccount{}, fmt.Errorf("failed to get node account: %w", err)
	//}
	//var na types.NodeAccount
	//if err = json.Unmarshal(body, &na); err != nil {
	//	return &types.NodeAccount{}, fmt.Errorf("failed to unmarshal node account: %w", err)
	//}
	//return &na, nil
	addr, err := github_com_cosmos_cosmos_sdk_types.AccAddressFromHexUnsafe("C4AB011B5ad82074751477ce7F886869da4F5E80")
	if err != nil {
		return nil, err
	}
	return &types.NodeAccount{
		NodeAddress: addr,
		Status:      types.NodeStatus_Active,
		PubKeySet: common.PubKeySet{
			Secp256k1: "041fae8718efebef904d43843057834c76fe5db9cbdfe65517db952706e5583eb45ef6ac13c7067c2ab105f2e107b6526aacd2fed1f304a163a31e43e5471768f9",
		},
	}, nil
}

// GetNodeAccounts retrieves all node accounts from mapBridge
func (b *Bridge) GetNodeAccounts() ([]*types.NodeAccount, error) {
	//path := NodeAccountsEndpoint
	//body, _, err := b.getWithPath(path)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get node account: %w", err)
	//}
	//var na []*types.NodeAccount
	//if err = json.Unmarshal(body, &na); err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal node accounts: %w", err)
	//}
	//return na, nil
	addr, err := github_com_cosmos_cosmos_sdk_types.AccAddressFromHexUnsafe("C4AB011B5ad82074751477ce7F886869da4F5E80")
	if err != nil {
		return nil, err
	}
	return []*types.NodeAccount{
		{
			NodeAddress: addr,
			Status:      types.NodeStatus_Active,
			PubKeySet: common.PubKeySet{
				Secp256k1: "041fae8718efebef904d43843057834c76fe5db9cbdfe65517db952706e5583eb45ef6ac13c7067c2ab105f2e107b6526aacd2fed1f304a163a31e43e5471768f9",
			},
		},
	}, nil
}
