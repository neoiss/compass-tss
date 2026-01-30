package mapo

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/ctx"
	keys2 "github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/tss"
	gotss "github.com/mapprotocol/compass-tss/tss/go-tss/tss"
	stypes "github.com/mapprotocol/compass-tss/x/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Endpoint urls
const (
	MimirEndpoint        = "/mapBridge/mimir"
	ChainVersionEndpoint = "/mapBridge/version"
	PoolsEndpoint        = "/mapBridge/pools"
)

// Bridge will be used to send tx to THORChain
type Bridge struct {
	logger                                                    zerolog.Logger
	cfg                                                       config.BifrostClientConfiguration
	keys                                                      *keys2.Keys
	errCounter                                                *prometheus.CounterVec
	m                                                         *metrics.Metrics
	blockHeight                                               int64
	chainID, gasPrice, epoch                                  *big.Int
	httpClient                                                *retryablehttp.Client
	broadcastLock                                             *sync.RWMutex
	ethClient                                                 *ethclient.Client
	blockScanner                                              *MapChainBlockScan
	stopChan                                                  chan struct{}
	wg                                                        *sync.WaitGroup
	gasCache                                                  []*big.Int
	ethPriKey                                                 *ecdsa.PrivateKey
	kw                                                        *evm.KeySignWrapper
	ethRpc                                                    *evm.EthRPC
	mainAbi, tssAbi, relayAbi, gasAbi, tokenRegistry, viewAbi *abi.ABI
	affiliateFeeAbi, fusionReceiverAbi                        *abi.ABI
	epochHash                                                 ecommon.Hash
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

	ethPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	signerAddr, err := k.GetEthAddress()
	if err != nil {
		return nil, err
	}
	logger.Info().Any("addr", signerAddr).Msg("Map signer address retrieved")

	rpcClient, err := evm.NewEthRPC(
		ethClient,
		time.Second*5,
		cfg.ChainID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create rpc : %w", err)
	}

	ret := &Bridge{
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
		ethRpc:        rpcClient,
		epoch:         big.NewInt(0),
		gasPrice:      big.NewInt(0),
		epochHash:     ecommon.Hash{},
	}
	err = InitAbi(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (b *Bridge) SetTssKeyManager(server *gotss.TssServer) error {
	priv, err := b.keys.GetPrivateKey()
	if err != nil {
		return fmt.Errorf("fail to get private key: %w", err)
	}
	temp, err := codec.ToCmtPubKeyInterface(priv.PubKey())
	if err != nil {
		return fmt.Errorf("fail to get tm pub key: %w", err)
	}
	pk, err := common.NewPubKeyFromCrypto(temp)
	if err != nil {
		return fmt.Errorf("fail to get pub key: %w", err)
	}

	tssKm, err := tss.NewKeySign(server, b)
	if err != nil {
		return fmt.Errorf("fail to create tss signer: %w", err)
	}
	tssKm.Start()

	keySignWrapper, err := evm.NewKeySignWrapper(b.ethPriKey, pk, tssKm, b.chainID, string(common.MAPChain))
	if err != nil {
		return fmt.Errorf("fail to create ETH key sign wrapper: %w", err)
	}
	b.kw = keySignWrapper
	return nil
}

// GetContext return a valid context with all relevant values set
func (b *Bridge) GetContext() ctx.Context {
	signerAddr, err := b.keys.GetEthAddress()
	if err != nil {
		panic(err)
	}
	ctx := ctx.Context{}
	ctx = ctx.WithKeyring(b.keys.GetKeybase())
	ctx = ctx.WithChainID(string(b.cfg.ChainID))
	ctx = ctx.WithHomeDir(b.cfg.ChainHomeFolder)
	ctx = ctx.WithFromName(b.cfg.SignerName)
	ctx = ctx.WithFromAddress(signerAddr.Hex())
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

func (b *Bridge) GetBlockScannerHeight() int64 {
	return b.blockHeight
}

func (b *Bridge) getWithPath(path string) ([]byte, int, error) {
	return b.get(b.getRelayChainURL(path))
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

// getRelayChainURL with the given path
func (b *Bridge) getRelayChainURL(path string) string {
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

// GetConfig return the configuration
func (b *Bridge) GetConfig() config.BifrostClientConfiguration {
	return b.cfg
}

// PostKeysignFailure generate and  post a keysign fail tx to thorchan
func (b *Bridge) PostKeysignFailure(blame stypes.Blame, height int64, memo string, coins common.Coins, pubkey common.PubKey) (string, error) {
	return b.Broadcast([]byte{})
}

// EnsureNodeWhitelistedWithTimeout check node is whitelisted with timeout retry
func (b *Bridge) EnsureNodeWhitelistedWithTimeout() error {
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
	if status == constants.NodeStatus_Unknown {
		return fmt.Errorf("node account status %d , will not be able to forward transaction to mapBridge", status)
	}
	return nil
}

// EnsureNodeWhitelisted will call to mapBridge to check whether the observer had been whitelist or not
func (b *Bridge) HeartBeat() error {
	go func() {
		b.logger.Info().Msg("Start heartbeat")
		for {
			select {
			case <-b.stopChan:
				b.logger.Info().Msg("Stop heartbeat")
			case <-time.After(time.Minute):
				input, err := b.mainAbi.Pack(constants.Heartbeat)
				if err != nil {
					b.logger.Error().Err(err).Msg("Fail to pack heartbeat input")
					continue
				}
				txBytes, err := b.assemblyTx(context.TODO(), input, 0, b.cfg.Maintainer)
				if err != nil {
					b.logger.Error().Err(err).Msg("Fail to assembly heartbeat tx")
					continue
				}
				txId, err := b.Broadcast(txBytes)
				if err != nil {
					b.logger.Error().Err(err).Msg("Fail to broadcast heartbeat tx")
					continue
				}
				b.logger.Debug().Msgf("Broadcast heartbeat tx %s", txId)
			}
		}
	}()
	return nil
}

// GetKeysignParty call into mapBridge to get the node accounts that should be join together to sign the message
func (b *Bridge) GetKeysignParty(vaultPubKey common.PubKey) (common.PubKeys, error) {
	return common.PubKeys{}, nil
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
		b.logger.Info().Msg("Map relay chain is syncing... waiting...")
		time.Sleep(constants.MAPRelayChainBlockTime)
	}
	return nil
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
