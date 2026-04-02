package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	tcconfig "github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	tcmetrics "github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/runners"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/signercache"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/tron/api"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/tron/rpc"
	"github.com/mapprotocol/compass-tss/tss"
	tctss "github.com/mapprotocol/compass-tss/tss/go-tss/tss"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// TimestampValidity is the expiration duration for an outbound transaction. We use
	// the thornode block time of the outbound height as the timestamp to preserve
	// determinism, so the 20 minute expiration means the built transaction will be valid
	// to broadcast for the first 20 minutes of the 30 minute signing period.
	TimestampValidity = 20 * time.Minute

	ConfirmationBlocks int64 = 1
	FinalityBlocks     int64 = 19
)

//go:embed abi/gateway.json
var gatewayABI []byte

type TronClient struct {
	logger              zerolog.Logger
	cfg                 tcconfig.BifrostChainConfiguration
	chainId             string
	blockScanner        *blockscanner.BlockScanner
	storage             *blockscanner.BlockScannerStorage
	signerCacheManager  *signercache.CacheManager
	tssKeyManager       *tss.KeySign
	localKeyManager     *KeyManager
	tronScanner         *TronBlockScanner
	api                 *api.TronApi
	rpc                 *rpc.TronRpc
	bridge              shareTypes.Bridge
	globalSolvencyQueue chan stypes.Solvency
	wg                  *sync.WaitGroup
	stopchan            chan struct{}
	gatewayAbi          *abi.ABI
}

// NewTronClient creates a new instance of a Tron chain client //
func NewTronClient(
	relayKey *keys.Keys,
	config tcconfig.BifrostChainConfiguration,
	server *tctss.TssServer,
	bridge shareTypes.Bridge,
	metrics *tcmetrics.Metrics,
) (*TronClient, error) {
	var err error
	logger := log.With().Str("module", config.ChainID.String()).Logger()

	client := TronClient{
		logger:   logger,
		chainId:  config.ChainID.String(),
		cfg:      config,
		bridge:   bridge,
		wg:       &sync.WaitGroup{},
		stopchan: make(chan struct{}),
		api:      api.NewTronApi(config.RPCHost, config.BlockScanner.HTTPRequestTimeout),
		rpc:      rpc.NewTronRpc(config.RPCHost, config.BlockScanner.HTTPRequestTimeout),
	}

	client.tssKeyManager, err = tss.NewKeySign(server, bridge)
	if err != nil {
		logger.Err(err).Msg("failed to create tss signer")
		return nil, err
	}

	client.localKeyManager, err = NewLocalKeyManager(relayKey)
	if err != nil {
		logger.Err(err).Msg("failed to create local key manager")
		return nil, err
	}

	var path string // if not set later, will in memory storage
	if len(client.cfg.BlockScanner.DBPath) > 0 {
		path = fmt.Sprintf(
			"%s/%s", config.BlockScanner.DBPath, config.BlockScanner.ChainID,
		)
	}

	client.storage, err = blockscanner.NewBlockScannerStorage(
		path,
		client.cfg.ScannerLevelDB,
	)
	if err != nil {
		logger.Err(err).Msg("failed to create scan storage")
		return nil, err
	}

	client.tronScanner, err = NewTronBlockScanner(
		config,
		client.bridge,
	)
	if err != nil {
		logger.Err(err).Msg("failed to create tron block scanner")
		return nil, err
	}

	client.blockScanner, err = blockscanner.NewBlockScanner(
		client.cfg.BlockScanner,
		client.storage,
		metrics,
		client.bridge,
		client.tronScanner,
	)
	if err != nil {
		logger.Err(err).Msg("failed to create block scanner")
		return nil, err
	}

	client.signerCacheManager, err = signercache.NewSignerCacheManager(
		client.storage.GetInternalDb(),
	)
	if err != nil {
		logger.Err(err).Msg("failed to create signer cache manager")
		return nil, err
	}
	relayKey.GetEthAddress()

	return &client, nil
}

// Start Tron chain client
func (c *TronClient) Start(
	globalTxsQueue chan types.TxIn,
	_ chan types.ErrataBlock,
	globalSolvencyQueue chan stypes.Solvency,
	globalNetworkFeeQueue chan stypes.NetworkFee,
) {
	c.globalSolvencyQueue = globalSolvencyQueue
	c.tronScanner.globalNetworkFeeQueue = globalNetworkFeeQueue
	c.blockScanner.Start(globalTxsQueue, globalNetworkFeeQueue)
	c.tssKeyManager.Start()

	c.wg.Add(1)
	go runners.SolvencyCheckRunner(
		c.GetChain(),
		c,
		c.bridge,
		c.stopchan,
		c.wg,
		time.Second,
	)
}

// Stop Tron chain client
func (c *TronClient) Stop() {
	c.tssKeyManager.Stop()
	c.blockScanner.Stop()
	close(c.stopchan)
	c.wg.Wait()
}

func (c *TronClient) IsBlockScannerHealthy() bool {
	return c.blockScanner.IsHealthy()
}

// GetChain returns the chain.
func (c *TronClient) GetChain() common.Chain {
	return c.cfg.ChainID
}

// GetConfig returns the chain client configuration
func (c *TronClient) GetConfig() tcconfig.BifrostChainConfiguration {
	return c.cfg
}

// GetHeight returns the current height of the chain.
func (c *TronClient) GetHeight() (int64, error) {
	return c.tronScanner.GetHeight()
}

// GetAddress returns the address for the given public key.
func (c *TronClient) GetAddress(pubKey common.PubKey) string {
	address, err := pubKey.GetAddress(c.GetChain())
	if err != nil {
		c.logger.Err(err).Msg("failed to get pool address")
		return ""
	}

	return address.String()
}

// GetAccount returns the account for the given public key.
func (c *TronClient) GetAccount(
	pubKey common.PubKey,
	height *big.Int,
) (common.Account, error) {
	address, err := pubKey.GetAddress(c.GetChain())
	if err != nil {
		c.logger.Err(err).
			Str("pubkey", pubKey.String()).
			Msg("failed to get pool address")
		return common.Account{}, err
	}
	return c.GetAccountByAddress(address.String(), height)
}

// GetAccountByAddress returns the account for the given address.
func (c *TronClient) GetAccountByAddress(
	address string,
	_ *big.Int,
) (common.Account, error) {
	account := common.Account{}
	return account, nil
}

// GetBlockScannerHeight returns block scanner height for chain
func (c *TronClient) GetBlockScannerHeight() (int64, error) {
	return c.blockScanner.PreviousHeight(), nil
}

// GetLatestTxForVault returns last observed and broadcasted tx for a particular vault and chain
func (c *TronClient) GetLatestTxForVault(vault string) (string, string, error) {
	lastObserved, err := c.signerCacheManager.GetLatestRecordedTx(
		types.InboundCacheKey(vault, c.GetChain().String()),
	)
	if err != nil {
		return "", "", err
	}
	lastBroadcasted, err := c.signerCacheManager.GetLatestRecordedTx(
		types.BroadcastCacheKey(vault, c.GetChain().String()),
	)
	return lastObserved, lastBroadcasted, err
}

// GetConfirmationCount returns the confirmation count for the given tx.
func (c *TronClient) GetConfirmationCount(_ types.TxIn) int64 {
	// https://developers.tron.network/docs/tron-protocol-transaction#transaction-lifecycle
	// We are scanning 19 blocks behind the actual tip, so returning 0 here
	return 0
}

func (c *TronClient) ConfirmationCountReady(_ types.TxIn) bool {
	return true
}

// OnObservedTxIn is called when a new observed tx is received.
func (c *TronClient) OnObservedTxIn(_ types.TxInItem, _ int64) {}

// SignTx returns the signed transaction.
func (c *TronClient) SignTx(
	txOutItem types.TxOutItem,
	_ int64,
) ([]byte, []byte, *types.TxInItem, error) {
	if c.signerCacheManager.HasSigned(txOutItem.CacheHash()) {
		c.logger.Info().Interface("txOutItem", txOutItem).Msg("transaction already signed, ignoring...")
		return nil, nil, nil, nil
	}
	// parse the chain and gas limit from the memo
	cgl, err := evm.ParseChainAndGasLimit(ethcommon.BytesToHash(common.Completion(txOutItem.ChainAndGasLimit.Bytes(), 32)))
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("fail to parse chain and gas limit")
		return nil, nil, nil, err
	}
	c.logger.Info().Str("relayHash", txOutItem.TxHash).Str("tx_rate", cgl.Third.String()).Str("tx_size", cgl.End.String())
	// ---------------------------------------------------------------------------
	fromAddr, err := c.localKeyManager.Pubkey().GetAddress(c.cfg.ChainID)
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("failed to get address from pubkey")
		return nil, nil, nil, err
	}
	inputBytes, err := c.buildInput(fromAddr, txOutItem)
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("failed to build input")
		return nil, nil, nil, err
	}

	contract, err := api.ConvertAddress(c.cfg.BlockScanner.Mos)
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("failed to convert contract address")
		return nil, nil, nil, err
	}

	apiTx, err := c.api.TriggerSmartContract(strings.Replace(fromAddr.String(), "0x", "41", 1), contract,
		"bridgeIn(address,bytes32,bytes,bytes)",
		hex.EncodeToString(inputBytes[4:]), cgl.Third.Uint64())
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("failed to trigger smart contract")
		return nil, nil, nil, err
	}
	err = c.signTx(&apiTx)
	if err != nil {
		c.logger.Err(err).Str("relayHash", txOutItem.TxHash).Msg("failed to sign transaction")
		return nil, nil, nil, err
	}

	txBytes, err := json.Marshal(&apiTx)
	if err != nil {
		c.logger.Err(err).Msg("failed to marshal tx")
		return nil, nil, nil, err
	}

	return txBytes, nil, nil, nil
}

func (c *TronClient) buildInput(sender common.Address, txOutItem stypes.TxOutItem) ([]byte, error) {
	var (
		err error
		ret []byte
	)
	switch txOutItem.Method {
	case constants.BridgeIn:
		ret, err = c.gatewayAbi.Pack(constants.BridgeIn,
			ethcommon.HexToAddress(sender.String()),
			txOutItem.OrderId, txOutItem.Data, txOutItem.Signature)
	default:
		return nil, fmt.Errorf("not support method(%s)", txOutItem.Method)
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *TronClient) base58ToHex(addr string) (string, error) {
	raw, err := base58.Decode(addr)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(raw[:21]), nil
}

func (c *TronClient) signTx(tx *api.Transaction) error {
	rawBytes, err := hex.DecodeString(tx.RawDataHex)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(rawBytes)
	signature, err := c.localKeyManager.Sign(hash[:])
	if err != nil {
		return err
	}

	tx.Signature = []string{hex.EncodeToString(signature)}
	return nil
}

// BroadcastTx sends the transaction to Tron chain
func (c *TronClient) BroadcastTx(
	txOutItem types.TxOutItem,
	txBytes []byte,
) (string, error) {
	response, err := c.api.BroadcastTransaction(txBytes)
	if err != nil {
		c.logger.Err(err).Msg("failed to broadcast tx")
		return "", err
	}

	// treat dup transaction as success - the transaction was already accepted
	if !response.Result && response.Code == "DUP_TRANSACTION_ERROR" {
		var tx api.Transaction
		if unmarshalErr := json.Unmarshal(txBytes, &tx); unmarshalErr != nil {
			return "", fmt.Errorf("failed to unmarshal tx for dup txid: %w", unmarshalErr)
		}
		c.logger.Info().Str("txid", tx.TxId).Msg("dup transaction treated as success")
		response.Result = true
		response.TxId = tx.TxId
	}

	// check the response
	if !response.Result {
		err = fmt.Errorf("failed to broadcast tx: %s - %s", response.Code, response.Message)
		c.logger.Err(err).Msg("")
		return "", err
	}

	err = c.signerCacheManager.SetSigned(
		txOutItem.CacheHash(),
		txOutItem.CacheVault(c.GetChain()),
		response.TxId,
	)
	if err != nil {
		c.logger.Err(err).
			Interface("tx_out_item", txOutItem).
			Msg("failed to mark tx out item as signed")
	}

	return response.TxId, nil
}

func (c *TronClient) ShouldReportSolvency(height int64) bool {
	return height%10 == 0
}

func (c *TronClient) ReportSolvency(height int64) error {
	if !c.ShouldReportSolvency(height) {
		return nil
	}

	return nil
}
