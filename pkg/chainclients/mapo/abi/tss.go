// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relay

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BridgeItem is an auto generated low-level Go binding around an user-defined struct.
type BridgeItem struct {
	ChainAndGasLimit *big.Int
	Vault            []byte
	TxType           uint8
	Sequence         *big.Int
	Token            []byte
	Amount           *big.Int
	From             []byte
	To               []byte
	Payload          []byte
}

// ITSSManagerKeyShare is an auto generated low-level Go binding around an user-defined struct.
type ITSSManagerKeyShare struct {
	Pubkey   []byte
	KeyShare []byte
}

// TSSManagerTssPoolParam is an auto generated low-level Go binding around an user-defined struct.
type TSSManagerTssPoolParam struct {
	Epoch     *big.Int
	Pubkey    []byte
	KeyShare  []byte
	Members   []common.Address
	Blames    []common.Address
	Signature []byte
}

// TxInItem is an auto generated low-level Go binding around an user-defined struct.
type TxInItem struct {
	OrderId    [32]byte
	BridgeItem BridgeItem
	Height     *big.Int
	RefundAddr []byte
}

// TxOutItem is an auto generated low-level Go binding around an user-defined struct.
type TxOutItem struct {
	OrderId    [32]byte
	BridgeItem BridgeItem
	Height     *big.Int
	GasUsed    *big.Int
	Sender     common.Address
}

// RelayMetaData contains all meta data concerning the Relay contract.
var RelayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AccessManagedInvalidAuthority\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"delay\",\"type\":\"uint32\"}],\"name\":\"AccessManagedRequiredDelay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AccessManagedUnauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"already_propose\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_members\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_sig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_status\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"no_access\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochId\",\"type\":\"uint256\"}],\"name\":\"MigrateCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"m\",\"type\":\"address\"}],\"name\":\"ResetSlashPoint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"retireEpochId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"activeEpochId\",\"type\":\"uint256\"}],\"name\":\"Retire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"currentId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nextId\",\"type\":\"uint256\"}],\"name\":\"Rotate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_maintainer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_relay\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_parameter\",\"type\":\"address\"}],\"name\":\"Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maitainer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"keyShare\",\"type\":\"bytes\"}],\"name\":\"UpdateKeyShare\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"VoteNetworkFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"height\",\"type\":\"uint128\"},{\"internalType\":\"bytes\",\"name\":\"refundAddr\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structTxInItem\",\"name\":\"txInItem\",\"type\":\"tuple\"}],\"name\":\"VoteTxIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"height\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"gasUsed\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"indexed\":false,\"internalType\":\"structTxOutItem\",\"name\":\"txOutItem\",\"type\":\"tuple\"}],\"name\":\"VoteTxOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyShare\",\"type\":\"bytes\"},{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"blames\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structTSSManager.TssPoolParam\",\"name\":\"param\",\"type\":\"tuple\"}],\"name\":\"VoteUpdateTssPool\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"pubkeyHash\",\"type\":\"bytes32\"}],\"name\":\"_publicKeyHashToAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"ms\",\"type\":\"address[]\"}],\"name\":\"batchGetSlashPoint\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"points\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_epochId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_maintainers\",\"type\":\"address[]\"}],\"name\":\"elect\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"m\",\"type\":\"address\"}],\"name\":\"getKeyShare\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyShare\",\"type\":\"bytes\"}],\"internalType\":\"structITSSManager.KeyShare\",\"name\":\"keyShare\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"}],\"name\":\"getMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"m\",\"type\":\"address\"}],\"name\":\"getSlashPoint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"point\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochId\",\"type\":\"uint256\"}],\"name\":\"getTSSStatus\",\"outputs\":[{\"internalType\":\"enumITSSManager.TSSStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_defaultAdmin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isConsumingScheduledOp\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maintainerManager\",\"outputs\":[{\"internalType\":\"contractIMaintainers\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"parameters\",\"outputs\":[{\"internalType\":\"contractIParameters\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relay\",\"outputs\":[{\"internalType\":\"contractIRelay\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"retireEpochId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activeEpochId\",\"type\":\"uint256\"}],\"name\":\"retire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"currentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nextId\",\"type\":\"uint256\"}],\"name\":\"rotate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_maintainer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_relay\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_parameter\",\"type\":\"address\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAuthority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionSizeWithCall\",\"type\":\"uint256\"}],\"name\":\"voteNetworkFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"height\",\"type\":\"uint128\"},{\"internalType\":\"bytes\",\"name\":\"refundAddr\",\"type\":\"bytes\"}],\"internalType\":\"structTxInItem[]\",\"name\":\"txInItems\",\"type\":\"tuple[]\"}],\"name\":\"voteTxIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"height\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"gasUsed\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTxOutItem[]\",\"name\":\"txOutItems\",\"type\":\"tuple[]\"}],\"name\":\"voteTxOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyShare\",\"type\":\"bytes\"},{\"internalType\":\"address[]\",\"name\":\"members\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"blames\",\"type\":\"address[]\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structTSSManager.TssPoolParam\",\"name\":\"param\",\"type\":\"tuple\"}],\"name\":\"voteUpdateTssPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RelayABI is the input ABI used to generate the binding from.
// Deprecated: Use RelayMetaData.ABI instead.
var RelayABI = RelayMetaData.ABI

// Relay is an auto generated Go binding around an Ethereum contract.
type Relay struct {
	RelayCaller     // Read-only binding to the contract
	RelayTransactor // Write-only binding to the contract
	RelayFilterer   // Log filterer for contract events
}

// RelayCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelaySession struct {
	Contract     *Relay            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayCallerSession struct {
	Contract *RelayCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RelayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayTransactorSession struct {
	Contract     *RelayTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayRaw struct {
	Contract *Relay // Generic contract binding to access the raw methods on
}

// RelayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayCallerRaw struct {
	Contract *RelayCaller // Generic read-only contract binding to access the raw methods on
}

// RelayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayTransactorRaw struct {
	Contract *RelayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelay creates a new instance of Relay, bound to a specific deployed contract.
func NewRelay(address common.Address, backend bind.ContractBackend) (*Relay, error) {
	contract, err := bindRelay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Relay{RelayCaller: RelayCaller{contract: contract}, RelayTransactor: RelayTransactor{contract: contract}, RelayFilterer: RelayFilterer{contract: contract}}, nil
}

// NewRelayCaller creates a new read-only instance of Relay, bound to a specific deployed contract.
func NewRelayCaller(address common.Address, caller bind.ContractCaller) (*RelayCaller, error) {
	contract, err := bindRelay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayCaller{contract: contract}, nil
}

// NewRelayTransactor creates a new write-only instance of Relay, bound to a specific deployed contract.
func NewRelayTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayTransactor, error) {
	contract, err := bindRelay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayTransactor{contract: contract}, nil
}

// NewRelayFilterer creates a new log filterer instance of Relay, bound to a specific deployed contract.
func NewRelayFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayFilterer, error) {
	contract, err := bindRelay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayFilterer{contract: contract}, nil
}

// bindRelay binds a generic wrapper to an already deployed contract.
func bindRelay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RelayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.RelayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Relay *RelayCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Relay *RelaySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Relay.Contract.UPGRADEINTERFACEVERSION(&_Relay.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Relay *RelayCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Relay.Contract.UPGRADEINTERFACEVERSION(&_Relay.CallOpts)
}

// PublicKeyHashToAddress is a free data retrieval call binding the contract method 0xead276b3.
//
// Solidity: function _publicKeyHashToAddress(bytes32 pubkeyHash) pure returns(address)
func (_Relay *RelayCaller) PublicKeyHashToAddress(opts *bind.CallOpts, pubkeyHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "_publicKeyHashToAddress", pubkeyHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PublicKeyHashToAddress is a free data retrieval call binding the contract method 0xead276b3.
//
// Solidity: function _publicKeyHashToAddress(bytes32 pubkeyHash) pure returns(address)
func (_Relay *RelaySession) PublicKeyHashToAddress(pubkeyHash [32]byte) (common.Address, error) {
	return _Relay.Contract.PublicKeyHashToAddress(&_Relay.CallOpts, pubkeyHash)
}

// PublicKeyHashToAddress is a free data retrieval call binding the contract method 0xead276b3.
//
// Solidity: function _publicKeyHashToAddress(bytes32 pubkeyHash) pure returns(address)
func (_Relay *RelayCallerSession) PublicKeyHashToAddress(pubkeyHash [32]byte) (common.Address, error) {
	return _Relay.Contract.PublicKeyHashToAddress(&_Relay.CallOpts, pubkeyHash)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Relay *RelayCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "authority")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Relay *RelaySession) Authority() (common.Address, error) {
	return _Relay.Contract.Authority(&_Relay.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() view returns(address)
func (_Relay *RelayCallerSession) Authority() (common.Address, error) {
	return _Relay.Contract.Authority(&_Relay.CallOpts)
}

// BatchGetSlashPoint is a free data retrieval call binding the contract method 0x069e79f3.
//
// Solidity: function batchGetSlashPoint(uint256 epoch, address[] ms) view returns(uint256[] points)
func (_Relay *RelayCaller) BatchGetSlashPoint(opts *bind.CallOpts, epoch *big.Int, ms []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "batchGetSlashPoint", epoch, ms)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BatchGetSlashPoint is a free data retrieval call binding the contract method 0x069e79f3.
//
// Solidity: function batchGetSlashPoint(uint256 epoch, address[] ms) view returns(uint256[] points)
func (_Relay *RelaySession) BatchGetSlashPoint(epoch *big.Int, ms []common.Address) ([]*big.Int, error) {
	return _Relay.Contract.BatchGetSlashPoint(&_Relay.CallOpts, epoch, ms)
}

// BatchGetSlashPoint is a free data retrieval call binding the contract method 0x069e79f3.
//
// Solidity: function batchGetSlashPoint(uint256 epoch, address[] ms) view returns(uint256[] points)
func (_Relay *RelayCallerSession) BatchGetSlashPoint(epoch *big.Int, ms []common.Address) ([]*big.Int, error) {
	return _Relay.Contract.BatchGetSlashPoint(&_Relay.CallOpts, epoch, ms)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Relay *RelayCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Relay *RelaySession) CurrentEpoch() (*big.Int, error) {
	return _Relay.Contract.CurrentEpoch(&_Relay.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Relay *RelayCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Relay.Contract.CurrentEpoch(&_Relay.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Relay *RelayCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Relay *RelaySession) GetImplementation() (common.Address, error) {
	return _Relay.Contract.GetImplementation(&_Relay.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Relay *RelayCallerSession) GetImplementation() (common.Address, error) {
	return _Relay.Contract.GetImplementation(&_Relay.CallOpts)
}

// GetKeyShare is a free data retrieval call binding the contract method 0x1d22fd76.
//
// Solidity: function getKeyShare(address m) view returns((bytes,bytes) keyShare)
func (_Relay *RelayCaller) GetKeyShare(opts *bind.CallOpts, m common.Address) (ITSSManagerKeyShare, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getKeyShare", m)

	if err != nil {
		return *new(ITSSManagerKeyShare), err
	}

	out0 := *abi.ConvertType(out[0], new(ITSSManagerKeyShare)).(*ITSSManagerKeyShare)

	return out0, err

}

// GetKeyShare is a free data retrieval call binding the contract method 0x1d22fd76.
//
// Solidity: function getKeyShare(address m) view returns((bytes,bytes) keyShare)
func (_Relay *RelaySession) GetKeyShare(m common.Address) (ITSSManagerKeyShare, error) {
	return _Relay.Contract.GetKeyShare(&_Relay.CallOpts, m)
}

// GetKeyShare is a free data retrieval call binding the contract method 0x1d22fd76.
//
// Solidity: function getKeyShare(address m) view returns((bytes,bytes) keyShare)
func (_Relay *RelayCallerSession) GetKeyShare(m common.Address) (ITSSManagerKeyShare, error) {
	return _Relay.Contract.GetKeyShare(&_Relay.CallOpts, m)
}

// GetMembers is a free data retrieval call binding the contract method 0xc680119d.
//
// Solidity: function getMembers(bytes pubkey) view returns(address[] members)
func (_Relay *RelayCaller) GetMembers(opts *bind.CallOpts, pubkey []byte) ([]common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getMembers", pubkey)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetMembers is a free data retrieval call binding the contract method 0xc680119d.
//
// Solidity: function getMembers(bytes pubkey) view returns(address[] members)
func (_Relay *RelaySession) GetMembers(pubkey []byte) ([]common.Address, error) {
	return _Relay.Contract.GetMembers(&_Relay.CallOpts, pubkey)
}

// GetMembers is a free data retrieval call binding the contract method 0xc680119d.
//
// Solidity: function getMembers(bytes pubkey) view returns(address[] members)
func (_Relay *RelayCallerSession) GetMembers(pubkey []byte) ([]common.Address, error) {
	return _Relay.Contract.GetMembers(&_Relay.CallOpts, pubkey)
}

// GetSlashPoint is a free data retrieval call binding the contract method 0x9465708b.
//
// Solidity: function getSlashPoint(uint256 epoch, address m) view returns(uint256 point)
func (_Relay *RelayCaller) GetSlashPoint(opts *bind.CallOpts, epoch *big.Int, m common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getSlashPoint", epoch, m)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlashPoint is a free data retrieval call binding the contract method 0x9465708b.
//
// Solidity: function getSlashPoint(uint256 epoch, address m) view returns(uint256 point)
func (_Relay *RelaySession) GetSlashPoint(epoch *big.Int, m common.Address) (*big.Int, error) {
	return _Relay.Contract.GetSlashPoint(&_Relay.CallOpts, epoch, m)
}

// GetSlashPoint is a free data retrieval call binding the contract method 0x9465708b.
//
// Solidity: function getSlashPoint(uint256 epoch, address m) view returns(uint256 point)
func (_Relay *RelayCallerSession) GetSlashPoint(epoch *big.Int, m common.Address) (*big.Int, error) {
	return _Relay.Contract.GetSlashPoint(&_Relay.CallOpts, epoch, m)
}

// GetTSSStatus is a free data retrieval call binding the contract method 0x128e9f9f.
//
// Solidity: function getTSSStatus(uint256 epochId) view returns(uint8)
func (_Relay *RelayCaller) GetTSSStatus(opts *bind.CallOpts, epochId *big.Int) (uint8, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getTSSStatus", epochId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetTSSStatus is a free data retrieval call binding the contract method 0x128e9f9f.
//
// Solidity: function getTSSStatus(uint256 epochId) view returns(uint8)
func (_Relay *RelaySession) GetTSSStatus(epochId *big.Int) (uint8, error) {
	return _Relay.Contract.GetTSSStatus(&_Relay.CallOpts, epochId)
}

// GetTSSStatus is a free data retrieval call binding the contract method 0x128e9f9f.
//
// Solidity: function getTSSStatus(uint256 epochId) view returns(uint8)
func (_Relay *RelayCallerSession) GetTSSStatus(epochId *big.Int) (uint8, error) {
	return _Relay.Contract.GetTSSStatus(&_Relay.CallOpts, epochId)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_Relay *RelayCaller) IsConsumingScheduledOp(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "isConsumingScheduledOp")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_Relay *RelaySession) IsConsumingScheduledOp() ([4]byte, error) {
	return _Relay.Contract.IsConsumingScheduledOp(&_Relay.CallOpts)
}

// IsConsumingScheduledOp is a free data retrieval call binding the contract method 0x8fb36037.
//
// Solidity: function isConsumingScheduledOp() view returns(bytes4)
func (_Relay *RelayCallerSession) IsConsumingScheduledOp() ([4]byte, error) {
	return _Relay.Contract.IsConsumingScheduledOp(&_Relay.CallOpts)
}

// MaintainerManager is a free data retrieval call binding the contract method 0xaeb541b4.
//
// Solidity: function maintainerManager() view returns(address)
func (_Relay *RelayCaller) MaintainerManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "maintainerManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MaintainerManager is a free data retrieval call binding the contract method 0xaeb541b4.
//
// Solidity: function maintainerManager() view returns(address)
func (_Relay *RelaySession) MaintainerManager() (common.Address, error) {
	return _Relay.Contract.MaintainerManager(&_Relay.CallOpts)
}

// MaintainerManager is a free data retrieval call binding the contract method 0xaeb541b4.
//
// Solidity: function maintainerManager() view returns(address)
func (_Relay *RelayCallerSession) MaintainerManager() (common.Address, error) {
	return _Relay.Contract.MaintainerManager(&_Relay.CallOpts)
}

// Parameters is a free data retrieval call binding the contract method 0x89035730.
//
// Solidity: function parameters() view returns(address)
func (_Relay *RelayCaller) Parameters(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "parameters")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Parameters is a free data retrieval call binding the contract method 0x89035730.
//
// Solidity: function parameters() view returns(address)
func (_Relay *RelaySession) Parameters() (common.Address, error) {
	return _Relay.Contract.Parameters(&_Relay.CallOpts)
}

// Parameters is a free data retrieval call binding the contract method 0x89035730.
//
// Solidity: function parameters() view returns(address)
func (_Relay *RelayCallerSession) Parameters() (common.Address, error) {
	return _Relay.Contract.Parameters(&_Relay.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Relay *RelayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Relay *RelaySession) Paused() (bool, error) {
	return _Relay.Contract.Paused(&_Relay.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Relay *RelayCallerSession) Paused() (bool, error) {
	return _Relay.Contract.Paused(&_Relay.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Relay *RelayCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Relay *RelaySession) ProxiableUUID() ([32]byte, error) {
	return _Relay.Contract.ProxiableUUID(&_Relay.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Relay *RelayCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Relay.Contract.ProxiableUUID(&_Relay.CallOpts)
}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_Relay *RelayCaller) Relay(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "relay")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_Relay *RelaySession) Relay() (common.Address, error) {
	return _Relay.Contract.Relay(&_Relay.CallOpts)
}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_Relay *RelayCallerSession) Relay() (common.Address, error) {
	return _Relay.Contract.Relay(&_Relay.CallOpts)
}

// Elect is a paid mutator transaction binding the contract method 0xfae34a46.
//
// Solidity: function elect(uint256 _epochId, address[] _maintainers) returns(bool)
func (_Relay *RelayTransactor) Elect(opts *bind.TransactOpts, _epochId *big.Int, _maintainers []common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "elect", _epochId, _maintainers)
}

// Elect is a paid mutator transaction binding the contract method 0xfae34a46.
//
// Solidity: function elect(uint256 _epochId, address[] _maintainers) returns(bool)
func (_Relay *RelaySession) Elect(_epochId *big.Int, _maintainers []common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Elect(&_Relay.TransactOpts, _epochId, _maintainers)
}

// Elect is a paid mutator transaction binding the contract method 0xfae34a46.
//
// Solidity: function elect(uint256 _epochId, address[] _maintainers) returns(bool)
func (_Relay *RelayTransactorSession) Elect(_epochId *big.Int, _maintainers []common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Elect(&_Relay.TransactOpts, _epochId, _maintainers)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _defaultAdmin) returns()
func (_Relay *RelayTransactor) Initialize(opts *bind.TransactOpts, _defaultAdmin common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "initialize", _defaultAdmin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _defaultAdmin) returns()
func (_Relay *RelaySession) Initialize(_defaultAdmin common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Initialize(&_Relay.TransactOpts, _defaultAdmin)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _defaultAdmin) returns()
func (_Relay *RelayTransactorSession) Initialize(_defaultAdmin common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Initialize(&_Relay.TransactOpts, _defaultAdmin)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_Relay *RelayTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_Relay *RelaySession) Migrate() (*types.Transaction, error) {
	return _Relay.Contract.Migrate(&_Relay.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns()
func (_Relay *RelayTransactorSession) Migrate() (*types.Transaction, error) {
	return _Relay.Contract.Migrate(&_Relay.TransactOpts)
}

// Retire is a paid mutator transaction binding the contract method 0x293a7f2e.
//
// Solidity: function retire(uint256 retireEpochId, uint256 activeEpochId) returns()
func (_Relay *RelayTransactor) Retire(opts *bind.TransactOpts, retireEpochId *big.Int, activeEpochId *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "retire", retireEpochId, activeEpochId)
}

// Retire is a paid mutator transaction binding the contract method 0x293a7f2e.
//
// Solidity: function retire(uint256 retireEpochId, uint256 activeEpochId) returns()
func (_Relay *RelaySession) Retire(retireEpochId *big.Int, activeEpochId *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Retire(&_Relay.TransactOpts, retireEpochId, activeEpochId)
}

// Retire is a paid mutator transaction binding the contract method 0x293a7f2e.
//
// Solidity: function retire(uint256 retireEpochId, uint256 activeEpochId) returns()
func (_Relay *RelayTransactorSession) Retire(retireEpochId *big.Int, activeEpochId *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Retire(&_Relay.TransactOpts, retireEpochId, activeEpochId)
}

// Rotate is a paid mutator transaction binding the contract method 0xb023e78d.
//
// Solidity: function rotate(uint256 currentId, uint256 nextId) returns()
func (_Relay *RelayTransactor) Rotate(opts *bind.TransactOpts, currentId *big.Int, nextId *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "rotate", currentId, nextId)
}

// Rotate is a paid mutator transaction binding the contract method 0xb023e78d.
//
// Solidity: function rotate(uint256 currentId, uint256 nextId) returns()
func (_Relay *RelaySession) Rotate(currentId *big.Int, nextId *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Rotate(&_Relay.TransactOpts, currentId, nextId)
}

// Rotate is a paid mutator transaction binding the contract method 0xb023e78d.
//
// Solidity: function rotate(uint256 currentId, uint256 nextId) returns()
func (_Relay *RelayTransactorSession) Rotate(currentId *big.Int, nextId *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Rotate(&_Relay.TransactOpts, currentId, nextId)
}

// Set is a paid mutator transaction binding the contract method 0x83209259.
//
// Solidity: function set(address _maintainer, address _relay, address _parameter) returns()
func (_Relay *RelayTransactor) Set(opts *bind.TransactOpts, _maintainer common.Address, _relay common.Address, _parameter common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "set", _maintainer, _relay, _parameter)
}

// Set is a paid mutator transaction binding the contract method 0x83209259.
//
// Solidity: function set(address _maintainer, address _relay, address _parameter) returns()
func (_Relay *RelaySession) Set(_maintainer common.Address, _relay common.Address, _parameter common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Set(&_Relay.TransactOpts, _maintainer, _relay, _parameter)
}

// Set is a paid mutator transaction binding the contract method 0x83209259.
//
// Solidity: function set(address _maintainer, address _relay, address _parameter) returns()
func (_Relay *RelayTransactorSession) Set(_maintainer common.Address, _relay common.Address, _parameter common.Address) (*types.Transaction, error) {
	return _Relay.Contract.Set(&_Relay.TransactOpts, _maintainer, _relay, _parameter)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_Relay *RelayTransactor) SetAuthority(opts *bind.TransactOpts, newAuthority common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setAuthority", newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_Relay *RelaySession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetAuthority(&_Relay.TransactOpts, newAuthority)
}

// SetAuthority is a paid mutator transaction binding the contract method 0x7a9e5e4b.
//
// Solidity: function setAuthority(address newAuthority) returns()
func (_Relay *RelayTransactorSession) SetAuthority(newAuthority common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetAuthority(&_Relay.TransactOpts, newAuthority)
}

// Trigger is a paid mutator transaction binding the contract method 0x7fec8d38.
//
// Solidity: function trigger() returns()
func (_Relay *RelayTransactor) Trigger(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "trigger")
}

// Trigger is a paid mutator transaction binding the contract method 0x7fec8d38.
//
// Solidity: function trigger() returns()
func (_Relay *RelaySession) Trigger() (*types.Transaction, error) {
	return _Relay.Contract.Trigger(&_Relay.TransactOpts)
}

// Trigger is a paid mutator transaction binding the contract method 0x7fec8d38.
//
// Solidity: function trigger() returns()
func (_Relay *RelayTransactorSession) Trigger() (*types.Transaction, error) {
	return _Relay.Contract.Trigger(&_Relay.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Relay *RelayTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Relay *RelaySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Relay.Contract.UpgradeToAndCall(&_Relay.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Relay *RelayTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Relay.Contract.UpgradeToAndCall(&_Relay.TransactOpts, newImplementation, data)
}

// VoteNetworkFee is a paid mutator transaction binding the contract method 0xf79bc017.
//
// Solidity: function voteNetworkFee(uint256 chain, uint256 height, uint256 transactionRate, uint256 transactionSize, uint256 transactionSizeWithCall) returns()
func (_Relay *RelayTransactor) VoteNetworkFee(opts *bind.TransactOpts, chain *big.Int, height *big.Int, transactionRate *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "voteNetworkFee", chain, height, transactionRate, transactionSize, transactionSizeWithCall)
}

// VoteNetworkFee is a paid mutator transaction binding the contract method 0xf79bc017.
//
// Solidity: function voteNetworkFee(uint256 chain, uint256 height, uint256 transactionRate, uint256 transactionSize, uint256 transactionSizeWithCall) returns()
func (_Relay *RelaySession) VoteNetworkFee(chain *big.Int, height *big.Int, transactionRate *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.VoteNetworkFee(&_Relay.TransactOpts, chain, height, transactionRate, transactionSize, transactionSizeWithCall)
}

// VoteNetworkFee is a paid mutator transaction binding the contract method 0xf79bc017.
//
// Solidity: function voteNetworkFee(uint256 chain, uint256 height, uint256 transactionRate, uint256 transactionSize, uint256 transactionSizeWithCall) returns()
func (_Relay *RelayTransactorSession) VoteNetworkFee(chain *big.Int, height *big.Int, transactionRate *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.VoteNetworkFee(&_Relay.TransactOpts, chain, height, transactionRate, transactionSize, transactionSizeWithCall)
}

// VoteTxIn is a paid mutator transaction binding the contract method 0xb87e124e.
//
// Solidity: function voteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes)[] txInItems) returns()
func (_Relay *RelayTransactor) VoteTxIn(opts *bind.TransactOpts, txInItems []TxInItem) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "voteTxIn", txInItems)
}

// VoteTxIn is a paid mutator transaction binding the contract method 0xb87e124e.
//
// Solidity: function voteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes)[] txInItems) returns()
func (_Relay *RelaySession) VoteTxIn(txInItems []TxInItem) (*types.Transaction, error) {
	return _Relay.Contract.VoteTxIn(&_Relay.TransactOpts, txInItems)
}

// VoteTxIn is a paid mutator transaction binding the contract method 0xb87e124e.
//
// Solidity: function voteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes)[] txInItems) returns()
func (_Relay *RelayTransactorSession) VoteTxIn(txInItems []TxInItem) (*types.Transaction, error) {
	return _Relay.Contract.VoteTxIn(&_Relay.TransactOpts, txInItems)
}

// VoteTxOut is a paid mutator transaction binding the contract method 0xa23f6a38.
//
// Solidity: function voteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address)[] txOutItems) returns()
func (_Relay *RelayTransactor) VoteTxOut(opts *bind.TransactOpts, txOutItems []TxOutItem) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "voteTxOut", txOutItems)
}

// VoteTxOut is a paid mutator transaction binding the contract method 0xa23f6a38.
//
// Solidity: function voteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address)[] txOutItems) returns()
func (_Relay *RelaySession) VoteTxOut(txOutItems []TxOutItem) (*types.Transaction, error) {
	return _Relay.Contract.VoteTxOut(&_Relay.TransactOpts, txOutItems)
}

// VoteTxOut is a paid mutator transaction binding the contract method 0xa23f6a38.
//
// Solidity: function voteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address)[] txOutItems) returns()
func (_Relay *RelayTransactorSession) VoteTxOut(txOutItems []TxOutItem) (*types.Transaction, error) {
	return _Relay.Contract.VoteTxOut(&_Relay.TransactOpts, txOutItems)
}

// VoteUpdateTssPool is a paid mutator transaction binding the contract method 0x134ec575.
//
// Solidity: function voteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param) returns()
func (_Relay *RelayTransactor) VoteUpdateTssPool(opts *bind.TransactOpts, param TSSManagerTssPoolParam) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "voteUpdateTssPool", param)
}

// VoteUpdateTssPool is a paid mutator transaction binding the contract method 0x134ec575.
//
// Solidity: function voteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param) returns()
func (_Relay *RelaySession) VoteUpdateTssPool(param TSSManagerTssPoolParam) (*types.Transaction, error) {
	return _Relay.Contract.VoteUpdateTssPool(&_Relay.TransactOpts, param)
}

// VoteUpdateTssPool is a paid mutator transaction binding the contract method 0x134ec575.
//
// Solidity: function voteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param) returns()
func (_Relay *RelayTransactorSession) VoteUpdateTssPool(param TSSManagerTssPoolParam) (*types.Transaction, error) {
	return _Relay.Contract.VoteUpdateTssPool(&_Relay.TransactOpts, param)
}

// RelayAuthorityUpdatedIterator is returned from FilterAuthorityUpdated and is used to iterate over the raw logs and unpacked data for AuthorityUpdated events raised by the Relay contract.
type RelayAuthorityUpdatedIterator struct {
	Event *RelayAuthorityUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayAuthorityUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayAuthorityUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayAuthorityUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayAuthorityUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayAuthorityUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayAuthorityUpdated represents a AuthorityUpdated event raised by the Relay contract.
type RelayAuthorityUpdated struct {
	Authority common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuthorityUpdated is a free log retrieval operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_Relay *RelayFilterer) FilterAuthorityUpdated(opts *bind.FilterOpts) (*RelayAuthorityUpdatedIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return &RelayAuthorityUpdatedIterator{contract: _Relay.contract, event: "AuthorityUpdated", logs: logs, sub: sub}, nil
}

// WatchAuthorityUpdated is a free log subscription operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_Relay *RelayFilterer) WatchAuthorityUpdated(opts *bind.WatchOpts, sink chan<- *RelayAuthorityUpdated) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "AuthorityUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayAuthorityUpdated)
				if err := _Relay.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAuthorityUpdated is a log parse operation binding the contract event 0x2f658b440c35314f52658ea8a740e05b284cdc84dc9ae01e891f21b8933e7cad.
//
// Solidity: event AuthorityUpdated(address authority)
func (_Relay *RelayFilterer) ParseAuthorityUpdated(log types.Log) (*RelayAuthorityUpdated, error) {
	event := new(RelayAuthorityUpdated)
	if err := _Relay.contract.UnpackLog(event, "AuthorityUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Relay contract.
type RelayInitializedIterator struct {
	Event *RelayInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayInitialized represents a Initialized event raised by the Relay contract.
type RelayInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Relay *RelayFilterer) FilterInitialized(opts *bind.FilterOpts) (*RelayInitializedIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RelayInitializedIterator{contract: _Relay.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Relay *RelayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RelayInitialized) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayInitialized)
				if err := _Relay.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Relay *RelayFilterer) ParseInitialized(log types.Log) (*RelayInitialized, error) {
	event := new(RelayInitialized)
	if err := _Relay.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayMigrateCompletedIterator is returned from FilterMigrateCompleted and is used to iterate over the raw logs and unpacked data for MigrateCompleted events raised by the Relay contract.
type RelayMigrateCompletedIterator struct {
	Event *RelayMigrateCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayMigrateCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayMigrateCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayMigrateCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayMigrateCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayMigrateCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayMigrateCompleted represents a MigrateCompleted event raised by the Relay contract.
type RelayMigrateCompleted struct {
	EpochId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMigrateCompleted is a free log retrieval operation binding the contract event 0xf26450c98e415ad46e2daa4d4ad23d6b61146c0d033f0f73d87f13a0486e82b0.
//
// Solidity: event MigrateCompleted(uint256 epochId)
func (_Relay *RelayFilterer) FilterMigrateCompleted(opts *bind.FilterOpts) (*RelayMigrateCompletedIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "MigrateCompleted")
	if err != nil {
		return nil, err
	}
	return &RelayMigrateCompletedIterator{contract: _Relay.contract, event: "MigrateCompleted", logs: logs, sub: sub}, nil
}

// WatchMigrateCompleted is a free log subscription operation binding the contract event 0xf26450c98e415ad46e2daa4d4ad23d6b61146c0d033f0f73d87f13a0486e82b0.
//
// Solidity: event MigrateCompleted(uint256 epochId)
func (_Relay *RelayFilterer) WatchMigrateCompleted(opts *bind.WatchOpts, sink chan<- *RelayMigrateCompleted) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "MigrateCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayMigrateCompleted)
				if err := _Relay.contract.UnpackLog(event, "MigrateCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMigrateCompleted is a log parse operation binding the contract event 0xf26450c98e415ad46e2daa4d4ad23d6b61146c0d033f0f73d87f13a0486e82b0.
//
// Solidity: event MigrateCompleted(uint256 epochId)
func (_Relay *RelayFilterer) ParseMigrateCompleted(log types.Log) (*RelayMigrateCompleted, error) {
	event := new(RelayMigrateCompleted)
	if err := _Relay.contract.UnpackLog(event, "MigrateCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Relay contract.
type RelayPausedIterator struct {
	Event *RelayPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayPaused represents a Paused event raised by the Relay contract.
type RelayPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Relay *RelayFilterer) FilterPaused(opts *bind.FilterOpts) (*RelayPausedIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RelayPausedIterator{contract: _Relay.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Relay *RelayFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RelayPaused) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayPaused)
				if err := _Relay.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Relay *RelayFilterer) ParsePaused(log types.Log) (*RelayPaused, error) {
	event := new(RelayPaused)
	if err := _Relay.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayResetSlashPointIterator is returned from FilterResetSlashPoint and is used to iterate over the raw logs and unpacked data for ResetSlashPoint events raised by the Relay contract.
type RelayResetSlashPointIterator struct {
	Event *RelayResetSlashPoint // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayResetSlashPointIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayResetSlashPoint)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayResetSlashPoint)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayResetSlashPointIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayResetSlashPointIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayResetSlashPoint represents a ResetSlashPoint event raised by the Relay contract.
type RelayResetSlashPoint struct {
	M   common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterResetSlashPoint is a free log retrieval operation binding the contract event 0xc3e4735e272b79f3817c816250f740d4e2799a475bc97708012c086f49b8ced0.
//
// Solidity: event ResetSlashPoint(address m)
func (_Relay *RelayFilterer) FilterResetSlashPoint(opts *bind.FilterOpts) (*RelayResetSlashPointIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "ResetSlashPoint")
	if err != nil {
		return nil, err
	}
	return &RelayResetSlashPointIterator{contract: _Relay.contract, event: "ResetSlashPoint", logs: logs, sub: sub}, nil
}

// WatchResetSlashPoint is a free log subscription operation binding the contract event 0xc3e4735e272b79f3817c816250f740d4e2799a475bc97708012c086f49b8ced0.
//
// Solidity: event ResetSlashPoint(address m)
func (_Relay *RelayFilterer) WatchResetSlashPoint(opts *bind.WatchOpts, sink chan<- *RelayResetSlashPoint) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "ResetSlashPoint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayResetSlashPoint)
				if err := _Relay.contract.UnpackLog(event, "ResetSlashPoint", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseResetSlashPoint is a log parse operation binding the contract event 0xc3e4735e272b79f3817c816250f740d4e2799a475bc97708012c086f49b8ced0.
//
// Solidity: event ResetSlashPoint(address m)
func (_Relay *RelayFilterer) ParseResetSlashPoint(log types.Log) (*RelayResetSlashPoint, error) {
	event := new(RelayResetSlashPoint)
	if err := _Relay.contract.UnpackLog(event, "ResetSlashPoint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayRetireIterator is returned from FilterRetire and is used to iterate over the raw logs and unpacked data for Retire events raised by the Relay contract.
type RelayRetireIterator struct {
	Event *RelayRetire // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayRetireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayRetire)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayRetire)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayRetireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayRetireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayRetire represents a Retire event raised by the Relay contract.
type RelayRetire struct {
	RetireEpochId *big.Int
	ActiveEpochId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRetire is a free log retrieval operation binding the contract event 0x439609f0617e4115fc11915679d1b8fa6071e011630e39131744cc297ceef5ef.
//
// Solidity: event Retire(uint256 retireEpochId, uint256 activeEpochId)
func (_Relay *RelayFilterer) FilterRetire(opts *bind.FilterOpts) (*RelayRetireIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Retire")
	if err != nil {
		return nil, err
	}
	return &RelayRetireIterator{contract: _Relay.contract, event: "Retire", logs: logs, sub: sub}, nil
}

// WatchRetire is a free log subscription operation binding the contract event 0x439609f0617e4115fc11915679d1b8fa6071e011630e39131744cc297ceef5ef.
//
// Solidity: event Retire(uint256 retireEpochId, uint256 activeEpochId)
func (_Relay *RelayFilterer) WatchRetire(opts *bind.WatchOpts, sink chan<- *RelayRetire) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Retire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayRetire)
				if err := _Relay.contract.UnpackLog(event, "Retire", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRetire is a log parse operation binding the contract event 0x439609f0617e4115fc11915679d1b8fa6071e011630e39131744cc297ceef5ef.
//
// Solidity: event Retire(uint256 retireEpochId, uint256 activeEpochId)
func (_Relay *RelayFilterer) ParseRetire(log types.Log) (*RelayRetire, error) {
	event := new(RelayRetire)
	if err := _Relay.contract.UnpackLog(event, "Retire", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayRotateIterator is returned from FilterRotate and is used to iterate over the raw logs and unpacked data for Rotate events raised by the Relay contract.
type RelayRotateIterator struct {
	Event *RelayRotate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayRotateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayRotate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayRotate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayRotateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayRotateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayRotate represents a Rotate event raised by the Relay contract.
type RelayRotate struct {
	CurrentId *big.Int
	NextId    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRotate is a free log retrieval operation binding the contract event 0xf5feff4c6b1d68649ea5159119be3b6e80b05e976f8bb16a8fa1202b2cf96db9.
//
// Solidity: event Rotate(uint256 currentId, uint256 nextId)
func (_Relay *RelayFilterer) FilterRotate(opts *bind.FilterOpts) (*RelayRotateIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Rotate")
	if err != nil {
		return nil, err
	}
	return &RelayRotateIterator{contract: _Relay.contract, event: "Rotate", logs: logs, sub: sub}, nil
}

// WatchRotate is a free log subscription operation binding the contract event 0xf5feff4c6b1d68649ea5159119be3b6e80b05e976f8bb16a8fa1202b2cf96db9.
//
// Solidity: event Rotate(uint256 currentId, uint256 nextId)
func (_Relay *RelayFilterer) WatchRotate(opts *bind.WatchOpts, sink chan<- *RelayRotate) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Rotate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayRotate)
				if err := _Relay.contract.UnpackLog(event, "Rotate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRotate is a log parse operation binding the contract event 0xf5feff4c6b1d68649ea5159119be3b6e80b05e976f8bb16a8fa1202b2cf96db9.
//
// Solidity: event Rotate(uint256 currentId, uint256 nextId)
func (_Relay *RelayFilterer) ParseRotate(log types.Log) (*RelayRotate, error) {
	event := new(RelayRotate)
	if err := _Relay.contract.UnpackLog(event, "Rotate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelaySetIterator is returned from FilterSet and is used to iterate over the raw logs and unpacked data for Set events raised by the Relay contract.
type RelaySetIterator struct {
	Event *RelaySet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelaySet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySet represents a Set event raised by the Relay contract.
type RelaySet struct {
	Maintainer common.Address
	Relay      common.Address
	Parameter  common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSet is a free log retrieval operation binding the contract event 0x2af5e77e7529216278d6036f07d924d3c5639a1b6b1f2164905f65c57165f726.
//
// Solidity: event Set(address _maintainer, address _relay, address _parameter)
func (_Relay *RelayFilterer) FilterSet(opts *bind.FilterOpts) (*RelaySetIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Set")
	if err != nil {
		return nil, err
	}
	return &RelaySetIterator{contract: _Relay.contract, event: "Set", logs: logs, sub: sub}, nil
}

// WatchSet is a free log subscription operation binding the contract event 0x2af5e77e7529216278d6036f07d924d3c5639a1b6b1f2164905f65c57165f726.
//
// Solidity: event Set(address _maintainer, address _relay, address _parameter)
func (_Relay *RelayFilterer) WatchSet(opts *bind.WatchOpts, sink chan<- *RelaySet) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Set")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySet)
				if err := _Relay.contract.UnpackLog(event, "Set", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSet is a log parse operation binding the contract event 0x2af5e77e7529216278d6036f07d924d3c5639a1b6b1f2164905f65c57165f726.
//
// Solidity: event Set(address _maintainer, address _relay, address _parameter)
func (_Relay *RelayFilterer) ParseSet(log types.Log) (*RelaySet, error) {
	event := new(RelaySet)
	if err := _Relay.contract.UnpackLog(event, "Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Relay contract.
type RelayUnpausedIterator struct {
	Event *RelayUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayUnpaused represents a Unpaused event raised by the Relay contract.
type RelayUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Relay *RelayFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RelayUnpausedIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RelayUnpausedIterator{contract: _Relay.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Relay *RelayFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RelayUnpaused) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayUnpaused)
				if err := _Relay.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Relay *RelayFilterer) ParseUnpaused(log types.Log) (*RelayUnpaused, error) {
	event := new(RelayUnpaused)
	if err := _Relay.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayUpdateKeyShareIterator is returned from FilterUpdateKeyShare and is used to iterate over the raw logs and unpacked data for UpdateKeyShare events raised by the Relay contract.
type RelayUpdateKeyShareIterator struct {
	Event *RelayUpdateKeyShare // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayUpdateKeyShareIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayUpdateKeyShare)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayUpdateKeyShare)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayUpdateKeyShareIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayUpdateKeyShareIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayUpdateKeyShare represents a UpdateKeyShare event raised by the Relay contract.
type RelayUpdateKeyShare struct {
	Maitainer common.Address
	Pubkey    []byte
	KeyShare  []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUpdateKeyShare is a free log retrieval operation binding the contract event 0xfedc1af7de1320c4cadf23429bfcbbb0abce84799c83e7971004f1468e376c4e.
//
// Solidity: event UpdateKeyShare(address maitainer, bytes pubkey, bytes keyShare)
func (_Relay *RelayFilterer) FilterUpdateKeyShare(opts *bind.FilterOpts) (*RelayUpdateKeyShareIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "UpdateKeyShare")
	if err != nil {
		return nil, err
	}
	return &RelayUpdateKeyShareIterator{contract: _Relay.contract, event: "UpdateKeyShare", logs: logs, sub: sub}, nil
}

// WatchUpdateKeyShare is a free log subscription operation binding the contract event 0xfedc1af7de1320c4cadf23429bfcbbb0abce84799c83e7971004f1468e376c4e.
//
// Solidity: event UpdateKeyShare(address maitainer, bytes pubkey, bytes keyShare)
func (_Relay *RelayFilterer) WatchUpdateKeyShare(opts *bind.WatchOpts, sink chan<- *RelayUpdateKeyShare) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "UpdateKeyShare")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayUpdateKeyShare)
				if err := _Relay.contract.UnpackLog(event, "UpdateKeyShare", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateKeyShare is a log parse operation binding the contract event 0xfedc1af7de1320c4cadf23429bfcbbb0abce84799c83e7971004f1468e376c4e.
//
// Solidity: event UpdateKeyShare(address maitainer, bytes pubkey, bytes keyShare)
func (_Relay *RelayFilterer) ParseUpdateKeyShare(log types.Log) (*RelayUpdateKeyShare, error) {
	event := new(RelayUpdateKeyShare)
	if err := _Relay.contract.UnpackLog(event, "UpdateKeyShare", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Relay contract.
type RelayUpgradedIterator struct {
	Event *RelayUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayUpgraded represents a Upgraded event raised by the Relay contract.
type RelayUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Relay *RelayFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RelayUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RelayUpgradedIterator{contract: _Relay.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Relay *RelayFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RelayUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayUpgraded)
				if err := _Relay.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Relay *RelayFilterer) ParseUpgraded(log types.Log) (*RelayUpgraded, error) {
	event := new(RelayUpgraded)
	if err := _Relay.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayVoteNetworkFeeIterator is returned from FilterVoteNetworkFee and is used to iterate over the raw logs and unpacked data for VoteNetworkFee events raised by the Relay contract.
type RelayVoteNetworkFeeIterator struct {
	Event *RelayVoteNetworkFee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayVoteNetworkFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayVoteNetworkFee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayVoteNetworkFee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayVoteNetworkFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayVoteNetworkFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayVoteNetworkFee represents a VoteNetworkFee event raised by the Relay contract.
type RelayVoteNetworkFee struct {
	Epoch  *big.Int
	Chain  *big.Int
	Height *big.Int
	Limit  *big.Int
	Price  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterVoteNetworkFee is a free log retrieval operation binding the contract event 0x54dbf9812f5917755f6626e3a7683913e5e9e1e2306ae993321ccb0872cb7217.
//
// Solidity: event VoteNetworkFee(uint256 epoch, uint256 chain, uint256 height, uint256 limit, uint256 price)
func (_Relay *RelayFilterer) FilterVoteNetworkFee(opts *bind.FilterOpts) (*RelayVoteNetworkFeeIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "VoteNetworkFee")
	if err != nil {
		return nil, err
	}
	return &RelayVoteNetworkFeeIterator{contract: _Relay.contract, event: "VoteNetworkFee", logs: logs, sub: sub}, nil
}

// WatchVoteNetworkFee is a free log subscription operation binding the contract event 0x54dbf9812f5917755f6626e3a7683913e5e9e1e2306ae993321ccb0872cb7217.
//
// Solidity: event VoteNetworkFee(uint256 epoch, uint256 chain, uint256 height, uint256 limit, uint256 price)
func (_Relay *RelayFilterer) WatchVoteNetworkFee(opts *bind.WatchOpts, sink chan<- *RelayVoteNetworkFee) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "VoteNetworkFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayVoteNetworkFee)
				if err := _Relay.contract.UnpackLog(event, "VoteNetworkFee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteNetworkFee is a log parse operation binding the contract event 0x54dbf9812f5917755f6626e3a7683913e5e9e1e2306ae993321ccb0872cb7217.
//
// Solidity: event VoteNetworkFee(uint256 epoch, uint256 chain, uint256 height, uint256 limit, uint256 price)
func (_Relay *RelayFilterer) ParseVoteNetworkFee(log types.Log) (*RelayVoteNetworkFee, error) {
	event := new(RelayVoteNetworkFee)
	if err := _Relay.contract.UnpackLog(event, "VoteNetworkFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayVoteTxInIterator is returned from FilterVoteTxIn and is used to iterate over the raw logs and unpacked data for VoteTxIn events raised by the Relay contract.
type RelayVoteTxInIterator struct {
	Event *RelayVoteTxIn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayVoteTxInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayVoteTxIn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayVoteTxIn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayVoteTxInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayVoteTxInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayVoteTxIn represents a VoteTxIn event raised by the Relay contract.
type RelayVoteTxIn struct {
	TxInItem TxInItem
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterVoteTxIn is a free log retrieval operation binding the contract event 0xdb251b03d781777e504f05b21ae8224e627da2a61a6a81eb17b6de3096e4b017.
//
// Solidity: event VoteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes) txInItem)
func (_Relay *RelayFilterer) FilterVoteTxIn(opts *bind.FilterOpts) (*RelayVoteTxInIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "VoteTxIn")
	if err != nil {
		return nil, err
	}
	return &RelayVoteTxInIterator{contract: _Relay.contract, event: "VoteTxIn", logs: logs, sub: sub}, nil
}

// WatchVoteTxIn is a free log subscription operation binding the contract event 0xdb251b03d781777e504f05b21ae8224e627da2a61a6a81eb17b6de3096e4b017.
//
// Solidity: event VoteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes) txInItem)
func (_Relay *RelayFilterer) WatchVoteTxIn(opts *bind.WatchOpts, sink chan<- *RelayVoteTxIn) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "VoteTxIn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayVoteTxIn)
				if err := _Relay.contract.UnpackLog(event, "VoteTxIn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteTxIn is a log parse operation binding the contract event 0xdb251b03d781777e504f05b21ae8224e627da2a61a6a81eb17b6de3096e4b017.
//
// Solidity: event VoteTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,bytes) txInItem)
func (_Relay *RelayFilterer) ParseVoteTxIn(log types.Log) (*RelayVoteTxIn, error) {
	event := new(RelayVoteTxIn)
	if err := _Relay.contract.UnpackLog(event, "VoteTxIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayVoteTxOutIterator is returned from FilterVoteTxOut and is used to iterate over the raw logs and unpacked data for VoteTxOut events raised by the Relay contract.
type RelayVoteTxOutIterator struct {
	Event *RelayVoteTxOut // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayVoteTxOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayVoteTxOut)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayVoteTxOut)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayVoteTxOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayVoteTxOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayVoteTxOut represents a VoteTxOut event raised by the Relay contract.
type RelayVoteTxOut struct {
	TxOutItem TxOutItem
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVoteTxOut is a free log retrieval operation binding the contract event 0xe4e9e5fd4f80156531cb7b404a661401b6cb8c1e3f91b7e3656207bbee8db567.
//
// Solidity: event VoteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address) txOutItem)
func (_Relay *RelayFilterer) FilterVoteTxOut(opts *bind.FilterOpts) (*RelayVoteTxOutIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "VoteTxOut")
	if err != nil {
		return nil, err
	}
	return &RelayVoteTxOutIterator{contract: _Relay.contract, event: "VoteTxOut", logs: logs, sub: sub}, nil
}

// WatchVoteTxOut is a free log subscription operation binding the contract event 0xe4e9e5fd4f80156531cb7b404a661401b6cb8c1e3f91b7e3656207bbee8db567.
//
// Solidity: event VoteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address) txOutItem)
func (_Relay *RelayFilterer) WatchVoteTxOut(opts *bind.WatchOpts, sink chan<- *RelayVoteTxOut) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "VoteTxOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayVoteTxOut)
				if err := _Relay.contract.UnpackLog(event, "VoteTxOut", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteTxOut is a log parse operation binding the contract event 0xe4e9e5fd4f80156531cb7b404a661401b6cb8c1e3f91b7e3656207bbee8db567.
//
// Solidity: event VoteTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint128,uint128,address) txOutItem)
func (_Relay *RelayFilterer) ParseVoteTxOut(log types.Log) (*RelayVoteTxOut, error) {
	event := new(RelayVoteTxOut)
	if err := _Relay.contract.UnpackLog(event, "VoteTxOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayVoteUpdateTssPoolIterator is returned from FilterVoteUpdateTssPool and is used to iterate over the raw logs and unpacked data for VoteUpdateTssPool events raised by the Relay contract.
type RelayVoteUpdateTssPoolIterator struct {
	Event *RelayVoteUpdateTssPool // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RelayVoteUpdateTssPoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayVoteUpdateTssPool)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RelayVoteUpdateTssPool)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RelayVoteUpdateTssPoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayVoteUpdateTssPoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayVoteUpdateTssPool represents a VoteUpdateTssPool event raised by the Relay contract.
type RelayVoteUpdateTssPool struct {
	Param TSSManagerTssPoolParam
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVoteUpdateTssPool is a free log retrieval operation binding the contract event 0x6a015f8d2cfd5d763d34c11d62f84e381b28ed461cd0fd63aaf749947e1a9f10.
//
// Solidity: event VoteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param)
func (_Relay *RelayFilterer) FilterVoteUpdateTssPool(opts *bind.FilterOpts) (*RelayVoteUpdateTssPoolIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "VoteUpdateTssPool")
	if err != nil {
		return nil, err
	}
	return &RelayVoteUpdateTssPoolIterator{contract: _Relay.contract, event: "VoteUpdateTssPool", logs: logs, sub: sub}, nil
}

// WatchVoteUpdateTssPool is a free log subscription operation binding the contract event 0x6a015f8d2cfd5d763d34c11d62f84e381b28ed461cd0fd63aaf749947e1a9f10.
//
// Solidity: event VoteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param)
func (_Relay *RelayFilterer) WatchVoteUpdateTssPool(opts *bind.WatchOpts, sink chan<- *RelayVoteUpdateTssPool) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "VoteUpdateTssPool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayVoteUpdateTssPool)
				if err := _Relay.contract.UnpackLog(event, "VoteUpdateTssPool", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteUpdateTssPool is a log parse operation binding the contract event 0x6a015f8d2cfd5d763d34c11d62f84e381b28ed461cd0fd63aaf749947e1a9f10.
//
// Solidity: event VoteUpdateTssPool((uint256,bytes,bytes,address[],address[],bytes) param)
func (_Relay *RelayFilterer) ParseVoteUpdateTssPool(log types.Log) (*RelayVoteUpdateTssPool, error) {
	event := new(RelayVoteUpdateTssPool)
	if err := _Relay.contract.UnpackLog(event, "VoteUpdateTssPool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
