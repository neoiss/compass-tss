package mapo

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	ecommon "github.com/ethereum/go-ethereum/common"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	keys2 "github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	openapi "github.com/mapprotocol/compass-tss/openapi/gen"
	selfAbi "github.com/mapprotocol/compass-tss/pkg/abi"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/contract"
	stypes "github.com/mapprotocol/compass-tss/x/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Endpoint urls
const (
	AuthAccountEndpoint      = "/cosmos/auth/v1beta1/accounts"
	LastBlockEndpoint        = "/mapBridge/lastblock"
	NodeAccountEndpoint      = "/mapBridge/node"
	SignerMembershipEndpoint = "/mapBridge/vaults/%s/signers"
	StatusEndpoint           = "/status"
	VaultEndpoint            = "/mapBridge/vault/%s"
	AsgardVault              = "/mapBridge/vaults/asgard"
	PubKeysEndpoint          = "/mapBridge/vaults/pubkeys"
	ThorchainConstants       = "/mapBridge/constants"
	RagnarokEndpoint         = "/mapBridge/ragnarok"
	MimirEndpoint            = "/mapBridge/mimir"
	ChainVersionEndpoint     = "/mapBridge/version"
	InboundAddressesEndpoint = "/mapBridge/inbound_addresses"
	PoolsEndpoint            = "/mapBridge/pools"
	THORNameEndpoint         = "/mapBridge/thorname/%s"
)

// Bridge will be used to send tx to THORChain
type Bridge struct {
	logger        zerolog.Logger
	cfg           config.BifrostClientConfiguration
	keys          *keys2.Keys
	errCounter    *prometheus.CounterVec
	m             *metrics.Metrics
	blockHeight   int64
	accountNumber uint64
	seqNumber     uint64
	chainID       *big.Int
	httpClient    *retryablehttp.Client
	broadcastLock *sync.RWMutex
	ethClient     *ethclient.Client
	blockScanner  *MapChainBlockScan
	stopChan      chan struct{}
	wg            *sync.WaitGroup
	gasPrice      *big.Int
	gasCache      []*big.Int
	ethPriKey     *ecdsa.PrivateKey
	kw            *evm.KeySignWrapper
	ethRpc        *evm.EthRPC
	mainAbi       *abi.ABI
	mainCall      *contract.Call
	epoch         *big.Int
}

// httpResponseCache used for caching HTTP responses for less frequent querying
type httpResponseCache struct {
	httpResponse        []byte
	httpResponseChecked time.Time
	httpResponseMu      *sync.Mutex
}

var (
	httpResponseCaches   = make(map[string]*httpResponseCache) // String-to-pointer map for quicker lookup
	httpResponseCachesMu = &sync.Mutex{}
)

// NewBridge create a new instance of Bridge
func NewBridge(cfg config.BifrostClientConfiguration, m *metrics.Metrics, k *keys2.Keys) (shareTypes.Bridge, error) {
	// main module logger
	logger := log.With().Str("module", "mapo_client").Logger()

	if len(cfg.ChainID) == 0 {
		return nil, errors.New("chain id is empty")
	}
	if len(cfg.ChainHost) == 0 {
		return nil, errors.New("chain host is empty")
	}

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	ethClient, err := ethclient.Dial(cfg.ChainHost)
	if err != nil {
		return nil, fmt.Errorf("fail to dial map rpc host(%s): %w", cfg.ChainHost, err)
	}

	chainID, err := getChainID(ethClient, time.Second*5)
	if err != nil {
		return nil, err
	}

	priv, err := k.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get private key: %w", err)
	}
	temp, err := codec.ToCmtPubKeyInterface(priv.PubKey())
	if err != nil {
		return nil, fmt.Errorf("fail to get tm pub key: %w", err)
	}
	pk, err := common.NewPubKeyFromCrypto(temp)
	if err != nil {
		return nil, fmt.Errorf("fail to get pub key: %w", err)
	}
	ethPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	mainAbi, err := newMaintainerABi()
	if err != nil {
		return nil, err
	}

	ai, err := selfAbi.New(maintainerAbi)
	if err != nil {
		return nil, err
	}

	mainCall := contract.New(ethClient, []ecommon.Address{ecommon.HexToAddress(cfg.Maintainer)}, ai)
	keySignWrapper, err := evm.NewKeySignWrapper(ethPrivateKey, pk, nil, chainID, "MAP")
	if err != nil {
		return nil, fmt.Errorf("fail to create ETH key sign wrapper: %w", err)
	}

	rpcClient, err := evm.NewEthRPC(
		ethClient,
		time.Second*5,
		cfg.ChainID.String(),
	)

	return &Bridge{
		logger:        logger,
		cfg:           cfg,
		keys:          k,
		errCounter:    m.GetCounterVec(metrics.MapChainClientError),
		httpClient:    httpClient,
		m:             m,
		chainID:       chainID,
		broadcastLock: &sync.RWMutex{},
		ethClient:     ethClient,
		stopChan:      make(chan struct{}),
		wg:            &sync.WaitGroup{},
		ethPriKey:     ethPrivateKey,
		kw:            keySignWrapper,
		ethRpc:        rpcClient,
		mainAbi:       mainAbi,
		mainCall:      mainCall,
		epoch:         big.NewInt(0),
		gasPrice:      big.NewInt(0),
	}, nil
}

// GetContext return a valid context with all relevant values set
func (b *Bridge) GetContext() client.Context {
	signerAddr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}
	ctx := client.Context{}
	ctx = ctx.WithKeyring(b.keys.GetKeybase())
	ctx = ctx.WithChainID(string(b.cfg.ChainID))
	ctx = ctx.WithHomeDir(b.cfg.ChainHomeFolder)
	ctx = ctx.WithFromName(b.cfg.SignerName)
	ctx = ctx.WithFromAddress(signerAddr)
	ctx = ctx.WithBroadcastMode("sync")

	remote := b.cfg.ChainRPC
	if !strings.HasPrefix(b.cfg.ChainHost, "http") {
		remote = fmt.Sprintf("tcp://%s", remote)
	}
	ctx = ctx.WithNodeURI(remote)
	client, err := rpchttp.New(remote, "/websocket")
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithClient(client)
	return ctx
}

func (b *Bridge) getWithPath(path string) ([]byte, int, error) {
	return b.get(b.getThorChainURL(path))
}

// get handle all the low level http GET calls using retryablehttp.Bridge
func (b *Bridge) get(url string) ([]byte, int, error) {
	// To reduce querying time and chance of "429 Too Many Requests",
	// do not query the same endpoint more than once per block time.
	httpResponseCachesMu.Lock()
	respCachePointer := httpResponseCaches[url]
	if respCachePointer == nil {
		// Since this is the first time using this endpoint, prepare a Mutex for it.
		respCachePointer = &httpResponseCache{httpResponseMu: &sync.Mutex{}}
		httpResponseCaches[url] = respCachePointer
	}
	httpResponseCachesMu.Unlock()

	// So lengthy queries don't hold up short queries, use query-specific mutexes.
	respCachePointer.httpResponseMu.Lock()
	defer respCachePointer.httpResponseMu.Unlock()

	// When the same endpoint has been checked within the span of a single block, return the cached response.
	if time.Since(respCachePointer.httpResponseChecked) < constants.MAPRelayChainBlockTime && respCachePointer.httpResponse != nil {
		return respCachePointer.httpResponse, http.StatusOK, nil
	}

	resp, err := b.httpClient.Get(url)
	if err != nil {
		b.errCounter.WithLabelValues("fail_get_from_thorchain", "").Inc()
		return nil, http.StatusNotFound, fmt.Errorf("failed to GET from mapBridge: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	buf, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return buf, resp.StatusCode, errors.New("Status code: " + resp.Status + " returned")
	}
	if err != nil {
		b.errCounter.WithLabelValues("fail_read_thorchain_resp", "").Inc()
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	// All being well with the response, save it to the cache.
	respCachePointer.httpResponse = buf
	respCachePointer.httpResponseChecked = time.Now()

	return buf, resp.StatusCode, nil
}

// getThorChainURL with the given path
func (b *Bridge) getThorChainURL(path string) string {
	if strings.HasPrefix(b.cfg.ChainHost, "http") {
		return fmt.Sprintf("%s/%s", b.cfg.ChainHost, path)
	}

	uri := url.URL{
		Scheme: "http",
		Host:   b.cfg.ChainHost,
		Path:   path,
	}
	return uri.String()
}

// getAccountNumberAndSequenceNumber returns account and Sequence number required to post into mapBridge
func (b *Bridge) getAccountNumberAndSequenceNumber() (uint64, uint64, error) {
	signerAddr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get signer address: %w", err)
	}
	path := fmt.Sprintf("%s/%s", AuthAccountEndpoint, signerAddr)

	body, _, err := b.getWithPath(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get auth accounts: %w", err)
	}

	var resp types.AccountResp
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal account resp: %w", err)
	}
	acc := resp.Account

	return acc.AccountNumber, acc.Sequence, nil
}

// GetConfig return the configuration
func (b *Bridge) GetConfig() config.BifrostClientConfiguration {
	return b.cfg
}

// PostKeysignFailure generate and  post a keysign fail tx to thorchan
func (b *Bridge) PostKeysignFailure(blame stypes.Blame, height int64, memo string, coins common.Coins, pubkey common.PubKey) (string, error) {
	return b.Broadcast([]byte{})
}

// GetErrataMsg get errata tx from params
func (b *Bridge) GetErrataMsg(txID common.TxID, chain common.Chain) sdk.Msg {
	signerAddr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}
	return stypes.NewMsgErrataTx(txID, chain, signerAddr)
}

// GetSolvencyMsg create MsgSolvency from the given parameters
func (b *Bridge) GetSolvencyMsg(height int64, chain common.Chain, pubKey common.PubKey, coins common.Coins) *stypes.MsgSolvency {
	// To prevent different MsgSolvency ID incompatibility between nodes with different coin-observation histories,
	// only report coins for which the amounts are not currently 0.
	coins = coins.NoneEmpty()
	signerAddr, err := b.keys.GetSignerInfo().GetAddress()
	if err != nil {
		panic(err)
	}
	msg, err := stypes.NewMsgSolvency(chain, pubKey, coins, height, signerAddr)
	if err != nil {
		b.logger.Err(err).Msg("fail to create MsgSolvency")
		return nil
	}
	return msg
}

// GetInboundOutbound separate the txs into inbound and outbound
func (b *Bridge) GetInboundOutbound(txIns common.ObservedTxs) (common.ObservedTxs, common.ObservedTxs, error) {
	if len(txIns) == 0 {
		return nil, nil, nil
	}
	inbound := common.ObservedTxs{}
	outbound := common.ObservedTxs{}

	// spilt our txs into inbound vs outbound txs
	for _, tx := range txIns {
		chain := common.EmptyChain
		if len(tx.Tx.Coins) > 0 {
			chain = tx.Tx.Coins[0].Asset.Chain
		}

		obAddr, err := tx.ObservedPubKey.GetAddress(chain)
		if err != nil {
			b.logger.Err(err).Msgf("fail to parse observed pool address: %s", tx.ObservedPubKey.String())
			continue
		}
		vaultToAddress := tx.Tx.ToAddress.Equals(obAddr)
		vaultFromAddress := tx.Tx.FromAddress.Equals(obAddr)
		var inInboundArray, inOutboundArray bool
		if vaultToAddress {
			inInboundArray = inbound.Contains(tx)
		}
		if vaultFromAddress {
			inOutboundArray = outbound.Contains(tx)
		}
		// for consolidate UTXO tx, both From & To address will be the asgard address
		// thus here we need to make sure that one add to inbound , the other add to outbound
		switch {
		case !vaultToAddress && !vaultFromAddress:
			// Neither ToAddress nor FromAddress matches obAddr, so drop it.
			b.logger.Error().Msgf("chain (%s) tx (%s) observedaddress (%s) does not match its toaddress (%s) or fromaddress (%s)", tx.Tx.Chain, tx.Tx.ID, obAddr, tx.Tx.ToAddress, tx.Tx.FromAddress)
		case vaultToAddress && !inInboundArray:
			inbound = append(inbound, tx)
		case vaultFromAddress && !inOutboundArray:
			outbound = append(outbound, tx)
		case inInboundArray && inOutboundArray:
			// It's already in both arrays, so drop it.
			b.logger.Error().Msgf("vault-to-vault chain (%s) tx (%s) is already in both inbound and outbound arrays", tx.Tx.Chain, tx.Tx.ID)
		case !vaultFromAddress && inInboundArray:
			// It's already in its only (inbound) array, so drop it.
			b.logger.Error().Msgf("observed tx in for chain (%s) tx (%s) is already in the inbound array", tx.Tx.Chain, tx.Tx.ID)
		case !vaultToAddress && inOutboundArray:
			// It's already in its only (outbound) array, so drop it.
			b.logger.Error().Msgf("observed tx out for chain (%s) tx (%s) is already in the outbound array", tx.Tx.Chain, tx.Tx.ID)
		default:
			// This should never happen; rather than dropping it, return an error.
			return nil, nil, fmt.Errorf("could not determine if chain (%s) tx (%s) was inbound or outbound", tx.Tx.Chain, tx.Tx.ID)
		}
	}

	return inbound, outbound, nil
}

// EnsureNodeWhitelistedWithTimeout check node is whitelisted with timeout retry
func (b *Bridge) EnsureNodeWhitelistedWithTimeout() error {
	// todo handler done
	for {
		select {
		case <-time.After(time.Hour):
			return errors.New("Observer is not whitelisted yet")
		default:
			err := b.EnsureNodeWhitelisted()
			if err == nil {
				// node had been whitelisted
				return nil
			}
			b.logger.Error().Err(err).Msg("observer is not whitelisted , will retry a bit later")
			time.Sleep(time.Second * 5)
		}
	}
}

// EnsureNodeWhitelisted will call to mapBridge to check whether the observer had been whitelist or not
func (b *Bridge) EnsureNodeWhitelisted() error {
	status, err := b.FetchNodeStatus()
	if err != nil {
		return fmt.Errorf("failed to get node status: %w", err)
	}
	if status == stypes.NodeStatus_Unknown {
		return fmt.Errorf("node account status %s , will not be able to forward transaction to mapBridge", status)
	}
	return nil
}

// GetKeysignParty call into mapBridge to get the node accounts that should be join together to sign the message
func (b *Bridge) GetKeysignParty(vaultPubKey common.PubKey) (common.PubKeys, error) {
	p := fmt.Sprintf(SignerMembershipEndpoint, vaultPubKey.String())
	result, _, err := b.getWithPath(p)
	if err != nil {
		return common.PubKeys{}, fmt.Errorf("fail to get key sign party from mapBridge: %w", err)
	}
	var keys common.PubKeys
	if err = json.Unmarshal(result, &keys); err != nil {
		return common.PubKeys{}, fmt.Errorf("fail to unmarshal result to pubkeys:%w", err)
	}
	return keys, nil
}

// IsSyncing returns bool for if map relay is catching up to the rest of the
// nodes. Returns yes, if it is, false if it is caught up.
func (b *Bridge) IsSyncing() (bool, error) {
	progress, err := b.ethClient.SyncProgress(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to get sync progress: %w", err)
	}
	//return progress == nil, nil
	if progress == nil {
		return false, nil
	}
	return progress.CurrentBlock < progress.HighestBlock, nil
}

// HasNetworkFee checks whether the given chain has set a network fee - determined by
// whether the `outbound_tx_size` for the inbound address response is non-zero.
func (b *Bridge) HasNetworkFee(chain common.Chain) (bool, error) {
	buf, s, err := b.getWithPath(InboundAddressesEndpoint)
	if err != nil {
		return false, fmt.Errorf("fail to get inbound addresses: %w", err)
	}
	if s != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", s)
	}

	var resp []openapi.InboundAddress
	if err = json.Unmarshal(buf, &resp); err != nil {
		return false, fmt.Errorf("fail to unmarshal inbound addresses: %w", err)
	}

	for _, addr := range resp {
		if addr.Chain != nil && *addr.Chain == chain.String() && addr.OutboundTxSize != nil {
			var size int64
			size, err = strconv.ParseInt(*addr.OutboundTxSize, 10, 64)
			if err != nil {
				return false, fmt.Errorf("fail to parse outbound_tx_size: %w", err)
			}
			return size > 0, nil
		}
	}

	return false, fmt.Errorf("no inbound address found for chain: %s", chain)
}

// GetNetworkFee get chain's network fee from THORNode.
func (b *Bridge) GetNetworkFee(chain common.Chain) (transactionSize, transactionFeeRate uint64, err error) {
	buf, s, err := b.getWithPath(InboundAddressesEndpoint)
	if err != nil {
		return 0, 0, fmt.Errorf("fail to get inbound addresses: %w", err)
	}
	if s != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status code: %d", s)
	}
	var resp []openapi.InboundAddress
	if err = json.Unmarshal(buf, &resp); err != nil {
		return 0, 0, fmt.Errorf("fail to unmarshal to json: %w", err)
	}

	for _, addr := range resp {
		if addr.Chain != nil && *addr.Chain == chain.String() {
			// Default values if nil or unfound are 0.
			if addr.OutboundTxSize != nil {
				transactionSize, err = strconv.ParseUint(*addr.OutboundTxSize, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("fail to parse outbound_tx_size: %w", err)
				}
			}
			if addr.ObservedFeeRate != nil {
				transactionFeeRate, err = strconv.ParseUint(*addr.ObservedFeeRate, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("fail to parse observed_fee_rate: %w", err)
				}
			}
			// Having found the chain, do not continue through the remaining chains.
			break
		}
	}
	return
}

// WaitSync wait for map relay chain to catch up
func (b *Bridge) WaitSync() error {
	for {
		yes, err := b.IsSyncing()
		if err != nil {
			return err
		}
		if !yes {
			break
		}
		b.logger.Info().Msg("map relay chain is syncing... waiting...")
		time.Sleep(constants.MAPRelayChainBlockTime)
	}
	return nil
}

// GetAsgards retrieve all the asgard vaults from mapBridge
func (b *Bridge) GetAsgards() (stypes.Vaults, error) {
	buf, s, err := b.getWithPath(AsgardVault)
	if err != nil {
		return nil, fmt.Errorf("fail to get asgard vaults: %w", err)
	}
	if s != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", s)
	}
	var vaults stypes.Vaults
	if err = json.Unmarshal(buf, &vaults); err != nil {
		return nil, fmt.Errorf("fail to unmarshal asgard vaults from json: %w", err)
	}
	return vaults, nil
}

// GetVault retrieves a specific vault from mapBridge.
func (b *Bridge) GetVault(pubkey string) (stypes.Vault, error) {
	buf, s, err := b.getWithPath(fmt.Sprintf(VaultEndpoint, pubkey))
	if err != nil {
		return stypes.Vault{}, fmt.Errorf("fail to get vault: %w", err)
	}
	if s != http.StatusOK {
		return stypes.Vault{}, fmt.Errorf("unexpected status code %d", s)
	}
	var vault stypes.Vault
	if err = json.Unmarshal(buf, &vault); err != nil {
		return stypes.Vault{}, fmt.Errorf("fail to unmarshal vault from json: %w", err)
	}
	return vault, nil
}

func (b *Bridge) getVaultPubkeys() ([]byte, error) {
	return nil, nil
	//buf, s, err := b.getWithPath(PubKeysEndpoint)
	//if err != nil {
	//	return nil, fmt.Errorf("fail to get asgard vaults: %w", err)
	//}
	//if s != http.StatusOK {
	//	return nil, fmt.Errorf("unexpected status code %d", s)
	//}
	//return buf, nil
}

// GetPubKeys retrieve vault pub keys and their relevant smart contracts
func (b *Bridge) GetPubKeys() ([]shareTypes.PubKeyContractAddressPair, error) {
	//// todo handler
	//buf, err := b.getVaultPubkeys()
	//if err != nil {
	//	return nil, fmt.Errorf("fail to get vault pubkeys ,err: %w", err)
	//}
	//var result openapi.VaultPubkeysResponse
	//if err = json.Unmarshal(buf, &result); err != nil {
	//	return nil, fmt.Errorf("fail to unmarshal pubkeys: %w", err)
	//}
	//var addressPairs []shareTypes.PubKeyContractAddressPair
	//for _, v := range append(result.Asgard, result.Inactive...) {
	//	kp := shareTypes.PubKeyContractAddressPair{
	//		PubKey:    common.PubKey(v.PubKey),
	//		Contracts: make(map[common.Chain]common.Address),
	//	}
	//	for _, item := range v.Routers {
	//		kp.Contracts[common.Chain(*item.Chain)] = common.Address(*item.Router)
	//	}
	//	addressPairs = append(addressPairs, kp)
	//}
	//return addressPairs, nil
	return nil, nil
}

// GetAsgardPubKeys retrieve asgard vaults, and it's relevant smart contracts
func (b *Bridge) GetAsgardPubKeys() ([]shareTypes.PubKeyContractAddressPair, error) {
	// todo temp code, will next 200
	return []shareTypes.PubKeyContractAddressPair{
		{
			PubKey: "029038a5cabb18c0bd3017b631d08feedf8107c816f3cd1783c26037516bfd7754",
		},
	}, nil

	//buf, err := b.getVaultPubkeys()
	//if err != nil {
	//	return nil, fmt.Errorf("fail to get vault pubkeys ,err: %w", err)
	//}
	//var result openapi.VaultPubkeysResponse
	//if err = json.Unmarshal(buf, &result); err != nil {
	//	return nil, fmt.Errorf("fail to unmarshal pubkeys: %w", err)
	//}
	//var addressPairs []shareTypes.PubKeyContractAddressPair
	//for _, v := range append(result.Asgard, result.Inactive...) {
	//	kp := shareTypes.PubKeyContractAddressPair{
	//		PubKey:    common.PubKey(v.PubKey),
	//		Contracts: make(map[common.Chain]common.Address),
	//	}
	//	for _, item := range v.Routers {
	//		kp.Contracts[common.Chain(*item.Chain)] = common.Address(*item.Router)
	//	}
	//	addressPairs = append(addressPairs, kp)
	//}
	//return addressPairs, nil
}

// PostNetworkFee send network fee message to THORNode
func (b *Bridge) PostNetworkFee(height int64, chain common.Chain, transactionSize, transactionRate, transactionSizeWithCall uint64) (string, error) {
	// todo next 1
	cId, err := chain.ChainID()
	if err != nil {
		return "", fmt.Errorf("fail to get chain id: %w", err)
	}
	b.mainAbi.Pack(constants.VoteNetworkFeeOfMaintainer, big.NewInt(1), cId, big.NewInt(height),
		big.NewInt(0).SetUint64(transactionSize),
		big.NewInt(0).SetUint64(transactionRate),
		big.NewInt(0).SetUint64(transactionSizeWithCall))

	return b.Broadcast([]byte{})
}

// GetConstants from thornode
func (b *Bridge) GetConstants() (map[string]int64, error) {
	var result struct {
		Int64Values map[string]int64 `json:"int_64_values"`
	}
	buf, s, err := b.getWithPath(ThorchainConstants)
	if err != nil {
		return nil, fmt.Errorf("fail to get constants: %w", err)
	}
	if s != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", s)
	}
	if err = json.Unmarshal(buf, &result); err != nil {
		return nil, fmt.Errorf("fail to unmarshal to json: %w", err)
	}
	return result.Int64Values, nil
}

// RagnarokInProgress is to query mapBridge to check whether ragnarok had been triggered
func (b *Bridge) RagnarokInProgress() (bool, error) {
	buf, s, err := b.getWithPath(RagnarokEndpoint)
	if err != nil {
		return false, fmt.Errorf("fail to get ragnarok status: %w", err)
	}
	if s != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", s)
	}
	var ragnarok bool
	if err = json.Unmarshal(buf, &ragnarok); err != nil {
		return false, fmt.Errorf("fail to unmarshal ragnarok status: %w", err)
	}
	return ragnarok, nil
}

// GetThorchainVersion retrieve mapBridge version
// func (b *Bridge) GetThorchainVersion() (semver.Version, error) {
func (b *Bridge) GetThorchainVersion() (string, error) {
	// todo handler
	return "1", nil
}

// GetMimir - get mimir settings
func (b *Bridge) GetMimir(key string) (int64, error) {
	//buf, s, err := b.getWithPath(MimirEndpoint + "/key/" + key)
	//if err != nil {
	//	return 0, fmt.Errorf("fail to get mimir: %w", err)
	//}
	//if s != http.StatusOK {
	//	return 0, fmt.Errorf("unexpected status code: %d", s)
	//}
	//var value int64
	//if err = json.Unmarshal(buf, &value); err != nil {
	//	return 0, fmt.Errorf("fail to unmarshal mimir: %w", err)
	//}
	//return value, nil
	// todo handler
	return 0, nil
}

// GetMimirWithRef is a helper function to more readably insert references (such as Asset MimirString or Chain) into Mimir key templates.
func (b *Bridge) GetMimirWithRef(template, ref string) (int64, error) {
	// 'template' should be something like "Halt%sChain" (to halt an arbitrary specified chain)
	// or "Ragnarok-%s" (to halt the pool of an arbitrary specified Asset (MimirString used for Assets to join Chain and Symbol with a hyphen).
	key := fmt.Sprintf(template, ref)
	return b.GetMimir(key)
}

// PubKeyContractAddressPair is an entry to map pubkey and contract addresses

// GetContractAddress retrieve the contract address from asgard
func (b *Bridge) GetContractAddress() ([]shareTypes.PubKeyContractAddressPair, error) {
	buf, s, err := b.getWithPath(InboundAddressesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("fail to get inbound addresses: %w", err)
	}
	if s != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", s)
	}
	type address struct {
		Chain   common.Chain   `json:"chain"`
		PubKey  common.PubKey  `json:"pub_key"`
		Address common.Address `json:"address"`
		Router  common.Address `json:"router"`
		Halted  bool           `json:"halted"`
	}
	var resp []address
	if err = json.Unmarshal(buf, &resp); err != nil {
		return nil, fmt.Errorf("fail to unmarshal response: %w", err)
	}
	var result []shareTypes.PubKeyContractAddressPair
	for _, item := range resp {
		exist := false
		for _, pair := range result {
			if item.PubKey.Equals(pair.PubKey) {
				pair.Contracts[item.Chain] = item.Router
				exist = true
				break
			}
		}
		if !exist {
			pair := shareTypes.PubKeyContractAddressPair{
				PubKey:    item.PubKey,
				Contracts: map[common.Chain]common.Address{},
			}
			pair.Contracts[item.Chain] = item.Router
			result = append(result, pair)
		}
	}
	return result, nil
}

// GetPools get pools from THORChain
func (b *Bridge) GetPools() (stypes.Pools, error) {
	buf, s, err := b.getWithPath(PoolsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("fail to get pools addresses: %w", err)
	}
	if s != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", s)
	}
	var pools stypes.Pools
	if err = json.Unmarshal(buf, &pools); err != nil {
		return nil, fmt.Errorf("fail to unmarshal pools from json: %w", err)
	}
	return pools, nil
}

// GetTHORName get THORName from THORChain
func (b *Bridge) GetTHORName(name string) (stypes.THORName, error) {
	p := fmt.Sprintf(THORNameEndpoint, name)
	buf, s, err := b.getWithPath(p)
	if err != nil {
		return stypes.THORName{}, fmt.Errorf("fail to get THORName: %w", err)
	}
	if s != http.StatusOK {
		return stypes.THORName{}, fmt.Errorf("unexpected status code: %d", s)
	}
	var tn stypes.THORName
	if err = json.Unmarshal(buf, &tn); err != nil {
		return stypes.THORName{}, fmt.Errorf("fail to unmarshal THORNames from json: %w", err)
	}
	return tn, nil
}

// GetChain returns the chain.
func (b *Bridge) GetChain() common.Chain {
	return b.cfg.ChainID
}

func WithBlockScanner(blockScan *MapChainBlockScan) shareTypes.BridgeOption {
	return func(bridge shareTypes.Bridge) error {
		mapBridge, ok := bridge.(*Bridge)
		if !ok {
			return fmt.Errorf("invalid bridge type(%v)", reflect.TypeOf(bridge))
		}
		mapBridge.blockScanner = blockScan
		return nil
	}
}

func (b *Bridge) InitBlockScanner(ops ...shareTypes.BridgeOption) error {
	for _, op := range ops {
		err := op(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bridge) GetObservationsStdTx(txIn *types.TxIn) ([]byte, error) {
	// todo check
	// Here we construct tx according to method， and return hex tx2bytes
	return nil, nil
}
