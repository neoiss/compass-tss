package mapo

import (
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/pkg/errors"
)

// GetLastObservedInHeight returns the lastobservedin value for the chain past in
func (b *Bridge) GetLastObservedInHeight(chain common.Chain) (int64, error) {
	// todo handler
	lastblock, err := b.getLastBlock(chain)
	if err != nil {
		return 0, err
	}
	if lastblock == nil {
		return 0, nil
	}

	return lastblock.LastObservedIn, nil
}

// GetBlockHeight returns the current height for mapBridge blocks
func (b *Bridge) GetBlockHeight() (int64, error) {
	return b.ethRpc.GetBlockHeight()
}

type LastBlock struct {
	Chain          string `json:"chain"`
	LastObservedIn int64  `json:"last_observed_in"`
	Relay          int64  `json:"relay"`
}

func (b *Bridge) getLastBlock(chain common.Chain) (*LastBlock, error) {
	method := constants.GetLastTxOutHeight
	input, err := b.viewAbi.Pack(method)
	if err != nil {
		return nil, errors.Wrap(err, "pack txOut failed")
	}

	var relayHeight *big.Int
	err = b.callContract(&relayHeight, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract")
	}

	method = constants.GetLastTxInHeight
	cId, _ := chain.ChainID()
	input, err = b.viewAbi.Pack(method, cId)
	if err != nil {
		return nil, errors.Wrap(err, "pack txIn failed")
	}

	var otherHeight *big.Int
	err = b.callContract(&otherHeight, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract")
	}

	return &LastBlock{
		Chain:          chain.String(),
		LastObservedIn: otherHeight.Int64(),
		Relay:          relayHeight.Int64(),
	}, nil
}

// GetAsgards retrieve all the vaults from mapBridge
func (b *Bridge) GetAsgards() (shareTypes.Vaults, error) {
	vaults, err := b.getPublickeys()
	if err != nil {
		return nil, err
	}

	ret := make(shareTypes.Vaults, 0, len(vaults))
	for _, ele := range vaults {
		vault, err := b.GetVault(ele.PubKey)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *vault)
	}

	return ret, nil
}

// GetPubKeys retrieve vault pub keys and their relevant smart contracts
func (b *Bridge) GetPubKeys() ([]shareTypes.PubKeyContractAddressPair, error) {
	vaults, err := b.getPublickeys()
	if err != nil {
		return nil, err
	}

	ret := make([]shareTypes.PubKeyContractAddressPair, 0, len(vaults))
	for _, ele := range vaults {
		compressedPubKey, err := common.CompressPubKey(ele.PubKey)
		if err != nil {
			continue
		}
		vault := shareTypes.PubKeyContractAddressPair{
			PubKey:           common.PubKey("04" + ecommon.Bytes2Hex(ele.PubKey)),
			CompressedPubKey: common.PubKey(compressedPubKey),
			Contracts:        make(map[common.Chain]common.Address),
		}
		for _, router := range ele.Routers {
			chain, ok := common.GetChainName(router.Chain)
			if !ok {
				continue
			}
			vault.Contracts[chain] = common.Address(ecommon.BytesToAddress(router.Router).String())
		}
		ret = append(ret, vault)
	}

	return ret, nil
}

// GetAsgardPubKeys retrieve vaults, and it's relevant smart contracts
func (b *Bridge) GetAsgardPubKeys() ([]shareTypes.PubKeyContractAddressPair, error) {
	vaults, err := b.getPublickeys()
	if err != nil {
		return nil, err
	}
	ret := make([]shareTypes.PubKeyContractAddressPair, 0, len(vaults))
	for _, ele := range vaults {
		compressedPubKey, err := common.CompressPubKey(ele.PubKey)
		if err != nil {
			continue
		}
		vault := shareTypes.PubKeyContractAddressPair{
			PubKey:           common.PubKey("04" + ecommon.Bytes2Hex(ele.PubKey)),
			CompressedPubKey: common.PubKey(compressedPubKey),
			Contracts:        make(map[common.Chain]common.Address),
		}
		for _, router := range ele.Routers {
			chain, ok := common.GetChainName(router.Chain)
			if !ok {
				continue
			}
			vault.Contracts[chain] = common.Address(ecommon.BytesToAddress(router.Router).String())
		}

		ret = append(ret, vault)
	}

	return ret, nil
}

type VaultInfo struct {
	PubKey  []byte
	Routers []VaultRouter
}

type VaultRouter struct {
	Chain  *big.Int
	Router []byte
}

func (b *Bridge) getPublickeys() ([]VaultInfo, error) {
	method := constants.GetPublickeys
	input, err := b.viewAbi.Pack(method)
	if err != nil {
		return nil, err
	}

	var ret []VaultInfo
	err = b.callContract(&ret, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (b *Bridge) GetVault(pubkey []byte) (*shareTypes.Vault, error) {
	method := constants.GetVault
	input, err := b.viewAbi.Pack(method, pubkey)
	if err != nil {
		return nil, err
	}

	ret := shareTypes.Vault{}
	err = b.callContract(&ret, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call %s", method)
	}
	return &ret, nil
}
