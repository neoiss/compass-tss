package mapo

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/constants"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/pkg/errors"
)

// callView executes a read-only contract call with a default RPC timeout.
func (b *Bridge) callView(to ecommon.Address, data []byte) ([]byte, error) {
	ctx, cancel := common.RPCContext()
	defer cancel()
	return b.ethClient.CallContract(ctx, ethereum.CallMsg{
		From: constants.ZeroAddress,
		To:   &to,
		Data: data,
	}, nil)
}

// //////////////////////////////////////////////////////////
// block related
// //////////////////////////////////////////////////////////

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

// //////////////////////////////////////////////////////////
// vaults related
// //////////////////////////////////////////////////////////

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

// //////////////////////////////////////////////////////////
// gas related
// //////////////////////////////////////////////////////////

// GetNetworkFee get chain's network fee from relay.
func (b *Bridge) GetNetworkFee(chain common.Chain) (uint64, uint64, uint64, error) {
	method := constants.GetNetworkFeeInfo
	cId, err := chain.ChainID()
	if err != nil {
		return 0, 0, 0, err
	}
	input, err := b.viewAbi.Pack(method, cId)
	if err != nil {
		return 0, 0, 0, err
	}

	ret := struct {
		TransactionRate         *big.Int
		TransactionSize         *big.Int
		TransactionSizeWithCall *big.Int
	}{}
	err = b.callContract(&ret, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return 0, 0, 0, errors.Wrapf(err, "unable to call %s", method)
	}

	return ret.TransactionSize.Uint64(),
		ret.TransactionSizeWithCall.Uint64(), ret.TransactionRate.Uint64(), nil
}

// HasNetworkFee checks whether the given chain has set a network fee - determined by
// whether the `outbound_tx_size` for the inbound address response is non-zero.
func (b *Bridge) HasNetworkFee(chain common.Chain) (bool, error) {
	size, sizeWIthCall, rate, err := b.GetNetworkFee(chain)
	if err != nil {
		return false, err
	}
	if rate != 0 && size != 0 && sizeWIthCall != 0 {
		return true, nil
	}

	return false, fmt.Errorf("no inbound address found for chain: %s", chain)
}

// //////////////////////////////////////////////////////////
// affiliate related
// //////////////////////////////////////////////////////////

func (b *Bridge) GetAffiliateIDByName(name string) (uint16, error) {
	if name == "" {
		return 0, errors.New("name is empty")
	}
	method := constants.GetInfoByNickname
	input, err := b.viewAbi.Pack(method, name)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to pack affiliate fee input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	affiliate := struct {
		Info struct {
			Id       uint16
			BaseRate uint16
			MaxRate  uint16
			Wallet   ecommon.Address
			Nickname string
		}
	}{}

	if err = outputs.Copy(&affiliate, unpack); err != nil {
		return 0, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	if affiliate.Info.Id == 0 {
		return 0, errors.New("affiliate not found")
	}
	return affiliate.Info.Id, nil
}

func (b *Bridge) GetAffiliateIDByAlias(name string) (uint16, error) {
	if name == "" {
		return 0, errors.New("name is empty")
	}
	method := constants.GetInfoByShortName
	input, err := b.viewAbi.Pack(method, name)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	affiliate := struct {
		Info struct {
			Id       uint16
			BaseRate uint16
			MaxRate  uint16
			Wallet   ecommon.Address
			Nickname string
		}
	}{}

	if err = outputs.Copy(&affiliate, unpack); err != nil {
		return 0, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	if affiliate.Info.Id == 0 {
		return 0, errors.New("affiliate not found")
	}
	return affiliate.Info.Id, nil
}

// //////////////////////////////////////////////////////////
// token related
// //////////////////////////////////////////////////////////

func (b *Bridge) GetChainID(name string) (*big.Int, error) {
	if name == "" {
		return nil, errors.New("chain name is empty")
	}

	method := constants.GetChainByName
	input, err := b.viewAbi.Pack(method, name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	chainID := big.NewInt(0)

	if err = outputs.Copy(&chainID, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return chainID, nil
}

func (b *Bridge) GetChainName(chain *big.Int) (string, error) {
	if chain == nil {
		return "", errors.New("chain is nil")
	}
	method := constants.GetChainName
	input, err := b.viewAbi.Pack(method, chain)
	if err != nil {
		return "", errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return "", errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	var name string
	if err = outputs.Copy(&name, unpack); err != nil {
		return "", errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return name, nil
}

func (b *Bridge) GetTokenAddress(chainID *big.Int, name string) ([]byte, error) {
	if chainID == nil {
		return nil, errors.New("chainID is nil")
	}
	if name == "" {
		return nil, errors.New("token name is empty")
	}
	method := constants.GetTokenAddressByNickname
	input, err := b.viewAbi.Pack(method, chainID, name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	address := make([]byte, 0)
	if err = outputs.Copy(&address, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	if len(address) == 0 {
		return nil, fmt.Errorf("unsupported token(%d:%s)", chainID, name)
	}
	return address, nil
}

func (b *Bridge) GetTokenDecimals(chainID *big.Int, address []byte) (*big.Int, error) {
	if chainID == nil {
		return nil, errors.New("chainID is nil")
	}
	if address == nil {
		return nil, errors.New("token address is nil")
	}
	method := constants.GetTokenDecimals
	input, err := b.viewAbi.Pack(method, chainID, address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to pack input of %s", method)
	}

	to := ecommon.HexToAddress(b.cfg.ViewController)
	output, err := b.callView(to, input)

	outputs := b.viewAbi.Methods[method].Outputs
	unpack, err := outputs.Unpack(output)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unpack output of %s", method)
	}

	decimals := big.NewInt(0)
	if err = outputs.Copy(&decimals, unpack); err != nil {
		return nil, errors.Wrapf(err, "unable to copy output of %s", method)
	}
	return decimals, nil
}
