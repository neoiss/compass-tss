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

// TxInItem is an auto generated low-level Go binding around an user-defined struct.
type TxInItem struct {
	OrderId    [32]byte
	BridgeItem BridgeItem
	Height     uint64
	RefundAddr []byte
}

// TxItem is an auto generated low-level Go binding around an user-defined struct.
type TxItem struct {
	OrderId   [32]byte
	Chain     *big.Int
	ChainType uint8
	Token     common.Address
	Amount    *big.Int
	To        common.Address
}

// TxOutItem is an auto generated low-level Go binding around an user-defined struct.
type TxOutItem struct {
	OrderId    [32]byte
	BridgeItem BridgeItem
	Height     uint64
	GasUsed    *big.Int
	Sender     common.Address
}

// RelayMetaData contains all meta data concerning the Relay contract.
var RelayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AccessManagedInvalidAuthority\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"delay\",\"type\":\"uint32\"}],\"name\":\"AccessManagedRequiredDelay\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"AccessManagedUnauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"expired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_refund_address\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_signature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_vault\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"invalid_vault_token\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"migration_not_completed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"no_access\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"not_bridge_able\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"order_executed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"transfer_in_failed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"transfer_out_failed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"transfer_token_out_failed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"vault_token_not_registered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"zero_address\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"zero_amount_out\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"authority\",\"type\":\"address\"}],\"name\":\"AuthorityUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumTxType\",\"name\":\"txOutType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"BridgeCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"BridgeFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"securityFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vaultFee\",\"type\":\"uint256\"}],\"name\":\"BridgeFeeCollected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumTxType\",\"name\":\"txInType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"BridgeIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumTxType\",\"name\":\"txOutType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"BridgeOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"BridgeRelay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"relayData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"BridgeRelaySigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromChain\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_affiliateFeeManager\",\"type\":\"address\"}],\"name\":\"SetAffiliateFeeManager\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_periphery\",\"type\":\"address\"}],\"name\":\"SetPeriphery\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_swap\",\"type\":\"address\"}],\"name\":\"SetSwap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_vaultManager\",\"type\":\"address\"}],\"name\":\"SetVaultManager\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_wToken\",\"type\":\"address\"}],\"name\":\"SetWToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"result\",\"type\":\"bool\"}],\"name\":\"TransferIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feature\",\"type\":\"uint256\"}],\"name\":\"UpdateTokens\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"reicerver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vaultAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"}],\"name\":\"addChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"affiliateFeeManager\",\"outputs\":[{\"internalType\":\"contractIAffiliateFeeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceFeeInfos\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toChain\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"bridgeOut\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"},{\"internalType\":\"enumChainType\",\"name\":\"chainType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"internalType\":\"structTxItem\",\"name\":\"txItem\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"fromChain\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"relayPayload\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"targetPayload\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"refundAddr\",\"type\":\"bytes\"}],\"internalType\":\"structTxInItem\",\"name\":\"txInItem\",\"type\":\"tuple\"}],\"name\":\"executeTxIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainAndGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"vault\",\"type\":\"bytes\"},{\"internalType\":\"enumTxType\",\"name\":\"txType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"internalType\":\"structBridgeItem\",\"name\":\"bridgeItem\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"gasUsed\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTxOutItem\",\"name\":\"txOutItem\",\"type\":\"tuple\"}],\"name\":\"executeTxOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"failedHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"}],\"name\":\"getChainLastScanBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_defaultAdmin\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"isBridgeable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isConsumingScheduledOp\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"isMintable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isTxIn\",\"type\":\"bool\"}],\"name\":\"isOrderExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orderInfos\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"signed\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"gasToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"estimateGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"periphery\",\"outputs\":[{\"internalType\":\"contractIPeriphery\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionSizeWithCall\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transactionRate\",\"type\":\"uint256\"}],\"name\":\"postNetworkFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"relayData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"relaySigned\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chain\",\"type\":\"uint256\"}],\"name\":\"removeChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"retiringVault\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"activeVault\",\"type\":\"bytes\"}],\"name\":\"rotate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"securityFeeReceiver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selfChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_affiliateFeeManager\",\"type\":\"address\"}],\"name\":\"setAffiliateFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAuthority\",\"type\":\"address\"}],\"name\":\"setAuthority\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_periphery\",\"type\":\"address\"}],\"name\":\"setPeriphery\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_swap\",\"type\":\"address\"}],\"name\":\"setSwap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vaultManager\",\"type\":\"address\"}],\"name\":\"setVaultManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wToken\",\"type\":\"address\"}],\"name\":\"setWtoken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"contractISwap\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenFeatureList\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trigger\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_feature\",\"type\":\"uint256\"}],\"name\":\"updateTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"vaultFeeInfos\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vaultManager\",\"outputs\":[{\"internalType\":\"contractIVaultManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_vaultToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_vaultAmount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// AffiliateFeeManager is a free data retrieval call binding the contract method 0xb2177ba4.
//
// Solidity: function affiliateFeeManager() view returns(address)
func (_Relay *RelayCaller) AffiliateFeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "affiliateFeeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AffiliateFeeManager is a free data retrieval call binding the contract method 0xb2177ba4.
//
// Solidity: function affiliateFeeManager() view returns(address)
func (_Relay *RelaySession) AffiliateFeeManager() (common.Address, error) {
	return _Relay.Contract.AffiliateFeeManager(&_Relay.CallOpts)
}

// AffiliateFeeManager is a free data retrieval call binding the contract method 0xb2177ba4.
//
// Solidity: function affiliateFeeManager() view returns(address)
func (_Relay *RelayCallerSession) AffiliateFeeManager() (common.Address, error) {
	return _Relay.Contract.AffiliateFeeManager(&_Relay.CallOpts)
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

// BalanceFeeInfos is a free data retrieval call binding the contract method 0x2035f1a2.
//
// Solidity: function balanceFeeInfos(address ) view returns(uint256)
func (_Relay *RelayCaller) BalanceFeeInfos(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "balanceFeeInfos", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceFeeInfos is a free data retrieval call binding the contract method 0x2035f1a2.
//
// Solidity: function balanceFeeInfos(address ) view returns(uint256)
func (_Relay *RelaySession) BalanceFeeInfos(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.BalanceFeeInfos(&_Relay.CallOpts, arg0)
}

// BalanceFeeInfos is a free data retrieval call binding the contract method 0x2035f1a2.
//
// Solidity: function balanceFeeInfos(address ) view returns(uint256)
func (_Relay *RelayCallerSession) BalanceFeeInfos(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.BalanceFeeInfos(&_Relay.CallOpts, arg0)
}

// FailedHash is a free data retrieval call binding the contract method 0x488e4c45.
//
// Solidity: function failedHash(bytes32 ) view returns(bool)
func (_Relay *RelayCaller) FailedHash(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "failedHash", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FailedHash is a free data retrieval call binding the contract method 0x488e4c45.
//
// Solidity: function failedHash(bytes32 ) view returns(bool)
func (_Relay *RelaySession) FailedHash(arg0 [32]byte) (bool, error) {
	return _Relay.Contract.FailedHash(&_Relay.CallOpts, arg0)
}

// FailedHash is a free data retrieval call binding the contract method 0x488e4c45.
//
// Solidity: function failedHash(bytes32 ) view returns(bool)
func (_Relay *RelayCallerSession) FailedHash(arg0 [32]byte) (bool, error) {
	return _Relay.Contract.FailedHash(&_Relay.CallOpts, arg0)
}

// GetChainLastScanBlock is a free data retrieval call binding the contract method 0xe5ff7a20.
//
// Solidity: function getChainLastScanBlock(uint256 chain) view returns(uint256)
func (_Relay *RelayCaller) GetChainLastScanBlock(opts *bind.CallOpts, chain *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "getChainLastScanBlock", chain)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChainLastScanBlock is a free data retrieval call binding the contract method 0xe5ff7a20.
//
// Solidity: function getChainLastScanBlock(uint256 chain) view returns(uint256)
func (_Relay *RelaySession) GetChainLastScanBlock(chain *big.Int) (*big.Int, error) {
	return _Relay.Contract.GetChainLastScanBlock(&_Relay.CallOpts, chain)
}

// GetChainLastScanBlock is a free data retrieval call binding the contract method 0xe5ff7a20.
//
// Solidity: function getChainLastScanBlock(uint256 chain) view returns(uint256)
func (_Relay *RelayCallerSession) GetChainLastScanBlock(chain *big.Int) (*big.Int, error) {
	return _Relay.Contract.GetChainLastScanBlock(&_Relay.CallOpts, chain)
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

// IsBridgeable is a free data retrieval call binding the contract method 0xbef920b2.
//
// Solidity: function isBridgeable(address _token) view returns(bool)
func (_Relay *RelayCaller) IsBridgeable(opts *bind.CallOpts, _token common.Address) (bool, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "isBridgeable", _token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBridgeable is a free data retrieval call binding the contract method 0xbef920b2.
//
// Solidity: function isBridgeable(address _token) view returns(bool)
func (_Relay *RelaySession) IsBridgeable(_token common.Address) (bool, error) {
	return _Relay.Contract.IsBridgeable(&_Relay.CallOpts, _token)
}

// IsBridgeable is a free data retrieval call binding the contract method 0xbef920b2.
//
// Solidity: function isBridgeable(address _token) view returns(bool)
func (_Relay *RelayCallerSession) IsBridgeable(_token common.Address) (bool, error) {
	return _Relay.Contract.IsBridgeable(&_Relay.CallOpts, _token)
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

// IsMintable is a free data retrieval call binding the contract method 0x222b15fb.
//
// Solidity: function isMintable(address _token) view returns(bool)
func (_Relay *RelayCaller) IsMintable(opts *bind.CallOpts, _token common.Address) (bool, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "isMintable", _token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMintable is a free data retrieval call binding the contract method 0x222b15fb.
//
// Solidity: function isMintable(address _token) view returns(bool)
func (_Relay *RelaySession) IsMintable(_token common.Address) (bool, error) {
	return _Relay.Contract.IsMintable(&_Relay.CallOpts, _token)
}

// IsMintable is a free data retrieval call binding the contract method 0x222b15fb.
//
// Solidity: function isMintable(address _token) view returns(bool)
func (_Relay *RelayCallerSession) IsMintable(_token common.Address) (bool, error) {
	return _Relay.Contract.IsMintable(&_Relay.CallOpts, _token)
}

// IsOrderExecuted is a free data retrieval call binding the contract method 0xb078dafb.
//
// Solidity: function isOrderExecuted(bytes32 orderId, bool isTxIn) view returns(bool executed)
func (_Relay *RelayCaller) IsOrderExecuted(opts *bind.CallOpts, orderId [32]byte, isTxIn bool) (bool, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "isOrderExecuted", orderId, isTxIn)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOrderExecuted is a free data retrieval call binding the contract method 0xb078dafb.
//
// Solidity: function isOrderExecuted(bytes32 orderId, bool isTxIn) view returns(bool executed)
func (_Relay *RelaySession) IsOrderExecuted(orderId [32]byte, isTxIn bool) (bool, error) {
	return _Relay.Contract.IsOrderExecuted(&_Relay.CallOpts, orderId, isTxIn)
}

// IsOrderExecuted is a free data retrieval call binding the contract method 0xb078dafb.
//
// Solidity: function isOrderExecuted(bytes32 orderId, bool isTxIn) view returns(bool executed)
func (_Relay *RelayCallerSession) IsOrderExecuted(orderId [32]byte, isTxIn bool) (bool, error) {
	return _Relay.Contract.IsOrderExecuted(&_Relay.CallOpts, orderId, isTxIn)
}

// OrderInfos is a free data retrieval call binding the contract method 0x266de649.
//
// Solidity: function orderInfos(bytes32 ) view returns(bool signed, uint64 height, address gasToken, uint256 estimateGas, bytes32 hash)
func (_Relay *RelayCaller) OrderInfos(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Signed      bool
	Height      uint64
	GasToken    common.Address
	EstimateGas *big.Int
	Hash        [32]byte
}, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "orderInfos", arg0)

	outstruct := new(struct {
		Signed      bool
		Height      uint64
		GasToken    common.Address
		EstimateGas *big.Int
		Hash        [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Signed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Height = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.GasToken = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.EstimateGas = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Hash = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// OrderInfos is a free data retrieval call binding the contract method 0x266de649.
//
// Solidity: function orderInfos(bytes32 ) view returns(bool signed, uint64 height, address gasToken, uint256 estimateGas, bytes32 hash)
func (_Relay *RelaySession) OrderInfos(arg0 [32]byte) (struct {
	Signed      bool
	Height      uint64
	GasToken    common.Address
	EstimateGas *big.Int
	Hash        [32]byte
}, error) {
	return _Relay.Contract.OrderInfos(&_Relay.CallOpts, arg0)
}

// OrderInfos is a free data retrieval call binding the contract method 0x266de649.
//
// Solidity: function orderInfos(bytes32 ) view returns(bool signed, uint64 height, address gasToken, uint256 estimateGas, bytes32 hash)
func (_Relay *RelayCallerSession) OrderInfos(arg0 [32]byte) (struct {
	Signed      bool
	Height      uint64
	GasToken    common.Address
	EstimateGas *big.Int
	Hash        [32]byte
}, error) {
	return _Relay.Contract.OrderInfos(&_Relay.CallOpts, arg0)
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

// Periphery is a free data retrieval call binding the contract method 0x77aace1a.
//
// Solidity: function periphery() view returns(address)
func (_Relay *RelayCaller) Periphery(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "periphery")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Periphery is a free data retrieval call binding the contract method 0x77aace1a.
//
// Solidity: function periphery() view returns(address)
func (_Relay *RelaySession) Periphery() (common.Address, error) {
	return _Relay.Contract.Periphery(&_Relay.CallOpts)
}

// Periphery is a free data retrieval call binding the contract method 0x77aace1a.
//
// Solidity: function periphery() view returns(address)
func (_Relay *RelayCallerSession) Periphery() (common.Address, error) {
	return _Relay.Contract.Periphery(&_Relay.CallOpts)
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

// SecurityFeeReceiver is a free data retrieval call binding the contract method 0xc2dae9a9.
//
// Solidity: function securityFeeReceiver() view returns(address)
func (_Relay *RelayCaller) SecurityFeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "securityFeeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SecurityFeeReceiver is a free data retrieval call binding the contract method 0xc2dae9a9.
//
// Solidity: function securityFeeReceiver() view returns(address)
func (_Relay *RelaySession) SecurityFeeReceiver() (common.Address, error) {
	return _Relay.Contract.SecurityFeeReceiver(&_Relay.CallOpts)
}

// SecurityFeeReceiver is a free data retrieval call binding the contract method 0xc2dae9a9.
//
// Solidity: function securityFeeReceiver() view returns(address)
func (_Relay *RelayCallerSession) SecurityFeeReceiver() (common.Address, error) {
	return _Relay.Contract.SecurityFeeReceiver(&_Relay.CallOpts)
}

// SelfChainId is a free data retrieval call binding the contract method 0xcc9e3e89.
//
// Solidity: function selfChainId() view returns(uint256)
func (_Relay *RelayCaller) SelfChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "selfChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SelfChainId is a free data retrieval call binding the contract method 0xcc9e3e89.
//
// Solidity: function selfChainId() view returns(uint256)
func (_Relay *RelaySession) SelfChainId() (*big.Int, error) {
	return _Relay.Contract.SelfChainId(&_Relay.CallOpts)
}

// SelfChainId is a free data retrieval call binding the contract method 0xcc9e3e89.
//
// Solidity: function selfChainId() view returns(uint256)
func (_Relay *RelayCallerSession) SelfChainId() (*big.Int, error) {
	return _Relay.Contract.SelfChainId(&_Relay.CallOpts)
}

// Swap is a free data retrieval call binding the contract method 0x8119c065.
//
// Solidity: function swap() view returns(address)
func (_Relay *RelayCaller) Swap(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "swap")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Swap is a free data retrieval call binding the contract method 0x8119c065.
//
// Solidity: function swap() view returns(address)
func (_Relay *RelaySession) Swap() (common.Address, error) {
	return _Relay.Contract.Swap(&_Relay.CallOpts)
}

// Swap is a free data retrieval call binding the contract method 0x8119c065.
//
// Solidity: function swap() view returns(address)
func (_Relay *RelayCallerSession) Swap() (common.Address, error) {
	return _Relay.Contract.Swap(&_Relay.CallOpts)
}

// TokenFeatureList is a free data retrieval call binding the contract method 0x0a83a4e9.
//
// Solidity: function tokenFeatureList(address ) view returns(uint256)
func (_Relay *RelayCaller) TokenFeatureList(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "tokenFeatureList", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenFeatureList is a free data retrieval call binding the contract method 0x0a83a4e9.
//
// Solidity: function tokenFeatureList(address ) view returns(uint256)
func (_Relay *RelaySession) TokenFeatureList(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.TokenFeatureList(&_Relay.CallOpts, arg0)
}

// TokenFeatureList is a free data retrieval call binding the contract method 0x0a83a4e9.
//
// Solidity: function tokenFeatureList(address ) view returns(uint256)
func (_Relay *RelayCallerSession) TokenFeatureList(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.TokenFeatureList(&_Relay.CallOpts, arg0)
}

// VaultFeeInfos is a free data retrieval call binding the contract method 0x7b15f4e8.
//
// Solidity: function vaultFeeInfos(address ) view returns(uint256)
func (_Relay *RelayCaller) VaultFeeInfos(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "vaultFeeInfos", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VaultFeeInfos is a free data retrieval call binding the contract method 0x7b15f4e8.
//
// Solidity: function vaultFeeInfos(address ) view returns(uint256)
func (_Relay *RelaySession) VaultFeeInfos(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.VaultFeeInfos(&_Relay.CallOpts, arg0)
}

// VaultFeeInfos is a free data retrieval call binding the contract method 0x7b15f4e8.
//
// Solidity: function vaultFeeInfos(address ) view returns(uint256)
func (_Relay *RelayCallerSession) VaultFeeInfos(arg0 common.Address) (*big.Int, error) {
	return _Relay.Contract.VaultFeeInfos(&_Relay.CallOpts, arg0)
}

// VaultManager is a free data retrieval call binding the contract method 0x8a4adf24.
//
// Solidity: function vaultManager() view returns(address)
func (_Relay *RelayCaller) VaultManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "vaultManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VaultManager is a free data retrieval call binding the contract method 0x8a4adf24.
//
// Solidity: function vaultManager() view returns(address)
func (_Relay *RelaySession) VaultManager() (common.Address, error) {
	return _Relay.Contract.VaultManager(&_Relay.CallOpts)
}

// VaultManager is a free data retrieval call binding the contract method 0x8a4adf24.
//
// Solidity: function vaultManager() view returns(address)
func (_Relay *RelayCallerSession) VaultManager() (common.Address, error) {
	return _Relay.Contract.VaultManager(&_Relay.CallOpts)
}

// WToken is a free data retrieval call binding the contract method 0x0babd864.
//
// Solidity: function wToken() view returns(address)
func (_Relay *RelayCaller) WToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Relay.contract.Call(opts, &out, "wToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WToken is a free data retrieval call binding the contract method 0x0babd864.
//
// Solidity: function wToken() view returns(address)
func (_Relay *RelaySession) WToken() (common.Address, error) {
	return _Relay.Contract.WToken(&_Relay.CallOpts)
}

// WToken is a free data retrieval call binding the contract method 0x0babd864.
//
// Solidity: function wToken() view returns(address)
func (_Relay *RelayCallerSession) WToken() (common.Address, error) {
	return _Relay.Contract.WToken(&_Relay.CallOpts)
}

// AddChain is a paid mutator transaction binding the contract method 0x5d21ef26.
//
// Solidity: function addChain(uint256 chain, uint256 startBlock) returns()
func (_Relay *RelayTransactor) AddChain(opts *bind.TransactOpts, chain *big.Int, startBlock *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "addChain", chain, startBlock)
}

// AddChain is a paid mutator transaction binding the contract method 0x5d21ef26.
//
// Solidity: function addChain(uint256 chain, uint256 startBlock) returns()
func (_Relay *RelaySession) AddChain(chain *big.Int, startBlock *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.AddChain(&_Relay.TransactOpts, chain, startBlock)
}

// AddChain is a paid mutator transaction binding the contract method 0x5d21ef26.
//
// Solidity: function addChain(uint256 chain, uint256 startBlock) returns()
func (_Relay *RelayTransactorSession) AddChain(chain *big.Int, startBlock *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.AddChain(&_Relay.TransactOpts, chain, startBlock)
}

// BridgeOut is a paid mutator transaction binding the contract method 0x34631986.
//
// Solidity: function bridgeOut(address token, uint256 amount, uint256 toChain, bytes to, address refundAddr, bytes payload, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelayTransactor) BridgeOut(opts *bind.TransactOpts, token common.Address, amount *big.Int, toChain *big.Int, to []byte, refundAddr common.Address, payload []byte, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "bridgeOut", token, amount, toChain, to, refundAddr, payload, deadline)
}

// BridgeOut is a paid mutator transaction binding the contract method 0x34631986.
//
// Solidity: function bridgeOut(address token, uint256 amount, uint256 toChain, bytes to, address refundAddr, bytes payload, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelaySession) BridgeOut(token common.Address, amount *big.Int, toChain *big.Int, to []byte, refundAddr common.Address, payload []byte, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.BridgeOut(&_Relay.TransactOpts, token, amount, toChain, to, refundAddr, payload, deadline)
}

// BridgeOut is a paid mutator transaction binding the contract method 0x34631986.
//
// Solidity: function bridgeOut(address token, uint256 amount, uint256 toChain, bytes to, address refundAddr, bytes payload, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelayTransactorSession) BridgeOut(token common.Address, amount *big.Int, toChain *big.Int, to []byte, refundAddr common.Address, payload []byte, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.BridgeOut(&_Relay.TransactOpts, token, amount, toChain, to, refundAddr, payload, deadline)
}

// Deposit is a paid mutator transaction binding the contract method 0x3fb7de52.
//
// Solidity: function deposit(address token, uint256 amount, address to, address refundAddr, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelayTransactor) Deposit(opts *bind.TransactOpts, token common.Address, amount *big.Int, to common.Address, refundAddr common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "deposit", token, amount, to, refundAddr, deadline)
}

// Deposit is a paid mutator transaction binding the contract method 0x3fb7de52.
//
// Solidity: function deposit(address token, uint256 amount, address to, address refundAddr, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelaySession) Deposit(token common.Address, amount *big.Int, to common.Address, refundAddr common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Deposit(&_Relay.TransactOpts, token, amount, to, refundAddr, deadline)
}

// Deposit is a paid mutator transaction binding the contract method 0x3fb7de52.
//
// Solidity: function deposit(address token, uint256 amount, address to, address refundAddr, uint256 deadline) payable returns(bytes32 orderId)
func (_Relay *RelayTransactorSession) Deposit(token common.Address, amount *big.Int, to common.Address, refundAddr common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Deposit(&_Relay.TransactOpts, token, amount, to, refundAddr, deadline)
}

// Execute is a paid mutator transaction binding the contract method 0x234684ba.
//
// Solidity: function execute((uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes) bridgeItem, (bytes32,uint256,uint8,address,uint256,address) txItem, uint256 fromChain, bytes relayPayload, bytes targetPayload) returns(uint256)
func (_Relay *RelayTransactor) Execute(opts *bind.TransactOpts, bridgeItem BridgeItem, txItem TxItem, fromChain *big.Int, relayPayload []byte, targetPayload []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "execute", bridgeItem, txItem, fromChain, relayPayload, targetPayload)
}

// Execute is a paid mutator transaction binding the contract method 0x234684ba.
//
// Solidity: function execute((uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes) bridgeItem, (bytes32,uint256,uint8,address,uint256,address) txItem, uint256 fromChain, bytes relayPayload, bytes targetPayload) returns(uint256)
func (_Relay *RelaySession) Execute(bridgeItem BridgeItem, txItem TxItem, fromChain *big.Int, relayPayload []byte, targetPayload []byte) (*types.Transaction, error) {
	return _Relay.Contract.Execute(&_Relay.TransactOpts, bridgeItem, txItem, fromChain, relayPayload, targetPayload)
}

// Execute is a paid mutator transaction binding the contract method 0x234684ba.
//
// Solidity: function execute((uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes) bridgeItem, (bytes32,uint256,uint8,address,uint256,address) txItem, uint256 fromChain, bytes relayPayload, bytes targetPayload) returns(uint256)
func (_Relay *RelayTransactorSession) Execute(bridgeItem BridgeItem, txItem TxItem, fromChain *big.Int, relayPayload []byte, targetPayload []byte) (*types.Transaction, error) {
	return _Relay.Contract.Execute(&_Relay.TransactOpts, bridgeItem, txItem, fromChain, relayPayload, targetPayload)
}

// ExecuteTxIn is a paid mutator transaction binding the contract method 0xeb482063.
//
// Solidity: function executeTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,bytes) txInItem) returns()
func (_Relay *RelayTransactor) ExecuteTxIn(opts *bind.TransactOpts, txInItem TxInItem) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "executeTxIn", txInItem)
}

// ExecuteTxIn is a paid mutator transaction binding the contract method 0xeb482063.
//
// Solidity: function executeTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,bytes) txInItem) returns()
func (_Relay *RelaySession) ExecuteTxIn(txInItem TxInItem) (*types.Transaction, error) {
	return _Relay.Contract.ExecuteTxIn(&_Relay.TransactOpts, txInItem)
}

// ExecuteTxIn is a paid mutator transaction binding the contract method 0xeb482063.
//
// Solidity: function executeTxIn((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,bytes) txInItem) returns()
func (_Relay *RelayTransactorSession) ExecuteTxIn(txInItem TxInItem) (*types.Transaction, error) {
	return _Relay.Contract.ExecuteTxIn(&_Relay.TransactOpts, txInItem)
}

// ExecuteTxOut is a paid mutator transaction binding the contract method 0x8dcc73ac.
//
// Solidity: function executeTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,uint128,address) txOutItem) returns()
func (_Relay *RelayTransactor) ExecuteTxOut(opts *bind.TransactOpts, txOutItem TxOutItem) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "executeTxOut", txOutItem)
}

// ExecuteTxOut is a paid mutator transaction binding the contract method 0x8dcc73ac.
//
// Solidity: function executeTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,uint128,address) txOutItem) returns()
func (_Relay *RelaySession) ExecuteTxOut(txOutItem TxOutItem) (*types.Transaction, error) {
	return _Relay.Contract.ExecuteTxOut(&_Relay.TransactOpts, txOutItem)
}

// ExecuteTxOut is a paid mutator transaction binding the contract method 0x8dcc73ac.
//
// Solidity: function executeTxOut((bytes32,(uint256,bytes,uint8,uint256,bytes,uint256,bytes,bytes,bytes),uint64,uint128,address) txOutItem) returns()
func (_Relay *RelayTransactorSession) ExecuteTxOut(txOutItem TxOutItem) (*types.Transaction, error) {
	return _Relay.Contract.ExecuteTxOut(&_Relay.TransactOpts, txOutItem)
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
// Solidity: function migrate() returns(bool)
func (_Relay *RelayTransactor) Migrate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "migrate")
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns(bool)
func (_Relay *RelaySession) Migrate() (*types.Transaction, error) {
	return _Relay.Contract.Migrate(&_Relay.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0x8fd3ab80.
//
// Solidity: function migrate() returns(bool)
func (_Relay *RelayTransactorSession) Migrate() (*types.Transaction, error) {
	return _Relay.Contract.Migrate(&_Relay.TransactOpts)
}

// PostNetworkFee is a paid mutator transaction binding the contract method 0x3ccf9acf.
//
// Solidity: function postNetworkFee(uint256 chain, uint256 height, uint256 transactionSize, uint256 transactionSizeWithCall, uint256 transactionRate) returns()
func (_Relay *RelayTransactor) PostNetworkFee(opts *bind.TransactOpts, chain *big.Int, height *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int, transactionRate *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "postNetworkFee", chain, height, transactionSize, transactionSizeWithCall, transactionRate)
}

// PostNetworkFee is a paid mutator transaction binding the contract method 0x3ccf9acf.
//
// Solidity: function postNetworkFee(uint256 chain, uint256 height, uint256 transactionSize, uint256 transactionSizeWithCall, uint256 transactionRate) returns()
func (_Relay *RelaySession) PostNetworkFee(chain *big.Int, height *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int, transactionRate *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.PostNetworkFee(&_Relay.TransactOpts, chain, height, transactionSize, transactionSizeWithCall, transactionRate)
}

// PostNetworkFee is a paid mutator transaction binding the contract method 0x3ccf9acf.
//
// Solidity: function postNetworkFee(uint256 chain, uint256 height, uint256 transactionSize, uint256 transactionSizeWithCall, uint256 transactionRate) returns()
func (_Relay *RelayTransactorSession) PostNetworkFee(chain *big.Int, height *big.Int, transactionSize *big.Int, transactionSizeWithCall *big.Int, transactionRate *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.PostNetworkFee(&_Relay.TransactOpts, chain, height, transactionSize, transactionSizeWithCall, transactionRate)
}

// RelaySigned is a paid mutator transaction binding the contract method 0xdba4364a.
//
// Solidity: function relaySigned(bytes32 orderId, bytes relayData, bytes signature) returns()
func (_Relay *RelayTransactor) RelaySigned(opts *bind.TransactOpts, orderId [32]byte, relayData []byte, signature []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "relaySigned", orderId, relayData, signature)
}

// RelaySigned is a paid mutator transaction binding the contract method 0xdba4364a.
//
// Solidity: function relaySigned(bytes32 orderId, bytes relayData, bytes signature) returns()
func (_Relay *RelaySession) RelaySigned(orderId [32]byte, relayData []byte, signature []byte) (*types.Transaction, error) {
	return _Relay.Contract.RelaySigned(&_Relay.TransactOpts, orderId, relayData, signature)
}

// RelaySigned is a paid mutator transaction binding the contract method 0xdba4364a.
//
// Solidity: function relaySigned(bytes32 orderId, bytes relayData, bytes signature) returns()
func (_Relay *RelayTransactorSession) RelaySigned(orderId [32]byte, relayData []byte, signature []byte) (*types.Transaction, error) {
	return _Relay.Contract.RelaySigned(&_Relay.TransactOpts, orderId, relayData, signature)
}

// RemoveChain is a paid mutator transaction binding the contract method 0xcbf50bb6.
//
// Solidity: function removeChain(uint256 chain) returns()
func (_Relay *RelayTransactor) RemoveChain(opts *bind.TransactOpts, chain *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "removeChain", chain)
}

// RemoveChain is a paid mutator transaction binding the contract method 0xcbf50bb6.
//
// Solidity: function removeChain(uint256 chain) returns()
func (_Relay *RelaySession) RemoveChain(chain *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.RemoveChain(&_Relay.TransactOpts, chain)
}

// RemoveChain is a paid mutator transaction binding the contract method 0xcbf50bb6.
//
// Solidity: function removeChain(uint256 chain) returns()
func (_Relay *RelayTransactorSession) RemoveChain(chain *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.RemoveChain(&_Relay.TransactOpts, chain)
}

// Rotate is a paid mutator transaction binding the contract method 0x75e98c4d.
//
// Solidity: function rotate(bytes retiringVault, bytes activeVault) returns()
func (_Relay *RelayTransactor) Rotate(opts *bind.TransactOpts, retiringVault []byte, activeVault []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "rotate", retiringVault, activeVault)
}

// Rotate is a paid mutator transaction binding the contract method 0x75e98c4d.
//
// Solidity: function rotate(bytes retiringVault, bytes activeVault) returns()
func (_Relay *RelaySession) Rotate(retiringVault []byte, activeVault []byte) (*types.Transaction, error) {
	return _Relay.Contract.Rotate(&_Relay.TransactOpts, retiringVault, activeVault)
}

// Rotate is a paid mutator transaction binding the contract method 0x75e98c4d.
//
// Solidity: function rotate(bytes retiringVault, bytes activeVault) returns()
func (_Relay *RelayTransactorSession) Rotate(retiringVault []byte, activeVault []byte) (*types.Transaction, error) {
	return _Relay.Contract.Rotate(&_Relay.TransactOpts, retiringVault, activeVault)
}

// SetAffiliateFeeManager is a paid mutator transaction binding the contract method 0x3fd77557.
//
// Solidity: function setAffiliateFeeManager(address _affiliateFeeManager) returns()
func (_Relay *RelayTransactor) SetAffiliateFeeManager(opts *bind.TransactOpts, _affiliateFeeManager common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setAffiliateFeeManager", _affiliateFeeManager)
}

// SetAffiliateFeeManager is a paid mutator transaction binding the contract method 0x3fd77557.
//
// Solidity: function setAffiliateFeeManager(address _affiliateFeeManager) returns()
func (_Relay *RelaySession) SetAffiliateFeeManager(_affiliateFeeManager common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetAffiliateFeeManager(&_Relay.TransactOpts, _affiliateFeeManager)
}

// SetAffiliateFeeManager is a paid mutator transaction binding the contract method 0x3fd77557.
//
// Solidity: function setAffiliateFeeManager(address _affiliateFeeManager) returns()
func (_Relay *RelayTransactorSession) SetAffiliateFeeManager(_affiliateFeeManager common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetAffiliateFeeManager(&_Relay.TransactOpts, _affiliateFeeManager)
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

// SetPeriphery is a paid mutator transaction binding the contract method 0xaeb22934.
//
// Solidity: function setPeriphery(address _periphery) returns()
func (_Relay *RelayTransactor) SetPeriphery(opts *bind.TransactOpts, _periphery common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setPeriphery", _periphery)
}

// SetPeriphery is a paid mutator transaction binding the contract method 0xaeb22934.
//
// Solidity: function setPeriphery(address _periphery) returns()
func (_Relay *RelaySession) SetPeriphery(_periphery common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetPeriphery(&_Relay.TransactOpts, _periphery)
}

// SetPeriphery is a paid mutator transaction binding the contract method 0xaeb22934.
//
// Solidity: function setPeriphery(address _periphery) returns()
func (_Relay *RelayTransactorSession) SetPeriphery(_periphery common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetPeriphery(&_Relay.TransactOpts, _periphery)
}

// SetSwap is a paid mutator transaction binding the contract method 0xacb3c073.
//
// Solidity: function setSwap(address _swap) returns()
func (_Relay *RelayTransactor) SetSwap(opts *bind.TransactOpts, _swap common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setSwap", _swap)
}

// SetSwap is a paid mutator transaction binding the contract method 0xacb3c073.
//
// Solidity: function setSwap(address _swap) returns()
func (_Relay *RelaySession) SetSwap(_swap common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetSwap(&_Relay.TransactOpts, _swap)
}

// SetSwap is a paid mutator transaction binding the contract method 0xacb3c073.
//
// Solidity: function setSwap(address _swap) returns()
func (_Relay *RelayTransactorSession) SetSwap(_swap common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetSwap(&_Relay.TransactOpts, _swap)
}

// SetVaultManager is a paid mutator transaction binding the contract method 0xb543503e.
//
// Solidity: function setVaultManager(address _vaultManager) returns()
func (_Relay *RelayTransactor) SetVaultManager(opts *bind.TransactOpts, _vaultManager common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setVaultManager", _vaultManager)
}

// SetVaultManager is a paid mutator transaction binding the contract method 0xb543503e.
//
// Solidity: function setVaultManager(address _vaultManager) returns()
func (_Relay *RelaySession) SetVaultManager(_vaultManager common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetVaultManager(&_Relay.TransactOpts, _vaultManager)
}

// SetVaultManager is a paid mutator transaction binding the contract method 0xb543503e.
//
// Solidity: function setVaultManager(address _vaultManager) returns()
func (_Relay *RelayTransactorSession) SetVaultManager(_vaultManager common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetVaultManager(&_Relay.TransactOpts, _vaultManager)
}

// SetWtoken is a paid mutator transaction binding the contract method 0x502843e7.
//
// Solidity: function setWtoken(address _wToken) returns()
func (_Relay *RelayTransactor) SetWtoken(opts *bind.TransactOpts, _wToken common.Address) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "setWtoken", _wToken)
}

// SetWtoken is a paid mutator transaction binding the contract method 0x502843e7.
//
// Solidity: function setWtoken(address _wToken) returns()
func (_Relay *RelaySession) SetWtoken(_wToken common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetWtoken(&_Relay.TransactOpts, _wToken)
}

// SetWtoken is a paid mutator transaction binding the contract method 0x502843e7.
//
// Solidity: function setWtoken(address _wToken) returns()
func (_Relay *RelayTransactorSession) SetWtoken(_wToken common.Address) (*types.Transaction, error) {
	return _Relay.Contract.SetWtoken(&_Relay.TransactOpts, _wToken)
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

// UpdateTokens is a paid mutator transaction binding the contract method 0xc8bc5fb0.
//
// Solidity: function updateTokens(address[] _tokens, uint256 _feature) returns()
func (_Relay *RelayTransactor) UpdateTokens(opts *bind.TransactOpts, _tokens []common.Address, _feature *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "updateTokens", _tokens, _feature)
}

// UpdateTokens is a paid mutator transaction binding the contract method 0xc8bc5fb0.
//
// Solidity: function updateTokens(address[] _tokens, uint256 _feature) returns()
func (_Relay *RelaySession) UpdateTokens(_tokens []common.Address, _feature *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.UpdateTokens(&_Relay.TransactOpts, _tokens, _feature)
}

// UpdateTokens is a paid mutator transaction binding the contract method 0xc8bc5fb0.
//
// Solidity: function updateTokens(address[] _tokens, uint256 _feature) returns()
func (_Relay *RelayTransactorSession) UpdateTokens(_tokens []common.Address, _feature *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.UpdateTokens(&_Relay.TransactOpts, _tokens, _feature)
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

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _vaultToken, uint256 _vaultAmount) returns()
func (_Relay *RelayTransactor) Withdraw(opts *bind.TransactOpts, _vaultToken common.Address, _vaultAmount *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "withdraw", _vaultToken, _vaultAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _vaultToken, uint256 _vaultAmount) returns()
func (_Relay *RelaySession) Withdraw(_vaultToken common.Address, _vaultAmount *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Withdraw(&_Relay.TransactOpts, _vaultToken, _vaultAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _vaultToken, uint256 _vaultAmount) returns()
func (_Relay *RelayTransactorSession) Withdraw(_vaultToken common.Address, _vaultAmount *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.Withdraw(&_Relay.TransactOpts, _vaultToken, _vaultAmount)
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

// RelayBridgeCompletedIterator is returned from FilterBridgeCompleted and is used to iterate over the raw logs and unpacked data for BridgeCompleted events raised by the Relay contract.
type RelayBridgeCompletedIterator struct {
	Event *RelayBridgeCompleted // Event containing the contract specifics and raw log

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
func (it *RelayBridgeCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeCompleted)
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
		it.Event = new(RelayBridgeCompleted)
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
func (it *RelayBridgeCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeCompleted represents a BridgeCompleted event raised by the Relay contract.
type RelayBridgeCompleted struct {
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Vault            []byte
	Sequence         *big.Int
	Sender           common.Address
	Data             []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBridgeCompleted is a free log retrieval operation binding the contract event 0x298a40641bd31f72c733761e0e85a6bd8a36909666ac2ed63a42c8015d025638.
//
// Solidity: event BridgeCompleted(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, uint256 sequence, address sender, bytes data)
func (_Relay *RelayFilterer) FilterBridgeCompleted(opts *bind.FilterOpts, orderId [][32]byte, chainAndGasLimit []*big.Int) (*RelayBridgeCompletedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeCompleted", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeCompletedIterator{contract: _Relay.contract, event: "BridgeCompleted", logs: logs, sub: sub}, nil
}

// WatchBridgeCompleted is a free log subscription operation binding the contract event 0x298a40641bd31f72c733761e0e85a6bd8a36909666ac2ed63a42c8015d025638.
//
// Solidity: event BridgeCompleted(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, uint256 sequence, address sender, bytes data)
func (_Relay *RelayFilterer) WatchBridgeCompleted(opts *bind.WatchOpts, sink chan<- *RelayBridgeCompleted, orderId [][32]byte, chainAndGasLimit []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeCompleted", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeCompleted)
				if err := _Relay.contract.UnpackLog(event, "BridgeCompleted", log); err != nil {
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

// ParseBridgeCompleted is a log parse operation binding the contract event 0x298a40641bd31f72c733761e0e85a6bd8a36909666ac2ed63a42c8015d025638.
//
// Solidity: event BridgeCompleted(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, uint256 sequence, address sender, bytes data)
func (_Relay *RelayFilterer) ParseBridgeCompleted(log types.Log) (*RelayBridgeCompleted, error) {
	event := new(RelayBridgeCompleted)
	if err := _Relay.contract.UnpackLog(event, "BridgeCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeFailedIterator is returned from FilterBridgeFailed and is used to iterate over the raw logs and unpacked data for BridgeFailed events raised by the Relay contract.
type RelayBridgeFailedIterator struct {
	Event *RelayBridgeFailed // Event containing the contract specifics and raw log

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
func (it *RelayBridgeFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeFailed)
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
		it.Event = new(RelayBridgeFailed)
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
func (it *RelayBridgeFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeFailed represents a BridgeFailed event raised by the Relay contract.
type RelayBridgeFailed struct {
	OrderId [32]byte
	Token   common.Address
	Amount  *big.Int
	From    []byte
	To      common.Address
	Data    []byte
	Reason  []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBridgeFailed is a free log retrieval operation binding the contract event 0x9b1e49b4414cebdf5b9320940891509783f35ca9ce8bb64156f913cadc7d4c53.
//
// Solidity: event BridgeFailed(bytes32 indexed orderId, address token, uint256 amount, bytes from, address to, bytes data, bytes reason)
func (_Relay *RelayFilterer) FilterBridgeFailed(opts *bind.FilterOpts, orderId [][32]byte) (*RelayBridgeFailedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeFailed", orderIdRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeFailedIterator{contract: _Relay.contract, event: "BridgeFailed", logs: logs, sub: sub}, nil
}

// WatchBridgeFailed is a free log subscription operation binding the contract event 0x9b1e49b4414cebdf5b9320940891509783f35ca9ce8bb64156f913cadc7d4c53.
//
// Solidity: event BridgeFailed(bytes32 indexed orderId, address token, uint256 amount, bytes from, address to, bytes data, bytes reason)
func (_Relay *RelayFilterer) WatchBridgeFailed(opts *bind.WatchOpts, sink chan<- *RelayBridgeFailed, orderId [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeFailed", orderIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeFailed)
				if err := _Relay.contract.UnpackLog(event, "BridgeFailed", log); err != nil {
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

// ParseBridgeFailed is a log parse operation binding the contract event 0x9b1e49b4414cebdf5b9320940891509783f35ca9ce8bb64156f913cadc7d4c53.
//
// Solidity: event BridgeFailed(bytes32 indexed orderId, address token, uint256 amount, bytes from, address to, bytes data, bytes reason)
func (_Relay *RelayFilterer) ParseBridgeFailed(log types.Log) (*RelayBridgeFailed, error) {
	event := new(RelayBridgeFailed)
	if err := _Relay.contract.UnpackLog(event, "BridgeFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeFeeCollectedIterator is returned from FilterBridgeFeeCollected and is used to iterate over the raw logs and unpacked data for BridgeFeeCollected events raised by the Relay contract.
type RelayBridgeFeeCollectedIterator struct {
	Event *RelayBridgeFeeCollected // Event containing the contract specifics and raw log

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
func (it *RelayBridgeFeeCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeFeeCollected)
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
		it.Event = new(RelayBridgeFeeCollected)
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
func (it *RelayBridgeFeeCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeFeeCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeFeeCollected represents a BridgeFeeCollected event raised by the Relay contract.
type RelayBridgeFeeCollected struct {
	OrderId     [32]byte
	Token       common.Address
	SecurityFee *big.Int
	VaultFee    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBridgeFeeCollected is a free log retrieval operation binding the contract event 0x90ac3f4e75a9d001d4b672a79bf7f25da05b6748856a0aeb61c7a840f89039ac.
//
// Solidity: event BridgeFeeCollected(bytes32 indexed orderId, address token, uint256 securityFee, uint256 vaultFee)
func (_Relay *RelayFilterer) FilterBridgeFeeCollected(opts *bind.FilterOpts, orderId [][32]byte) (*RelayBridgeFeeCollectedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeFeeCollected", orderIdRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeFeeCollectedIterator{contract: _Relay.contract, event: "BridgeFeeCollected", logs: logs, sub: sub}, nil
}

// WatchBridgeFeeCollected is a free log subscription operation binding the contract event 0x90ac3f4e75a9d001d4b672a79bf7f25da05b6748856a0aeb61c7a840f89039ac.
//
// Solidity: event BridgeFeeCollected(bytes32 indexed orderId, address token, uint256 securityFee, uint256 vaultFee)
func (_Relay *RelayFilterer) WatchBridgeFeeCollected(opts *bind.WatchOpts, sink chan<- *RelayBridgeFeeCollected, orderId [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeFeeCollected", orderIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeFeeCollected)
				if err := _Relay.contract.UnpackLog(event, "BridgeFeeCollected", log); err != nil {
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

// ParseBridgeFeeCollected is a log parse operation binding the contract event 0x90ac3f4e75a9d001d4b672a79bf7f25da05b6748856a0aeb61c7a840f89039ac.
//
// Solidity: event BridgeFeeCollected(bytes32 indexed orderId, address token, uint256 securityFee, uint256 vaultFee)
func (_Relay *RelayFilterer) ParseBridgeFeeCollected(log types.Log) (*RelayBridgeFeeCollected, error) {
	event := new(RelayBridgeFeeCollected)
	if err := _Relay.contract.UnpackLog(event, "BridgeFeeCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeInIterator is returned from FilterBridgeIn and is used to iterate over the raw logs and unpacked data for BridgeIn events raised by the Relay contract.
type RelayBridgeInIterator struct {
	Event *RelayBridgeIn // Event containing the contract specifics and raw log

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
func (it *RelayBridgeInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeIn)
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
		it.Event = new(RelayBridgeIn)
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
func (it *RelayBridgeInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeIn represents a BridgeIn event raised by the Relay contract.
type RelayBridgeIn struct {
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	TxInType         uint8
	Vault            []byte
	Sequence         *big.Int
	Sender           common.Address
	Token            common.Address
	Amount           *big.Int
	To               common.Address
	Data             []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBridgeIn is a free log retrieval operation binding the contract event 0x8104943fdd0997a3240b59b381251572ac6ac81941e1af29845de70edca938a4.
//
// Solidity: event BridgeIn(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txInType, bytes vault, uint256 sequence, address sender, address token, uint256 amount, address to, bytes data)
func (_Relay *RelayFilterer) FilterBridgeIn(opts *bind.FilterOpts, orderId [][32]byte, chainAndGasLimit []*big.Int) (*RelayBridgeInIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeIn", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeInIterator{contract: _Relay.contract, event: "BridgeIn", logs: logs, sub: sub}, nil
}

// WatchBridgeIn is a free log subscription operation binding the contract event 0x8104943fdd0997a3240b59b381251572ac6ac81941e1af29845de70edca938a4.
//
// Solidity: event BridgeIn(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txInType, bytes vault, uint256 sequence, address sender, address token, uint256 amount, address to, bytes data)
func (_Relay *RelayFilterer) WatchBridgeIn(opts *bind.WatchOpts, sink chan<- *RelayBridgeIn, orderId [][32]byte, chainAndGasLimit []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeIn", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeIn)
				if err := _Relay.contract.UnpackLog(event, "BridgeIn", log); err != nil {
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

// ParseBridgeIn is a log parse operation binding the contract event 0x8104943fdd0997a3240b59b381251572ac6ac81941e1af29845de70edca938a4.
//
// Solidity: event BridgeIn(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txInType, bytes vault, uint256 sequence, address sender, address token, uint256 amount, address to, bytes data)
func (_Relay *RelayFilterer) ParseBridgeIn(log types.Log) (*RelayBridgeIn, error) {
	event := new(RelayBridgeIn)
	if err := _Relay.contract.UnpackLog(event, "BridgeIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeOutIterator is returned from FilterBridgeOut and is used to iterate over the raw logs and unpacked data for BridgeOut events raised by the Relay contract.
type RelayBridgeOutIterator struct {
	Event *RelayBridgeOut // Event containing the contract specifics and raw log

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
func (it *RelayBridgeOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeOut)
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
		it.Event = new(RelayBridgeOut)
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
func (it *RelayBridgeOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeOut represents a BridgeOut event raised by the Relay contract.
type RelayBridgeOut struct {
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	TxOutType        uint8
	Vault            []byte
	Token            common.Address
	Amount           *big.Int
	From             common.Address
	RefundAddr       common.Address
	To               []byte
	Payload          []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBridgeOut is a free log retrieval operation binding the contract event 0x70e58c596a5c2186f7e5546898d023e26870b3c0c19cabc29b5c86dfd0690a2f.
//
// Solidity: event BridgeOut(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, address token, uint256 amount, address from, address refundAddr, bytes to, bytes payload)
func (_Relay *RelayFilterer) FilterBridgeOut(opts *bind.FilterOpts, orderId [][32]byte, chainAndGasLimit []*big.Int) (*RelayBridgeOutIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeOut", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeOutIterator{contract: _Relay.contract, event: "BridgeOut", logs: logs, sub: sub}, nil
}

// WatchBridgeOut is a free log subscription operation binding the contract event 0x70e58c596a5c2186f7e5546898d023e26870b3c0c19cabc29b5c86dfd0690a2f.
//
// Solidity: event BridgeOut(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, address token, uint256 amount, address from, address refundAddr, bytes to, bytes payload)
func (_Relay *RelayFilterer) WatchBridgeOut(opts *bind.WatchOpts, sink chan<- *RelayBridgeOut, orderId [][32]byte, chainAndGasLimit []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeOut", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeOut)
				if err := _Relay.contract.UnpackLog(event, "BridgeOut", log); err != nil {
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

// ParseBridgeOut is a log parse operation binding the contract event 0x70e58c596a5c2186f7e5546898d023e26870b3c0c19cabc29b5c86dfd0690a2f.
//
// Solidity: event BridgeOut(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txOutType, bytes vault, address token, uint256 amount, address from, address refundAddr, bytes to, bytes payload)
func (_Relay *RelayFilterer) ParseBridgeOut(log types.Log) (*RelayBridgeOut, error) {
	event := new(RelayBridgeOut)
	if err := _Relay.contract.UnpackLog(event, "BridgeOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeRelayIterator is returned from FilterBridgeRelay and is used to iterate over the raw logs and unpacked data for BridgeRelay events raised by the Relay contract.
type RelayBridgeRelayIterator struct {
	Event *RelayBridgeRelay // Event containing the contract specifics and raw log

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
func (it *RelayBridgeRelayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeRelay)
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
		it.Event = new(RelayBridgeRelay)
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
func (it *RelayBridgeRelayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeRelayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeRelay represents a BridgeRelay event raised by the Relay contract.
type RelayBridgeRelay struct {
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	TxType           uint8
	Vault            []byte
	To               []byte
	Token            []byte
	Amount           *big.Int
	Sequence         *big.Int
	Hash             [32]byte
	From             []byte
	Data             []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBridgeRelay is a free log retrieval operation binding the contract event 0x2b44f5da40771b7770f8469714d202f01aa69c410bd79a5dd86f9562413c3fca.
//
// Solidity: event BridgeRelay(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txType, bytes vault, bytes to, bytes token, uint256 amount, uint256 sequence, bytes32 hash, bytes from, bytes data)
func (_Relay *RelayFilterer) FilterBridgeRelay(opts *bind.FilterOpts, orderId [][32]byte, chainAndGasLimit []*big.Int) (*RelayBridgeRelayIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeRelay", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeRelayIterator{contract: _Relay.contract, event: "BridgeRelay", logs: logs, sub: sub}, nil
}

// WatchBridgeRelay is a free log subscription operation binding the contract event 0x2b44f5da40771b7770f8469714d202f01aa69c410bd79a5dd86f9562413c3fca.
//
// Solidity: event BridgeRelay(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txType, bytes vault, bytes to, bytes token, uint256 amount, uint256 sequence, bytes32 hash, bytes from, bytes data)
func (_Relay *RelayFilterer) WatchBridgeRelay(opts *bind.WatchOpts, sink chan<- *RelayBridgeRelay, orderId [][32]byte, chainAndGasLimit []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeRelay", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeRelay)
				if err := _Relay.contract.UnpackLog(event, "BridgeRelay", log); err != nil {
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

// ParseBridgeRelay is a log parse operation binding the contract event 0x2b44f5da40771b7770f8469714d202f01aa69c410bd79a5dd86f9562413c3fca.
//
// Solidity: event BridgeRelay(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, uint8 txType, bytes vault, bytes to, bytes token, uint256 amount, uint256 sequence, bytes32 hash, bytes from, bytes data)
func (_Relay *RelayFilterer) ParseBridgeRelay(log types.Log) (*RelayBridgeRelay, error) {
	event := new(RelayBridgeRelay)
	if err := _Relay.contract.UnpackLog(event, "BridgeRelay", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayBridgeRelaySignedIterator is returned from FilterBridgeRelaySigned and is used to iterate over the raw logs and unpacked data for BridgeRelaySigned events raised by the Relay contract.
type RelayBridgeRelaySignedIterator struct {
	Event *RelayBridgeRelaySigned // Event containing the contract specifics and raw log

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
func (it *RelayBridgeRelaySignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayBridgeRelaySigned)
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
		it.Event = new(RelayBridgeRelaySigned)
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
func (it *RelayBridgeRelaySignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayBridgeRelaySignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayBridgeRelaySigned represents a BridgeRelaySigned event raised by the Relay contract.
type RelayBridgeRelaySigned struct {
	OrderId          [32]byte
	ChainAndGasLimit *big.Int
	Vault            []byte
	RelayData        []byte
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBridgeRelaySigned is a free log retrieval operation binding the contract event 0x0bbc5d146426ff8f68e7bcd94bd95e273fabd90609e118badb81428405e38107.
//
// Solidity: event BridgeRelaySigned(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, bytes vault, bytes relayData, bytes signature)
func (_Relay *RelayFilterer) FilterBridgeRelaySigned(opts *bind.FilterOpts, orderId [][32]byte, chainAndGasLimit []*big.Int) (*RelayBridgeRelaySignedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "BridgeRelaySigned", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return &RelayBridgeRelaySignedIterator{contract: _Relay.contract, event: "BridgeRelaySigned", logs: logs, sub: sub}, nil
}

// WatchBridgeRelaySigned is a free log subscription operation binding the contract event 0x0bbc5d146426ff8f68e7bcd94bd95e273fabd90609e118badb81428405e38107.
//
// Solidity: event BridgeRelaySigned(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, bytes vault, bytes relayData, bytes signature)
func (_Relay *RelayFilterer) WatchBridgeRelaySigned(opts *bind.WatchOpts, sink chan<- *RelayBridgeRelaySigned, orderId [][32]byte, chainAndGasLimit []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var chainAndGasLimitRule []interface{}
	for _, chainAndGasLimitItem := range chainAndGasLimit {
		chainAndGasLimitRule = append(chainAndGasLimitRule, chainAndGasLimitItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "BridgeRelaySigned", orderIdRule, chainAndGasLimitRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayBridgeRelaySigned)
				if err := _Relay.contract.UnpackLog(event, "BridgeRelaySigned", log); err != nil {
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

// ParseBridgeRelaySigned is a log parse operation binding the contract event 0x0bbc5d146426ff8f68e7bcd94bd95e273fabd90609e118badb81428405e38107.
//
// Solidity: event BridgeRelaySigned(bytes32 indexed orderId, uint256 indexed chainAndGasLimit, bytes vault, bytes relayData, bytes signature)
func (_Relay *RelayFilterer) ParseBridgeRelaySigned(log types.Log) (*RelayBridgeRelaySigned, error) {
	event := new(RelayBridgeRelaySigned)
	if err := _Relay.contract.UnpackLog(event, "BridgeRelaySigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Relay contract.
type RelayDepositIterator struct {
	Event *RelayDeposit // Event containing the contract specifics and raw log

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
func (it *RelayDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayDeposit)
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
		it.Event = new(RelayDeposit)
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
func (it *RelayDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayDeposit represents a Deposit event raised by the Relay contract.
type RelayDeposit struct {
	OrderId   [32]byte
	FromChain *big.Int
	Token     common.Address
	Amount    *big.Int
	To        common.Address
	From      []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x3d57f75a1f744f3d68cdb6feff30d8e091336afa058aa79931b4648c13ad3cd4.
//
// Solidity: event Deposit(bytes32 orderId, uint256 fromChain, address token, uint256 amount, address to, bytes from)
func (_Relay *RelayFilterer) FilterDeposit(opts *bind.FilterOpts) (*RelayDepositIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &RelayDepositIterator{contract: _Relay.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x3d57f75a1f744f3d68cdb6feff30d8e091336afa058aa79931b4648c13ad3cd4.
//
// Solidity: event Deposit(bytes32 orderId, uint256 fromChain, address token, uint256 amount, address to, bytes from)
func (_Relay *RelayFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RelayDeposit) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayDeposit)
				if err := _Relay.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x3d57f75a1f744f3d68cdb6feff30d8e091336afa058aa79931b4648c13ad3cd4.
//
// Solidity: event Deposit(bytes32 orderId, uint256 fromChain, address token, uint256 amount, address to, bytes from)
func (_Relay *RelayFilterer) ParseDeposit(log types.Log) (*RelayDeposit, error) {
	event := new(RelayDeposit)
	if err := _Relay.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// RelaySetAffiliateFeeManagerIterator is returned from FilterSetAffiliateFeeManager and is used to iterate over the raw logs and unpacked data for SetAffiliateFeeManager events raised by the Relay contract.
type RelaySetAffiliateFeeManagerIterator struct {
	Event *RelaySetAffiliateFeeManager // Event containing the contract specifics and raw log

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
func (it *RelaySetAffiliateFeeManagerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySetAffiliateFeeManager)
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
		it.Event = new(RelaySetAffiliateFeeManager)
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
func (it *RelaySetAffiliateFeeManagerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetAffiliateFeeManagerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySetAffiliateFeeManager represents a SetAffiliateFeeManager event raised by the Relay contract.
type RelaySetAffiliateFeeManager struct {
	AffiliateFeeManager common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetAffiliateFeeManager is a free log retrieval operation binding the contract event 0x91e20f942fdf0e8f25581b4dbed114ec7479cc0568fb916fe209f0f594fd4d37.
//
// Solidity: event SetAffiliateFeeManager(address _affiliateFeeManager)
func (_Relay *RelayFilterer) FilterSetAffiliateFeeManager(opts *bind.FilterOpts) (*RelaySetAffiliateFeeManagerIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "SetAffiliateFeeManager")
	if err != nil {
		return nil, err
	}
	return &RelaySetAffiliateFeeManagerIterator{contract: _Relay.contract, event: "SetAffiliateFeeManager", logs: logs, sub: sub}, nil
}

// WatchSetAffiliateFeeManager is a free log subscription operation binding the contract event 0x91e20f942fdf0e8f25581b4dbed114ec7479cc0568fb916fe209f0f594fd4d37.
//
// Solidity: event SetAffiliateFeeManager(address _affiliateFeeManager)
func (_Relay *RelayFilterer) WatchSetAffiliateFeeManager(opts *bind.WatchOpts, sink chan<- *RelaySetAffiliateFeeManager) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "SetAffiliateFeeManager")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySetAffiliateFeeManager)
				if err := _Relay.contract.UnpackLog(event, "SetAffiliateFeeManager", log); err != nil {
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

// ParseSetAffiliateFeeManager is a log parse operation binding the contract event 0x91e20f942fdf0e8f25581b4dbed114ec7479cc0568fb916fe209f0f594fd4d37.
//
// Solidity: event SetAffiliateFeeManager(address _affiliateFeeManager)
func (_Relay *RelayFilterer) ParseSetAffiliateFeeManager(log types.Log) (*RelaySetAffiliateFeeManager, error) {
	event := new(RelaySetAffiliateFeeManager)
	if err := _Relay.contract.UnpackLog(event, "SetAffiliateFeeManager", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelaySetPeripheryIterator is returned from FilterSetPeriphery and is used to iterate over the raw logs and unpacked data for SetPeriphery events raised by the Relay contract.
type RelaySetPeripheryIterator struct {
	Event *RelaySetPeriphery // Event containing the contract specifics and raw log

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
func (it *RelaySetPeripheryIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySetPeriphery)
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
		it.Event = new(RelaySetPeriphery)
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
func (it *RelaySetPeripheryIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetPeripheryIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySetPeriphery represents a SetPeriphery event raised by the Relay contract.
type RelaySetPeriphery struct {
	Periphery common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetPeriphery is a free log retrieval operation binding the contract event 0xb62d2ad211763879f4b426c82f30fda8ef562b3aedf50e09ae9aa33dbf4b82f3.
//
// Solidity: event SetPeriphery(address _periphery)
func (_Relay *RelayFilterer) FilterSetPeriphery(opts *bind.FilterOpts) (*RelaySetPeripheryIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "SetPeriphery")
	if err != nil {
		return nil, err
	}
	return &RelaySetPeripheryIterator{contract: _Relay.contract, event: "SetPeriphery", logs: logs, sub: sub}, nil
}

// WatchSetPeriphery is a free log subscription operation binding the contract event 0xb62d2ad211763879f4b426c82f30fda8ef562b3aedf50e09ae9aa33dbf4b82f3.
//
// Solidity: event SetPeriphery(address _periphery)
func (_Relay *RelayFilterer) WatchSetPeriphery(opts *bind.WatchOpts, sink chan<- *RelaySetPeriphery) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "SetPeriphery")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySetPeriphery)
				if err := _Relay.contract.UnpackLog(event, "SetPeriphery", log); err != nil {
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

// ParseSetPeriphery is a log parse operation binding the contract event 0xb62d2ad211763879f4b426c82f30fda8ef562b3aedf50e09ae9aa33dbf4b82f3.
//
// Solidity: event SetPeriphery(address _periphery)
func (_Relay *RelayFilterer) ParseSetPeriphery(log types.Log) (*RelaySetPeriphery, error) {
	event := new(RelaySetPeriphery)
	if err := _Relay.contract.UnpackLog(event, "SetPeriphery", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelaySetSwapIterator is returned from FilterSetSwap and is used to iterate over the raw logs and unpacked data for SetSwap events raised by the Relay contract.
type RelaySetSwapIterator struct {
	Event *RelaySetSwap // Event containing the contract specifics and raw log

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
func (it *RelaySetSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySetSwap)
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
		it.Event = new(RelaySetSwap)
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
func (it *RelaySetSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySetSwap represents a SetSwap event raised by the Relay contract.
type RelaySetSwap struct {
	Swap common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetSwap is a free log retrieval operation binding the contract event 0x0640d773eb4eb9fe289b6fbb8bf82f5418820ce8a52592a70358dba1a9781655.
//
// Solidity: event SetSwap(address _swap)
func (_Relay *RelayFilterer) FilterSetSwap(opts *bind.FilterOpts) (*RelaySetSwapIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "SetSwap")
	if err != nil {
		return nil, err
	}
	return &RelaySetSwapIterator{contract: _Relay.contract, event: "SetSwap", logs: logs, sub: sub}, nil
}

// WatchSetSwap is a free log subscription operation binding the contract event 0x0640d773eb4eb9fe289b6fbb8bf82f5418820ce8a52592a70358dba1a9781655.
//
// Solidity: event SetSwap(address _swap)
func (_Relay *RelayFilterer) WatchSetSwap(opts *bind.WatchOpts, sink chan<- *RelaySetSwap) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "SetSwap")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySetSwap)
				if err := _Relay.contract.UnpackLog(event, "SetSwap", log); err != nil {
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

// ParseSetSwap is a log parse operation binding the contract event 0x0640d773eb4eb9fe289b6fbb8bf82f5418820ce8a52592a70358dba1a9781655.
//
// Solidity: event SetSwap(address _swap)
func (_Relay *RelayFilterer) ParseSetSwap(log types.Log) (*RelaySetSwap, error) {
	event := new(RelaySetSwap)
	if err := _Relay.contract.UnpackLog(event, "SetSwap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelaySetVaultManagerIterator is returned from FilterSetVaultManager and is used to iterate over the raw logs and unpacked data for SetVaultManager events raised by the Relay contract.
type RelaySetVaultManagerIterator struct {
	Event *RelaySetVaultManager // Event containing the contract specifics and raw log

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
func (it *RelaySetVaultManagerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySetVaultManager)
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
		it.Event = new(RelaySetVaultManager)
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
func (it *RelaySetVaultManagerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetVaultManagerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySetVaultManager represents a SetVaultManager event raised by the Relay contract.
type RelaySetVaultManager struct {
	VaultManager common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSetVaultManager is a free log retrieval operation binding the contract event 0xf936751e3ffe55eec34636f5166fff2a5714505f4a3aa0862829ddb1eb75f101.
//
// Solidity: event SetVaultManager(address _vaultManager)
func (_Relay *RelayFilterer) FilterSetVaultManager(opts *bind.FilterOpts) (*RelaySetVaultManagerIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "SetVaultManager")
	if err != nil {
		return nil, err
	}
	return &RelaySetVaultManagerIterator{contract: _Relay.contract, event: "SetVaultManager", logs: logs, sub: sub}, nil
}

// WatchSetVaultManager is a free log subscription operation binding the contract event 0xf936751e3ffe55eec34636f5166fff2a5714505f4a3aa0862829ddb1eb75f101.
//
// Solidity: event SetVaultManager(address _vaultManager)
func (_Relay *RelayFilterer) WatchSetVaultManager(opts *bind.WatchOpts, sink chan<- *RelaySetVaultManager) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "SetVaultManager")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySetVaultManager)
				if err := _Relay.contract.UnpackLog(event, "SetVaultManager", log); err != nil {
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

// ParseSetVaultManager is a log parse operation binding the contract event 0xf936751e3ffe55eec34636f5166fff2a5714505f4a3aa0862829ddb1eb75f101.
//
// Solidity: event SetVaultManager(address _vaultManager)
func (_Relay *RelayFilterer) ParseSetVaultManager(log types.Log) (*RelaySetVaultManager, error) {
	event := new(RelaySetVaultManager)
	if err := _Relay.contract.UnpackLog(event, "SetVaultManager", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelaySetWTokenIterator is returned from FilterSetWToken and is used to iterate over the raw logs and unpacked data for SetWToken events raised by the Relay contract.
type RelaySetWTokenIterator struct {
	Event *RelaySetWToken // Event containing the contract specifics and raw log

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
func (it *RelaySetWTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelaySetWToken)
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
		it.Event = new(RelaySetWToken)
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
func (it *RelaySetWTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelaySetWTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelaySetWToken represents a SetWToken event raised by the Relay contract.
type RelaySetWToken struct {
	WToken common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSetWToken is a free log retrieval operation binding the contract event 0x84eb9082159e970b869283d8b9070a08ca1769321bce2f3dbfbee2c73bf94bed.
//
// Solidity: event SetWToken(address _wToken)
func (_Relay *RelayFilterer) FilterSetWToken(opts *bind.FilterOpts) (*RelaySetWTokenIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "SetWToken")
	if err != nil {
		return nil, err
	}
	return &RelaySetWTokenIterator{contract: _Relay.contract, event: "SetWToken", logs: logs, sub: sub}, nil
}

// WatchSetWToken is a free log subscription operation binding the contract event 0x84eb9082159e970b869283d8b9070a08ca1769321bce2f3dbfbee2c73bf94bed.
//
// Solidity: event SetWToken(address _wToken)
func (_Relay *RelayFilterer) WatchSetWToken(opts *bind.WatchOpts, sink chan<- *RelaySetWToken) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "SetWToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelaySetWToken)
				if err := _Relay.contract.UnpackLog(event, "SetWToken", log); err != nil {
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

// ParseSetWToken is a log parse operation binding the contract event 0x84eb9082159e970b869283d8b9070a08ca1769321bce2f3dbfbee2c73bf94bed.
//
// Solidity: event SetWToken(address _wToken)
func (_Relay *RelayFilterer) ParseSetWToken(log types.Log) (*RelaySetWToken, error) {
	event := new(RelaySetWToken)
	if err := _Relay.contract.UnpackLog(event, "SetWToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RelayTransferInIterator is returned from FilterTransferIn and is used to iterate over the raw logs and unpacked data for TransferIn events raised by the Relay contract.
type RelayTransferInIterator struct {
	Event *RelayTransferIn // Event containing the contract specifics and raw log

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
func (it *RelayTransferInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayTransferIn)
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
		it.Event = new(RelayTransferIn)
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
func (it *RelayTransferInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayTransferInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayTransferIn represents a TransferIn event raised by the Relay contract.
type RelayTransferIn struct {
	OrderId [32]byte
	Token   common.Address
	Amount  *big.Int
	To      common.Address
	Result  bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferIn is a free log retrieval operation binding the contract event 0x9b5f743d5b09dd6daaf5ec5b369f6139e222351a3e4c2b4542043b0544ac35ae.
//
// Solidity: event TransferIn(bytes32 orderId, address token, uint256 amount, address to, bool result)
func (_Relay *RelayFilterer) FilterTransferIn(opts *bind.FilterOpts) (*RelayTransferInIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return &RelayTransferInIterator{contract: _Relay.contract, event: "TransferIn", logs: logs, sub: sub}, nil
}

// WatchTransferIn is a free log subscription operation binding the contract event 0x9b5f743d5b09dd6daaf5ec5b369f6139e222351a3e4c2b4542043b0544ac35ae.
//
// Solidity: event TransferIn(bytes32 orderId, address token, uint256 amount, address to, bool result)
func (_Relay *RelayFilterer) WatchTransferIn(opts *bind.WatchOpts, sink chan<- *RelayTransferIn) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayTransferIn)
				if err := _Relay.contract.UnpackLog(event, "TransferIn", log); err != nil {
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

// ParseTransferIn is a log parse operation binding the contract event 0x9b5f743d5b09dd6daaf5ec5b369f6139e222351a3e4c2b4542043b0544ac35ae.
//
// Solidity: event TransferIn(bytes32 orderId, address token, uint256 amount, address to, bool result)
func (_Relay *RelayFilterer) ParseTransferIn(log types.Log) (*RelayTransferIn, error) {
	event := new(RelayTransferIn)
	if err := _Relay.contract.UnpackLog(event, "TransferIn", log); err != nil {
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

// RelayUpdateTokensIterator is returned from FilterUpdateTokens and is used to iterate over the raw logs and unpacked data for UpdateTokens events raised by the Relay contract.
type RelayUpdateTokensIterator struct {
	Event *RelayUpdateTokens // Event containing the contract specifics and raw log

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
func (it *RelayUpdateTokensIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayUpdateTokens)
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
		it.Event = new(RelayUpdateTokens)
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
func (it *RelayUpdateTokensIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayUpdateTokensIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayUpdateTokens represents a UpdateTokens event raised by the Relay contract.
type RelayUpdateTokens struct {
	Token   common.Address
	Feature *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdateTokens is a free log retrieval operation binding the contract event 0xe4e40aeaf1c91c8af9cf0a8aa8ea7c0a889ef4057c409e4d14fadd0526e685aa.
//
// Solidity: event UpdateTokens(address token, uint256 feature)
func (_Relay *RelayFilterer) FilterUpdateTokens(opts *bind.FilterOpts) (*RelayUpdateTokensIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "UpdateTokens")
	if err != nil {
		return nil, err
	}
	return &RelayUpdateTokensIterator{contract: _Relay.contract, event: "UpdateTokens", logs: logs, sub: sub}, nil
}

// WatchUpdateTokens is a free log subscription operation binding the contract event 0xe4e40aeaf1c91c8af9cf0a8aa8ea7c0a889ef4057c409e4d14fadd0526e685aa.
//
// Solidity: event UpdateTokens(address token, uint256 feature)
func (_Relay *RelayFilterer) WatchUpdateTokens(opts *bind.WatchOpts, sink chan<- *RelayUpdateTokens) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "UpdateTokens")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayUpdateTokens)
				if err := _Relay.contract.UnpackLog(event, "UpdateTokens", log); err != nil {
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

// ParseUpdateTokens is a log parse operation binding the contract event 0xe4e40aeaf1c91c8af9cf0a8aa8ea7c0a889ef4057c409e4d14fadd0526e685aa.
//
// Solidity: event UpdateTokens(address token, uint256 feature)
func (_Relay *RelayFilterer) ParseUpdateTokens(log types.Log) (*RelayUpdateTokens, error) {
	event := new(RelayUpdateTokens)
	if err := _Relay.contract.UnpackLog(event, "UpdateTokens", log); err != nil {
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

// RelayWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Relay contract.
type RelayWithdrawIterator struct {
	Event *RelayWithdraw // Event containing the contract specifics and raw log

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
func (it *RelayWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayWithdraw)
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
		it.Event = new(RelayWithdraw)
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
func (it *RelayWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayWithdraw represents a Withdraw event raised by the Relay contract.
type RelayWithdraw struct {
	Token       common.Address
	Reicerver   common.Address
	VaultAmount *big.Int
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address token, address reicerver, uint256 vaultAmount, uint256 tokenAmount)
func (_Relay *RelayFilterer) FilterWithdraw(opts *bind.FilterOpts) (*RelayWithdrawIterator, error) {

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &RelayWithdrawIterator{contract: _Relay.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address token, address reicerver, uint256 vaultAmount, uint256 tokenAmount)
func (_Relay *RelayFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *RelayWithdraw) (event.Subscription, error) {

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayWithdraw)
				if err := _Relay.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xf341246adaac6f497bc2a656f546ab9e182111d630394f0c57c710a59a2cb567.
//
// Solidity: event Withdraw(address token, address reicerver, uint256 vaultAmount, uint256 tokenAmount)
func (_Relay *RelayFilterer) ParseWithdraw(log types.Log) (*RelayWithdraw, error) {
	event := new(RelayWithdraw)
	if err := _Relay.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
