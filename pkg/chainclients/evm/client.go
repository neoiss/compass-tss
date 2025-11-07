package evm

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"

	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/runners"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/signercache"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/mapprotocol/compass-tss/tss"
	tssp "github.com/mapprotocol/compass-tss/tss/go-tss/tss"
)

////////////////////////////////////////////////////////////////////////////////////////
// EVMClient
////////////////////////////////////////////////////////////////////////////////////////

// EVMClient is a generic client for interacting with EVM chains.
type EVMClient struct {
	logger                  zerolog.Logger
	cfg                     config.BifrostChainConfiguration
	localPubKey             common.PubKey
	kw                      *evm.KeySignWrapper
	ethClient               *ethclient.Client
	evmScanner              *EVMScanner
	bridge                  shareTypes.Bridge
	blockScanner            *blockscanner.BlockScanner
	gatewayAbi              *abi.ABI
	pubkeyMgr               pubkeymanager.PubKeyValidator
	poolMgr                 mapo.PoolManager
	tssKeySigner            *tss.KeySign
	wg                      *sync.WaitGroup
	stopchan                chan struct{}
	globalSolvencyQueue     chan stypes.Solvency
	signerCacheManager      *signercache.CacheManager
	lastSolvencyCheckHeight int64
}

// NewEVMClient creates a new EVMClient.
func NewEVMClient(
	relayKey *keys.Keys,
	cfg config.BifrostChainConfiguration,
	server *tssp.TssServer,
	bridge shareTypes.Bridge,
	m *metrics.Metrics,
	pubkeyMgr pubkeymanager.PubKeyValidator,
	poolMgr mapo.PoolManager,
) (*EVMClient, error) {
	// check required arguments
	if relayKey == nil {
		return nil, fmt.Errorf("failed to create EVM client, thor keys empty")
	}
	if bridge == nil {
		return nil, errors.New("thorchain bridge is nil")
	}
	if pubkeyMgr == nil {
		return nil, errors.New("pubkey manager is nil")
	}
	if poolMgr == nil {
		return nil, errors.New("pool manager is nil")
	}

	// create keys
	tssKm, err := tss.NewKeySign(server, bridge)
	if err != nil {
		return nil, fmt.Errorf("failed to create tss signer: %w", err)
	}
	priv, err := relayKey.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}
	temp, err := codec.ToCmtPubKeyInterface(priv.PubKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get tm pub key: %w", err)
	}
	pk, err := common.NewPubKeyFromCrypto(temp)
	if err != nil {
		return nil, fmt.Errorf("failed to get pub key: %w", err)
	}
	evmPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	clog := log.With().Str("module", "evm").Stringer("chain", cfg.ChainID).Logger()

	// create rpc client based on what authentication config is set
	var ethClient *ethclient.Client
	switch {
	case cfg.AuthorizationBearer != "":
		clog.Info().Msg("initializing evm client with bearer token")
		authFn := func(h http.Header) error {
			h.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.AuthorizationBearer))
			return nil
		}
		var rpcClient *rpc.Client
		rpcClient, err = rpc.DialOptions(
			context.Background(),
			cfg.RPCHost,
			rpc.WithHTTPAuth(authFn),
		)
		if err != nil {
			return nil, err
		}
		ethClient = ethclient.NewClient(rpcClient)

	case cfg.UserName != "" && cfg.Password != "":
		clog.Info().Msg("initializing evm client with http basic auth")

		authFn := func(h http.Header) error {
			auth := base64.StdEncoding.EncodeToString([]byte(cfg.UserName + ":" + cfg.Password))
			h.Set("Authorization", fmt.Sprintf("Basic %s", auth))
			return nil
		}
		var rpcClient *rpc.Client
		rpcClient, err = rpc.DialOptions(
			context.Background(),
			cfg.RPCHost,
			rpc.WithHTTPAuth(authFn),
		)
		if err != nil {
			return nil, err
		}
		ethClient = ethclient.NewClient(rpcClient)

	default:
		ethClient, err = ethclient.Dial(cfg.RPCHost)
		if err != nil {
			return nil, fmt.Errorf("fail to dial ETH rpc host(%s): %w", cfg.RPCHost, err)
		}
	}

	rpcClient, err := evm.NewEthRPC(
		ethClient,
		cfg.BlockScanner.HTTPRequestTimeout,
		cfg.ChainID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create ETH rpc host(%s): %w", cfg.RPCHost, err)
	}

	// get chain id
	chainID, err := getChainID(ethClient, cfg.BlockScanner.HTTPRequestTimeout)
	if err != nil {
		return nil, err
	}
	if chainID.Uint64() == 0 {
		return nil, fmt.Errorf("chain id is: %d , invalid", chainID.Uint64())
	}

	// create keysign wrapper
	keysignWrapper, err := evm.NewKeySignWrapper(evmPrivateKey, pk, tssKm, chainID, cfg.ChainID.String())
	if err != nil {
		return nil, fmt.Errorf("fail to create %s key sign wrapper: %w", cfg.ChainID, err)
	}

	// load vault abi
	vaultABI, _, err := evm.GetContractABI(gatewayContractABI, erc20ContractABI)
	if err != nil {
		return nil, fmt.Errorf("fail to get contract abi: %w", err)
	}

	c := &EVMClient{
		logger:       clog,
		cfg:          cfg,
		ethClient:    ethClient,
		localPubKey:  pk,
		kw:           keysignWrapper,
		bridge:       bridge,
		gatewayAbi:   vaultABI,
		pubkeyMgr:    pubkeyMgr,
		poolMgr:      poolMgr,
		tssKeySigner: tssKm,
		wg:           &sync.WaitGroup{},
		stopchan:     make(chan struct{}),
	}

	// initialize storage
	var path string // if not set later, will in memory storage
	if len(c.cfg.BlockScanner.DBPath) > 0 {
		path = fmt.Sprintf("%s/%s", c.cfg.BlockScanner.DBPath, c.cfg.BlockScanner.ChainID)
	}
	storage, err := blockscanner.NewBlockScannerStorage(path, c.cfg.ScannerLevelDB)
	if err != nil {
		return c, fmt.Errorf("fail to create blockscanner storage: %w", err)
	}
	signerCacheManager, err := signercache.NewSignerCacheManager(storage.GetInternalDb())
	if err != nil {
		return nil, fmt.Errorf("fail to create signer cache manager")
	}
	c.signerCacheManager = signerCacheManager

	// create block scanner
	c.evmScanner, err = NewEVMScanner(
		c.cfg.BlockScanner,
		storage,
		chainID,
		ethClient,
		rpcClient,
		c.bridge,
		m,
		pubkeyMgr,
		c.ReportSolvency,
		signerCacheManager,
	)
	if err != nil {
		return c, fmt.Errorf("fail to create evm block scanner: %w", err)
	}

	// initialize block scanner
	c.blockScanner, err = blockscanner.NewBlockScanner(
		c.cfg.BlockScanner, storage, m, c.bridge, c.evmScanner,
	)
	if err != nil {
		return c, fmt.Errorf("fail to create block scanner: %w", err)
	}

	// TODO: Is this necessary?
	localNodeAddress, err := c.localPubKey.GetAddress(cfg.ChainID)
	if err != nil {
		c.logger.Err(err).Stringer("chain", cfg.ChainID).Msg("failed to get local node address")
	}
	c.logger.Info().
		Stringer("chain", cfg.ChainID).
		Stringer("address", localNodeAddress).
		Msg("local node address")

	return c, nil
}

// Start starts the chain client with the given queues.
func (c *EVMClient) Start(
	globalTxsQueue chan stypes.TxIn,
	globalErrataQueue chan stypes.ErrataBlock,
	globalSolvencyQueue chan stypes.Solvency,
	globalNetworkFeeQueue chan stypes.NetworkFee,
) {
	c.evmScanner.globalErrataQueue = globalErrataQueue
	c.evmScanner.globalNetworkFeeQueue = globalNetworkFeeQueue
	c.globalSolvencyQueue = globalSolvencyQueue
	c.tssKeySigner.Start()
	c.blockScanner.Start(globalTxsQueue, globalNetworkFeeQueue)
	c.wg.Add(1)
	go c.unstuck()
	c.wg.Add(1)
	go runners.SolvencyCheckRunner(c.GetChain(), c, c.bridge, c.stopchan, c.wg, constants.MAPRelayChainBlockTime)
}

// Stop stops the chain client.
func (c *EVMClient) Stop() {
	c.tssKeySigner.Stop()
	c.blockScanner.Stop()
	close(c.stopchan)
	c.wg.Wait()
}

// IsBlockScannerHealthy returns true if the block scanner is healthy.
func (c *EVMClient) IsBlockScannerHealthy() bool {
	return c.blockScanner.IsHealthy()
}

// --------------------------------- config ---------------------------------

// GetConfig returns the chain configuration.
func (c *EVMClient) GetConfig() config.BifrostChainConfiguration {
	return c.cfg
}

// GetChain returns the chain.
func (c *EVMClient) GetChain() common.Chain {
	return c.cfg.ChainID
}

// --------------------------------- status ---------------------------------

// GetHeight returns the current height of the chain.
func (c *EVMClient) GetHeight() (int64, error) {
	return c.evmScanner.GetHeight()
}

// GetBlockScannerHeight returns blockscanner height
func (c *EVMClient) GetBlockScannerHeight() (int64, error) {
	return c.blockScanner.PreviousHeight(), nil
}

func (c *EVMClient) GetLatestTxForVault(vault string) (string, string, error) {
	lastObserved, err := c.signerCacheManager.GetLatestRecordedTx(stypes.InboundCacheKey(vault, c.GetChain().String()))
	if err != nil {
		return "", "", err
	}
	lastBroadCasted, err := c.signerCacheManager.GetLatestRecordedTx(stypes.BroadcastCacheKey(vault, c.GetChain().String()))
	return lastObserved, lastBroadCasted, err
}

// --------------------------------- addresses ---------------------------------

// GetAddress returns the address for the given public key.
func (c *EVMClient) GetAddress(poolPubKey common.PubKey) string {
	addr, err := poolPubKey.GetAddress(c.cfg.ChainID)
	if err != nil {
		c.logger.Error().Err(err).Str("pool_pub_key", poolPubKey.String()).Msg("fail to get pool address")
		return ""
	}
	return addr.String()
}

// GetAccount returns the account for the given public key.
func (c *EVMClient) GetAccount(pk common.PubKey, height *big.Int) (common.Account, error) {
	addr := c.GetAddress(pk)
	nonce, err := c.evmScanner.GetNonce(addr)
	if err != nil {
		return common.Account{}, err
	}
	coins, err := c.GetBalances(addr, height)
	if err != nil {
		return common.Account{}, err
	}
	account := common.NewAccount(int64(nonce), 0, coins, false)
	return account, nil
}

// GetAccountByAddress returns the account for the given address.
func (c *EVMClient) GetAccountByAddress(address string, height *big.Int) (common.Account, error) {
	nonce, err := c.evmScanner.GetNonce(address)
	if err != nil {
		return common.Account{}, err
	}
	coins, err := c.GetBalances(address, height)
	if err != nil {
		return common.Account{}, err
	}
	account := common.NewAccount(int64(nonce), 0, coins, false)
	return account, nil
}

func (c *EVMClient) getSmartContractAddr(pubkey common.PubKey) common.Address {
	return c.pubkeyMgr.GetContract(c.cfg.ChainID, pubkey)
}

func (c *EVMClient) getTokenAddressFromAsset(asset common.Asset) string {
	if asset.Equals(c.cfg.ChainID.GetGasAsset()) {
		return evm.NativeTokenAddr
	}
	allParts := strings.Split(asset.Symbol.String(), "-")
	return allParts[len(allParts)-1]
}

// --------------------------------- balances ---------------------------------

// GetBalance returns the balance of the provided address.
func (c *EVMClient) GetBalance(addr, token string, height *big.Int) (*big.Int, error) {
	contractAddresses := c.pubkeyMgr.GetContracts(c.cfg.ChainID)
	c.logger.Debug().Interface("contractAddresses", contractAddresses).Msg("got contracts")
	if len(contractAddresses) == 0 {
		return nil, fmt.Errorf("fail to get contract address")
	}

	return c.evmScanner.tokenManager.GetBalance(addr, token, height, contractAddresses[0].String())
}

// GetBalances returns the balances of the provided address.
func (c *EVMClient) GetBalances(addr string, height *big.Int) (common.Coins, error) {
	// for all the tokens the chain client has dealt with before
	tokens, err := c.evmScanner.GetTokens()
	if err != nil {
		return nil, fmt.Errorf("fail to get all the tokens: %w", err)
	}
	coins := common.Coins{}
	for _, token := range tokens {
		var balance *big.Int
		balance, err = c.GetBalance(addr, token.Address, height)
		if err != nil {
			c.logger.Err(err).Str("token", token.Address).Msg("fail to get balance for token")
			continue
		}
		asset := c.cfg.ChainID.GetGasAsset()
		if !strings.EqualFold(token.Address, evm.NativeTokenAddr) {
			asset, err = common.NewAsset(fmt.Sprintf("%s.%s-%s", c.GetChain(), token.Symbol, token.Address))
			if err != nil {
				return nil, err
			}
		}
		bal := c.evmScanner.tokenManager.ConvertAmount(token.Address, balance)
		coins = append(coins, common.NewCoin(asset, bal))
	}

	return coins.Distinct(), nil
}

// --------------------------------- gas ---------------------------------

// GetGasFee returns the gas fee based on the current gas price.
func (c *EVMClient) GetGasFee(gas uint64) common.Gas {
	return common.GetEVMGasFee(c.cfg.ChainID, c.GetGasPrice(), gas)
}

// GetGasPrice returns the current gas price.
func (c *EVMClient) GetGasPrice() *big.Int {
	gasPrice := c.evmScanner.GetGasPrice()
	return gasPrice
}

// --------------------------------- build transaction ---------------------------------

// getOutboundTxData generates the tx data and tx value of the outbound Router Contract call, and checks if the router contract has been updated
func (c *EVMClient) getOutboundTxData(sender common.Address, txOutItem stypes.TxOutItem) ([]byte, error) {
	var (
		err error
		ret []byte
	)
	switch txOutItem.Method {
	case constants.BridgeIn:
		ret, err = c.gatewayAbi.Pack(constants.BridgeIn,
			ecommon.HexToAddress(sender.String()), txOutItem.OrderId, txOutItem.Data, txOutItem.Signature)
	default:
		return nil, fmt.Errorf("not support method(%s)", txOutItem.Method)
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *EVMClient) buildOutboundTx(txOutItem stypes.TxOutItem, nonce uint64) (*etypes.Transaction, error) {
	fromAddr, err := c.localPubKey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return nil, err
	}
	cgl, err := evm.ParseChainAndGasLimit(ecommon.BytesToHash(common.Completion(txOutItem.ChainAndGasLimit.Bytes(), 32)))
	if err != nil {
		c.logger.Err(err).Msg("fail to parse chain and gas limit")
		return nil, err
	}

	txData, err := c.getOutboundTxData(fromAddr, txOutItem)
	if err != nil {
		return nil, fmt.Errorf("failed to get outbound tx data %w", err)
	}

	gasRate := c.GetGasPrice()
	if c.cfg.BlockScanner.FixedGasRate > 0 && gasRate.Cmp(big.NewInt(0)) == 0 {
		// if chain gas is zero we are still filling our gas price buffer, use outbound rate
		gasRate = big.NewInt(c.cfg.BlockScanner.FixedGasRate)
	} else if gasRate.Cmp(big.NewInt(0)) == 0 {
		// todo will next 2
	}

	if gasRate.Cmp(cgl.Third) != 0 {
		c.logger.Info().Str("inHash", txOutItem.InTxHash).
			Str("outboundRate", cgl.Third.String()).
			Str("currentRate", c.GetGasPrice().String()).
			Str("effectiveRate", gasRate.String()).
			Msg("gas rate")
	}

	// outbound tx always send to smart contract address
	createdTx := etypes.NewTransaction(nonce, ecommon.HexToAddress(c.cfg.BlockScanner.Mos), big.NewInt(0), c.cfg.BlockScanner.MaxGasLimit, gasRate, txData)
	estimatedGas, err := c.evmScanner.ethRpc.EstimateGas(fromAddr.String(), createdTx)
	if err != nil {
		c.logger.Err(err).Str("input", ecommon.Bytes2Hex(createdTx.Data())).
			Msg("fail to estimate gas")
		return nil, err
	}

	// if estimated gas is more than the planned gas, abort and let thornode reschedule
	if estimatedGas > cgl.End.Uint64() {
		c.logger.Warn().
			Str("in_hash", txOutItem.InTxHash).
			Stringer("rate", gasRate).
			Uint64("estimated_gas_units", estimatedGas).
			Uint64("max_gas_units", cgl.End.Uint64()).
			// Str("scheduled_max_fee", scheduledMaxFee.String()).
			Msg("max gas exceeded, aborting to let thornode reschedule")
		return nil, nil
	}

	// before signing, confirm the vault has enough gas asset
	estimatedFee := big.NewInt(int64(estimatedGas))
	estimatedFee.Mul(estimatedFee, gasRate)

	createdTx = etypes.NewTransaction(
		nonce,
		ecommon.HexToAddress(c.cfg.BlockScanner.Mos),
		big.NewInt(0),
		estimatedGas,
		gasRate,
		txData,
	)

	return createdTx, nil
}

// --------------------------------- sign ---------------------------------

// SignTx returns the signed transaction.
func (c *EVMClient) SignTx(tx stypes.TxOutItem, height int64) ([]byte, []byte, *stypes.TxInItem, error) {
	selfId, _ := c.cfg.ChainID.ChainID()
	if tx.Chain.Cmp(selfId) != 0 {
		return nil, nil, nil, fmt.Errorf("chain %s is not support by evm chain client", tx.Chain)
	}

	if c.signerCacheManager.HasSigned(tx.CacheHash()) {
		c.logger.Info().Interface("tx", tx).Msg("transaction signed before, ignore")
		return nil, nil, nil, nil
	}

	// the nonce is stored as the transaction checkpoint, if it is set deserialize it
	// so we only retry with the same nonce to avoid double spend
	var (
		nonce    uint64
		fromAddr common.Address
	)
	fromAddr, err := c.localPubKey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return nil, nil, nil, err
	}
	nonce, err = c.evmScanner.GetNonce(fromAddr.String())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to fetch account(%s) nonce: %w", fromAddr, err)
	}

	// abort signing if the pending nonce is too far in the future
	var finalizedNonce uint64
	finalizedNonce, err = c.evmScanner.GetNonceFinalized(fromAddr.String())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to fetch account(%s) finalized nonce: %w", fromAddr, err)
	}
	if (nonce - finalizedNonce) > c.cfg.MaxPendingNonces {
		c.logger.Warn().
			Uint64("nonce", nonce).
			Uint64("finalizedNonce", finalizedNonce).
			Msg("pending nonce too far in future")
		return nil, nil, nil, fmt.Errorf("pending nonce too far in future")
	}

	// serialize nonce for later
	nonceBytes, err := json.Marshal(nonce)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to marshal nonce: %w", err)
	}

	outboundTx, err := c.buildOutboundTx(tx, nonce)
	if err != nil {
		c.logger.Err(err).Msg("fail to build outbound tx")
		return nil, nil, nil, err
	}

	// if transaction is nil, abort to allow thornode reschedule
	if outboundTx == nil {
		return nil, nil, nil, nil
	}

	rawTx, err := c.sign(outboundTx, common.EmptyPubKey, height, tx) // tx.VaultPubKey
	if err != nil || len(rawTx) == 0 {
		return nil, nonceBytes, nil, fmt.Errorf("fail to sign message: %w", err)
	}

	signedTx := &etypes.Transaction{}
	if err = signedTx.UnmarshalJSON(rawTx); err != nil {
		return nil, rawTx, nil, fmt.Errorf("fail to unmarshal signed tx: %w", err)
	}

	var txIn *stypes.TxInItem
	return rawTx, nil, txIn, nil
}

// sign is design to sign a given message with keysign party and keysign wrapper
func (c *EVMClient) sign(tx *etypes.Transaction, poolPubKey common.PubKey, height int64, txOutItem stypes.TxOutItem) ([]byte, error) {
	rawBytes, err := c.kw.LocalSign(tx)
	if err == nil && rawBytes != nil {
		return rawBytes, nil
	}
	// var keysignError tss.KeysignError
	// if errors.As(err, &keysignError) {
	// 	if len(keysignError.Blame.BlameNodes) == 0 {
	// 		// TSS doesn't know which node to blame
	// 		return nil, fmt.Errorf("fail to sign tx: %w", err)
	// 	}
	// 	// key sign error forward the keysign blame to thorchain
	// 	txID, errPostKeysignFail := c.bridge.PostKeysignFailure(keysignError.Blame, height,
	// 		txOutItem.Memo, nil, common.EmptyPubKey) // txOutItem.Coins, txOutItem.VaultPubKey
	// 	if errPostKeysignFail != nil {
	// 		return nil, multierror.Append(err, errPostKeysignFail)
	// 	}
	// 	c.logger.Info().Str("tx_id", txID).Msg("post keysign failure to thorchain")
	// }
	return nil, fmt.Errorf("fail to sign tx: %w", err)
}

// --------------------------------- broadcast ---------------------------------

// BroadcastTx broadcasts the transaction and returns the transaction hash.
func (c *EVMClient) BroadcastTx(txOutItem stypes.TxOutItem, hexTx []byte) (string, error) {
	// decode the transaction
	tx := &etypes.Transaction{}
	if err := tx.UnmarshalJSON(hexTx); err != nil {
		return "", err
	}
	txID := tx.Hash().String()

	// get context with default timeout
	ctx, cancel := c.getTimeoutContext()
	defer cancel()

	// send the transaction
	if err := c.ethClient.SendTransaction(ctx, tx); !isAcceptableError(err) {
		c.logger.Error().Str("txid", txID).Err(err).Msg("failed to send transaction")
		return "", err
	}
	c.logger.Info().Str("memo", txOutItem.Memo).Str("txid", txID).Msg("broadcast tx")

	// update the signer cache
	if err := c.signerCacheManager.SetSigned(txOutItem.CacheHash(), txOutItem.CacheVault(c.GetChain()), txID); err != nil {
		c.logger.Err(err).Interface("txOutItem", txOutItem).Msg("fail to mark tx out item as signed")
	}

	blockHeight, err := c.bridge.GetBlockHeight()
	if err != nil {
		c.logger.Err(err).Msg("fail to get current THORChain block height")
		// at this point , the tx already broadcast successfully , don't return an error
		// otherwise will cause the same tx to retry
	} else if err = c.AddSignedTxItem(txID, blockHeight, string(common.EmptyPubKey), &txOutItem); err != nil { //txOutItem.VaultPubKey.String()
		c.logger.Err(err).Str("hash", txID).Msg("fail to add signed tx item")
	}

	return txID, nil
}

// --------------------------------- observe ---------------------------------

// OnObservedTxIn is called when a new observed tx is received.
func (c *EVMClient) OnObservedTxIn(txIn stypes.TxInItem, blockHeight int64) {
	//m, err := mem.ParseMemo(common.LatestVersion, txIn.Memo)
	//if err != nil {
	//	// Debug log only as ParseMemo error is expected for THORName inbounds.
	//	c.logger.Debug().Err(err).Str("memo", txIn.Memo).Msg("fail to parse memo")
	//	return
	//}
	//if !m.IsOutbound() {
	//	return
	//}
	//if m.GetTxID().IsEmpty() {
	//	return
	//}
	if err := c.signerCacheManager.SetSigned(txIn.CacheHash(c.GetChain(), txIn.Tx), txIn.CacheVault(c.GetChain()), txIn.Tx); err != nil {
		c.logger.Err(err).Msg("fail to update signer cache")
	}
}

// GetConfirmationCount returns the confirmation count for the given tx.
func (c *EVMClient) GetConfirmationCount(txIn stypes.TxIn) int64 {
	switch c.cfg.ChainID {
	case common.AVAXChain: // instant finality
		return 0
	case common.BASEChain:
		return 12 // ~2 Ethereum blocks for parity with the 2 block minimum in eth client
	case common.BSCChain:
		return 3 // round up from 2.5 blocks required for finality
	default:
		c.logger.Fatal().Msgf("unsupported chain: %s", c.cfg.ChainID)
		return 0
	}
}

// ConfirmationCountReady returns true if the confirmation count is ready.
func (c *EVMClient) ConfirmationCountReady(txIn stypes.TxIn) bool {
	switch c.cfg.ChainID {
	case common.AVAXChain: // instant finality
		return true
	case common.BSCChain:
		if len(txIn.TxArray) == 0 {
			return true
		}
		blockHeight := txIn.TxArray[0].Height.Int64()
		confirm := txIn.ConfirmationRequired
		c.logger.Info().Msgf("confirmation required: %d", confirm)
		return (c.evmScanner.currentBlockHeight - blockHeight) >= confirm
	case common.BASEChain:
		// block is already finalized(settled to l1)
		return true
	default:
		c.logger.Fatal().Msgf("unsupported chain: %s", c.cfg.ChainID)
		return false
	}
}

// --------------------------------- solvency ---------------------------------

// ReportSolvency reports solvency once per configured solvency blocks.
func (c *EVMClient) ReportSolvency(height int64) error {
	if !c.ShouldReportSolvency(height) {
		return nil
	}

	// when block scanner is not healthy, only report from auto-unhalt SolvencyCheckRunner
	// (FetchTxs passes currentBlockHeight, while SolvencyCheckRunner passes chainHeight)
	if !c.IsBlockScannerHealthy() && height == c.evmScanner.currentBlockHeight {
		return nil
	}

	// // todo this report dont need
	// // fetch all asgard vaults
	// asgardVaults, err := c.bridge.GetAsgards()
	// if err != nil {
	// 	return fmt.Errorf("fail to get asgards, err: %w", err)
	// }

	// currentGasFee := cosmos.NewUint(3 * c.cfg.BlockScanner.MaxGasLimit * c.evmScanner.lastReportedGasPrice)

	// // report insolvent asgard vaults,
	// // or else all if the chain is halted and all are solvent
	// msgs := make([]stypes.Solvency, 0, len(asgardVaults))
	// solventMsgs := make([]stypes.Solvency, 0, len(asgardVaults))
	// for i := range asgardVaults {
	// 	var acct common.Account
	// 	acct, err = c.GetAccount(asgardVaults[i].PubKey, new(big.Int).SetInt64(height))
	// 	if err != nil {
	// 		c.logger.Err(err).Msg("fail to get account balance")
	// 		continue
	// 	}

	// 	msg := stypes.Solvency{
	// 		Height: height,
	// 		Chain:  c.cfg.ChainID,
	// 		PubKey: asgardVaults[i].PubKey,
	// 		Coins:  acct.Coins,
	// 	}

	// 	if runners.IsVaultSolvent(acct, asgardVaults[i], currentGasFee) {
	// 		solventMsgs = append(solventMsgs, msg) // Solvent-vault message
	// 		continue
	// 	}
	// 	msgs = append(msgs, msg) // Insolvent-vault message
	// }

	// // Only if the block scanner is unhealthy (e.g. solvency-halted) and all vaults are solvent,
	// // report that all the vaults are solvent.
	// // If there are any insolvent vaults, report only them.
	// // Not reporting both solvent and insolvent vaults is to avoid noise (spam):
	// // Reporting both could halt-and-unhalt SolvencyHalt in the same THOR block
	// // (resetting its height), plus making it harder to know at a glance from solvency reports which vaults were insolvent.
	// solvent := false
	// if !c.IsBlockScannerHealthy() && len(solventMsgs) == len(asgardVaults) {
	// 	msgs = solventMsgs
	// 	solvent = true
	// }

	// for i := range msgs {
	// 	c.logger.Info().
	// 		Stringer("asgard", msgs[i].PubKey).
	// 		Interface("coins", msgs[i].Coins).
	// 		Bool("solvent", solvent).
	// 		Msg("reporting solvency")

	// 	// send solvency to map via global queue consumed by the observer
	// 	select {
	// 	case c.globalSolvencyQueue <- msgs[i]:
	// 	case <-time.After(constants.MAPRelayChainBlockTime):
	// 		c.logger.Info().Msg("fail to send solvency info to thorchain, timeout")
	// 	}
	// }
	// c.lastSolvencyCheckHeight = height
	return nil
}

// ShouldReportSolvency returns true if the given height is a solvency report height.
func (c *EVMClient) ShouldReportSolvency(height int64) bool {
	return height%c.cfg.SolvencyBlocks == 0
}

// --------------------------------- helpers ---------------------------------

func (c *EVMClient) getTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.cfg.BlockScanner.HTTPRequestTimeout)
}
