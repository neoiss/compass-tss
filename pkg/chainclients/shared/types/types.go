package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/internal/structure"
	stypes "github.com/mapprotocol/compass-tss/x/types"
	"math/big"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/mapclient/types"
)

// ChainClient is the interface for chain clients.
type ChainClient interface {
	// Start starts the chain client with the given queues.
	Start(
		globalTxsQueue chan types.TxIn,
		globalErrataQueue chan types.ErrataBlock,
		globalSolvencyQueue chan types.Solvency,
		globalNetworkFeeQueue chan types.NetworkFee,
	)

	// Stop stops the chain client.
	Stop()

	// IsBlockScannerHealthy returns true if the block scanner is healthy.
	IsBlockScannerHealthy() bool

	// SignTx returns the signed transaction.
	SignTx(tx types.TxOutItem, height int64) ([]byte, []byte, *types.TxInItem, error)

	// BroadcastTx broadcasts the transaction and returns the transaction hash.
	BroadcastTx(_ types.TxOutItem, _ []byte) (string, error)

	// GetHeight returns the current height of the chain.
	GetHeight() (int64, error)

	// GetAddress returns the address for the given public key.
	GetAddress(poolPubKey common.PubKey) string

	// GetAccount returns the account for the given public key.
	GetAccount(poolPubKey common.PubKey, height *big.Int) (common.Account, error)

	// GetAccountByAddress returns the account for the given address.
	GetAccountByAddress(address string, height *big.Int) (common.Account, error)

	// GetChain returns the chain.
	GetChain() common.Chain

	// GetConfig returns the chain configuration.
	GetConfig() config.BifrostChainConfiguration

	// OnObservedTxIn is called when a new observed tx is received.
	OnObservedTxIn(txIn types.TxInItem, blockHeight int64)

	// GetConfirmationCount returns the confirmation count for the given tx.
	GetConfirmationCount(txIn types.TxIn) int64

	// ConfirmationCountReady returns true if the confirmation count is ready.
	ConfirmationCountReady(txIn types.TxIn) bool

	// GetBlockScannerHeight returns block scanner height for chain
	GetBlockScannerHeight() (int64, error)

	// GetLatestTxForVault returns last observed and broadcasted tx for a particular vault and chain
	GetLatestTxForVault(vault string) (string, string, error)
}

// SolvencyReporter reports the solvency of the chain at the given height.
type SolvencyReporter func(height int64) error

// Bridge is compass 2 map
type Bridge interface {
	EnsureNodeWhitelisted() error
	EnsureNodeWhitelistedWithTimeout() error
	FetchNodeStatus() (stypes.NodeStatus, error)
	FetchActiveNodes() ([]common.PubKey, error)
	GetNodeAccount(string) (*structure.MaintainerInfo, error)
	GetNodeAccounts(*big.Int) ([]structure.MaintainerInfo, error)
	GetKeygenBlock() (*structure.KeyGen, error)
	GetPools() (stypes.Pools, error) //
	GetVault(pubkey string) (stypes.Vault, error)
	GetPubKeys() ([]PubKeyContractAddressPair, error)
	GetContractAddress() ([]PubKeyContractAddressPair, error)
	GetLastObservedInHeight(chain common.Chain) (int64, error)
	GetLastSignedOutHeight(chain common.Chain) (int64, error)
	GetAsgards() (stypes.Vaults, error)
	GetConstants() (map[string]int64, error)
	GetMimir(key string) (int64, error)
	GetMimirWithRef(template, ref string) (int64, error)
	GetTHORName(name string) (stypes.THORName, error)
	GetThorchainVersion() (string, error)
	HasNetworkFee(chain common.Chain) (bool, error)
	GetNetworkFee(chain common.Chain) (transactionSize, transactionSwapSize, transactionFeeRate uint64, err error)
	PostNetworkFee(ctx context.Context, height int64, chainId *big.Int, transactionSize, transactionRate, transactionSizeWithCall uint64) (string, error)
	RagnarokInProgress() (bool, error)
	GetAsgardPubKeys() ([]PubKeyContractAddressPair, error)
	IsSyncing() (bool, error)
	WaitSync() error
	GetBlockHeight() (int64, error)
	Broadcast(hexTx []byte) (string, error)
	InitBlockScanner(...BridgeOption) error
	GetConfig() config.BifrostClientConfiguration
	GetContext() client.Context
	GetTxByBlockNumber(blockHeight int64, pk string) (types.TxOut, error)
	GetErrataMsg(txID common.TxID, chain common.Chain) sdk.Msg
	SendKeyGenStdTx(epoch *big.Int, poolPubKey common.PubKey, signature, keyShares []byte, blame []ecommon.Address,
		members []ecommon.Address) (string, error)
	GetKeyShare() ([]byte, error)
	GetKeysignParty(vaultPubKey common.PubKey) (common.PubKeys, error)
	GetInboundOutbound(txIns common.ObservedTxs) (common.ObservedTxs, common.ObservedTxs, error)
	GetSolvencyMsg(height int64, chain common.Chain, pubKey common.PubKey, coins common.Coins) *stypes.MsgSolvency
	PostKeysignFailure(blame stypes.Blame, height int64, memo string, coins common.Coins, pubkey common.PubKey) (string, error)
	GetEpochInfo(epoch *big.Int) (*structure.EpochInfo, error)
	GetChainID(name string) (*big.Int, error)
	GetTokenAddress(chainID *big.Int, name string) ([]byte, error)
	GetObservationsStdTx(txIn *types.TxIn) ([]byte, error)
}

type BridgeOption func(Bridge) error

type PubKeyContractAddressPair struct {
	PubKey    common.PubKey
	Contracts map[common.Chain]common.Address // todo will next 100
}
