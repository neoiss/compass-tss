package common

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	dogchaincfg "github.com/eager7/dogd/chaincfg"
	"github.com/hashicorp/go-multierror"
	ltcchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/constants"
)

const (
	EmptyChain = Chain("")
	BSCChain   = Chain("BSC")
	ETHChain   = Chain("ETH")
	BTCChain   = Chain("BTC")
	LTCChain   = Chain("LTC")
	BCHChain   = Chain("BCH")
	DOGEChain  = Chain("DOGE")
	THORChain  = Chain("THOR")
	GAIAChain  = Chain("GAIA")
	AVAXChain  = Chain("AVAX")
	BASEChain  = Chain("BASE")
	XRPChain   = Chain("XRP")
	MAPChain   = Chain("MAP")

	SigningAlgoSecp256k1 = SigningAlgo("secp256k1")
	SigningAlgoEd25519   = SigningAlgo("ed25519")
)

var AllChains = [...]Chain{
	BSCChain,
	ETHChain,
	BTCChain,
	LTCChain,
	BCHChain,
	DOGEChain,
	THORChain,
	GAIAChain,
	AVAXChain,
	BASEChain,
	XRPChain,
	MAPChain,
}

var chainToChainID = map[string]*big.Int{
	// test network
	getChainKey(BSCChain, TestNet):  big.NewInt(97),
	getChainKey(ETHChain, TestNet):  big.NewInt(11155111),
	getChainKey(BTCChain, TestNet):  big.NewInt(1360095883558914),
	getChainKey(DOGEChain, TestNet): big.NewInt(1360095883558916),
	getChainKey(AVAXChain, TestNet): big.NewInt(43113),
	getChainKey(BASEChain, TestNet): big.NewInt(84532),
	getChainKey(MAPChain, TestNet):  big.NewInt(212),

	// main network
	getChainKey(BSCChain, MainNet):  big.NewInt(56),
	getChainKey(ETHChain, MainNet):  big.NewInt(1),
	getChainKey(BTCChain, MainNet):  big.NewInt(1360095883558913),
	getChainKey(DOGEChain, MainNet): big.NewInt(1360095883558915),
	getChainKey(AVAXChain, MainNet): big.NewInt(43114),
	getChainKey(BASEChain, MainNet): big.NewInt(8453),
	getChainKey(MAPChain, MainNet):  big.NewInt(22776),
}

var chainIDToChain = map[string]Chain{
	// test network
	big.NewInt(97).String():               BSCChain,
	big.NewInt(11155111).String():         ETHChain,
	big.NewInt(1360095883558914).String(): BTCChain,
	big.NewInt(1360095883558916).String(): DOGEChain,
	big.NewInt(43113).String():            AVAXChain,
	big.NewInt(84532).String():            BASEChain,
	big.NewInt(97).String():               MAPChain,

	// main network
	big.NewInt(56).String():               BSCChain,
	big.NewInt(1).String():                ETHChain,
	big.NewInt(1360095883558913).String(): BTCChain,
	big.NewInt(1360095883558915).String(): DOGEChain,
	big.NewInt(43114).String():            AVAXChain,
	big.NewInt(8453).String():             BASEChain,
	big.NewInt(22776).String():            MAPChain,
}

func GetChainName(key *big.Int) (Chain, bool) {
	if key == nil {
		return EmptyChain, false
	}
	chain, ok := chainIDToChain[key.String()]
	if !ok {
		return EmptyChain, false
	}
	return chain, ok
}

type SigningAlgo string

type Chain string

// Chains represent a slice of Chain
type Chains []Chain

func getChainKey(chain Chain, network ChainNetwork) string {
	return fmt.Sprintf("%v-%v", chain, network)
}

// Valid validates chain format, should consist only of uppercase letters
func (c Chain) Valid() error {
	_, ok := chainToChainID[getChainKey(c, CurrentChainNetwork)]
	if !ok {
		return UnsupportedChain
	}
	return nil
}

// NewChain create a new Chain and default the siging_algo to Secp256k1
func NewChain(chainID string) (Chain, error) {
	chain := Chain(strings.ToUpper(chainID))
	if err := chain.Valid(); err != nil {
		return chain, err
	}
	return chain, nil
}

// Equals compare two chain to see whether they represent the same chain
func (c Chain) Equals(c2 Chain) bool {
	return strings.EqualFold(c.String(), c2.String())
}

func (c Chain) IsTHORChain() bool {
	return c.Equals(THORChain)
}

func (c Chain) IsBSCChain() bool {
	return c.Equals(BSCChain)
}

// IsEVM returns true if given chain is an EVM chain.
// See working definition of an "EVM" chain in the
// `GetEVMChains` function description
func (c Chain) IsEVM() bool {
	evmChains := GetEVMChains()
	for _, evm := range evmChains {
		if c.Equals(evm) {
			return true
		}
	}
	return false
}

// IsUTXO returns true if given chain is a UTXO chain.
func (c Chain) IsUTXO() bool {
	utxoChains := GetUTXOChains()
	for _, utxo := range utxoChains {
		if c.Equals(utxo) {
			return true
		}
	}
	return false
}

// IsEmpty is to determinate whether the chain is empty
func (c Chain) IsEmpty() bool {
	return strings.TrimSpace(c.String()) == ""
}

// String implement fmt.Stringer
func (c Chain) String() string {
	return strings.ToUpper(string(c))
}

func (c Chain) ChainID() (*big.Int, error) {
	id, ok := chainToChainID[getChainKey(c, CurrentChainNetwork)]
	if !ok {
		return nil, UnsupportedChain
	}
	return id, nil
}

// GetSigningAlgo get the signing algorithm for the given chain
func (c Chain) GetSigningAlgo() SigningAlgo {
	// Only SigningAlgoSecp256k1 is supported for now
	return SigningAlgoSecp256k1
}

// GetGasAsset chain's base asset
func (c Chain) GetGasAsset() Asset {
	switch c {
	case THORChain:
		return RuneNative
	case BSCChain:
		return BNBBEP20Asset
	case BTCChain:
		return BTCAsset
	case LTCChain:
		return LTCAsset
	case BCHChain:
		return BCHAsset
	case DOGEChain:
		return DOGEAsset
	case ETHChain:
		return ETHAsset
	case AVAXChain:
		return AVAXAsset
	case GAIAChain:
		return ATOMAsset
	case BASEChain:
		return BaseETHAsset
	case XRPChain:
		return XRPAsset
	default:
		return EmptyAsset
	}
}

// GetGasUnits returns name of the gas unit for each chain
func (c Chain) GetGasUnits() string {
	switch c {
	case AVAXChain:
		return "nAVAX"
	case BTCChain:
		return "satsperbyte"
	case BCHChain:
		return "satsperbyte"
	case DOGEChain:
		return "satsperbyte"
	case ETHChain, BSCChain, BASEChain:
		return "gwei"
	case GAIAChain:
		return "uatom"
	case LTCChain:
		return "satsperbyte"
	case XRPChain:
		return "drop"
	default:
		return ""
	}
}

// GetGasAssetDecimal returns decimals for the gas asset of the given chain. Currently
// Gaia is 1e6 and all others are 1e8. If an external chain's gas asset is larger than
// 1e8, just return cosmos.DefaultCoinDecimals.
func (c Chain) GetGasAssetDecimal() int64 {
	switch c {
	case GAIAChain:
		return 6
	case XRPChain:
		return 6
	default:
		return cosmos.DefaultCoinDecimals
	}
}

// IsValidAddress make sure the address is correct for the chain
// And this also make sure mocknet doesn't use mainnet address vice versa
func (c Chain) IsValidAddress(addr Address) bool {
	network := CurrentChainNetwork
	prefix := c.AddressPrefix(network)
	return strings.HasPrefix(addr.String(), prefix)
}

// AddressPrefix return the address prefix used by the given network (mocknet/mainnet)
func (c Chain) AddressPrefix(cn ChainNetwork) string {
	if c.IsEVM() {
		return "0x"
	}
	switch cn {
	case TestNet:
		switch c {
		case GAIAChain:
			return "cosmos"
		//case THORChain:
		//	// TODO update this to use mocknet address prefix
		//	return types.GetConfig().GetBech32AccountAddrPrefix()
		case BTCChain:
			return chaincfg.TestNet3Params.Bech32HRPSegwit
		case LTCChain:
			return ltcchaincfg.TestNet4Params.Bech32HRPSegwit
		case DOGEChain:
			return dogchaincfg.TestNet3Params.Bech32HRPSegwit
		}
	case MainNet:
		switch c {
		case GAIAChain:
			return "cosmos"
		//case THORChain:
		//	return types.GetConfig().GetBech32AccountAddrPrefix()
		case BTCChain:
			return chaincfg.MainNetParams.Bech32HRPSegwit
		case LTCChain:
			return ltcchaincfg.MainNetParams.Bech32HRPSegwit
		case DOGEChain:
			return dogchaincfg.MainNetParams.Bech32HRPSegwit
		}
	}
	return ""
}

// DustThreshold returns the min dust threshold for each chain
// The min dust threshold defines the lower end of the withdraw range of memoless savers txs
// The native coin value provided in a memoless tx defines a basis points amount of Withdraw or Add to a savers position as follows:
// Withdraw range: (dust_threshold + 1) -> (dust_threshold + 10_000)
// Add range: dust_threshold -> Inf
// NOTE: these should all be in 8 decimal places
func (c Chain) DustThreshold() cosmos.Uint {
	switch c {
	case BTCChain, LTCChain, BCHChain:
		return cosmos.NewUint(10_000)
	case DOGEChain:
		return cosmos.NewUint(100_000_000)
	case ETHChain, AVAXChain, GAIAChain, BSCChain, BASEChain:
		return cosmos.OneUint()
	case XRPChain:
		// XRP's dust threshold is being set to 1 XRP. This is the base reserve requirement on XRP's ledger.
		// It is set to this value for two reasons:
		//    1. to prevent edge cases of outbound XRP to new addresses where this is the minimum that must be transferred
		//    2. to burn this amount on churns of each XRP vault, effectively leaving it behind as it cannot be transferred, but still transferring all other XRP
		// On churns, we can optionally delete the account to recover an additional .8 XRP, but would increases code complexity and will remove related ledger entries
		// Comparing to BTC, this dust threshold should be reasonable.
		return cosmos.NewUint(One) // 1 XRP
	default:
		return cosmos.ZeroUint()
	}
}

// MaxMemoLength returns the max memo length for each chain.
func (c Chain) MaxMemoLength() int {
	switch c {
	case BTCChain, LTCChain, BCHChain, DOGEChain:
		return 80
	default:
		// Default to the max memo size that we will process, regardless
		// of any higher memo size capable on other chains.
		return constants.MaxMemoSize
	}
}

// DefaultCoinbase returns the default coinbase address for each chain, returns 0 if no
// coinbase emission is used. This is used used at the time of writing as a fallback
// value in Bifrost, and for inbound confirmation count estimates in the quote APIs.
func (c Chain) DefaultCoinbase() float64 {
	switch c {
	case BTCChain:
		return 3.125
	case LTCChain:
		return 6.25
	case BCHChain:
		return 3.125
	case DOGEChain:
		return 10000
	default:
		return 0
	}
}

func (c Chain) ApproximateBlockMilliseconds() int64 {
	switch c {
	case BTCChain:
		return 600_000
	case LTCChain:
		return 150_000
	case BCHChain:
		return 600_000
	case DOGEChain:
		return 60_000
	case ETHChain:
		return 12_000
	case AVAXChain:
		return 3_000
	case BSCChain:
		return 3_000
	case GAIAChain:
		return 6_000
	case THORChain:
		return 6_000
	case BASEChain:
		return 2_000
	case MAPChain:
		return 5_000
	case XRPChain:
		return 4_000 // approx 3-5 seconds
	default:
		return 0
	}
}

func (c Chain) InboundNotes() string {
	switch c {
	case BTCChain, LTCChain, BCHChain, DOGEChain:
		return "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats."
	case ETHChain, AVAXChain, BSCChain, BASEChain:
		return "Base Asset: Send the inbound_address the asset with the memo encoded in hex in the data field. Tokens: First approve router to spend tokens from user: asset.approve(router, amount). Then call router.depositWithExpiry(inbound_address, asset, amount, memo, expiry). Asset is the token contract address. Amount should be in native asset decimals (eg 1e18 for most tokens). Do not swap to smart contract addresses."
	case GAIAChain:
		return "Transfer the inbound_address the asset with the memo. Do not use multi-in, multi-out transactions."
	case THORChain:
		return "Broadcast a MsgDeposit to the THORChain network with the appropriate memo. Do not use multi-in, multi-out transactions."
	case XRPChain:
		return "Transfer the inbound_address the asset with the memo. Only a single memo is supported and only MemoData is used."
	default:
		return ""
	}
}

// GetEVMChains returns all "EVM" chains connected to THORChain
// "EVM" is defined, in thornode's context, as a chain that:
// - uses 0x as an address prefix
// - has a "Router" Smart Contract
func GetEVMChains() []Chain {
	return []Chain{ETHChain, AVAXChain, BSCChain, BASEChain}
}

// GetUTXOChains returns all "UTXO" chains connected to THORChain.
func GetUTXOChains() []Chain {
	return []Chain{BTCChain, LTCChain, BCHChain, DOGEChain}
}

func NewChains(raw []string) (Chains, error) {
	var returnErr error
	var chains Chains
	for _, c := range raw {
		chain, err := NewChain(c)
		if err == nil {
			chains = append(chains, chain)
		} else {
			returnErr = multierror.Append(returnErr, err)
		}
	}
	return chains, returnErr
}

// Has check whether chain c is in the list
func (chains Chains) Has(c Chain) bool {
	for _, ch := range chains {
		if ch.Equals(c) {
			return true
		}
	}
	return false
}

// Distinct return a distinct set of chains, no duplicates
func (chains Chains) Distinct() Chains {
	var newChains Chains
	for _, chain := range chains {
		if !newChains.Has(chain) {
			newChains = append(newChains, chain)
		}
	}
	return newChains
}

func (chains Chains) Strings() []string {
	strings := make([]string, len(chains))
	for i, c := range chains {
		strings[i] = c.String()
	}
	return strings
}
