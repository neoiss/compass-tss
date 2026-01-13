package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	_ "embed"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/tokenlist"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm/types"
	evmtypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/signercache"
	. "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/mapprotocol/compass-tss/x/aggregators"
)

////////////////////////////////////////////////////////////////////////////////////////
// EVMScanner
////////////////////////////////////////////////////////////////////////////////////////

type EVMScanner struct {
	cfg                   config.BifrostBlockScannerConfiguration
	logger                zerolog.Logger
	db                    blockscanner.ScannerStorage
	m                     *metrics.Metrics
	errCounter            *prometheus.CounterVec
	gasPriceChanged       bool
	gasPrice              *big.Int
	lastReportedGasPrice  uint64
	ethClient             *ethclient.Client
	ethRpc                *evm.EthRPC
	blockMetaAccessor     evm.BlockMetaAccessor
	globalErrataQueue     chan<- stypes.ErrataBlock
	globalNetworkFeeQueue chan<- stypes.NetworkFee
	bridge                shareTypes.Bridge
	pubkeyMgr             pubkeymanager.PubKeyValidator
	eipSigner             etypes.Signer
	currentBlockHeight    int64
	gasCache              []*big.Int
	solvencyReporter      SolvencyReporter
	whitelistTokens       []tokenlist.ERC20Token
	whitelistContracts    []common.Address
	signerCacheManager    *signercache.CacheManager
	tokenManager          *evm.TokenManager
	gatewayABI, erc20ABI  *abi.ABI
}

// NewEVMScanner create a new instance of EVMScanner.
func NewEVMScanner(cfg config.BifrostBlockScannerConfiguration,
	storage blockscanner.ScannerStorage,
	chainID *big.Int,
	ethClient *ethclient.Client,
	ethRpc *evm.EthRPC,
	bridge shareTypes.Bridge,
	m *metrics.Metrics,
	pubkeyMgr pubkeymanager.PubKeyValidator,
	solvencyReporter SolvencyReporter,
	signerCacheManager *signercache.CacheManager,
) (*EVMScanner, error) {
	// check required arguments
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	if m == nil {
		return nil, errors.New("metrics manager is nil")
	}
	if ethClient == nil {
		return nil, errors.New("ETH RPC client is nil")
	}
	if pubkeyMgr == nil {
		return nil, errors.New("pubkey manager is nil")
	}

	// set storage prefixes
	prefixBlockMeta := fmt.Sprintf("%s-blockmeta-", strings.ToLower(cfg.ChainID.String()))
	prefixSignedMeta := fmt.Sprintf("%s-signedtx-", strings.ToLower(cfg.ChainID.String()))
	prefixTokenMeta := fmt.Sprintf("%s-tokenmeta-", strings.ToLower(cfg.ChainID.String()))

	// create block meta accessor
	blockMetaAccessor, err := evm.NewLevelDBBlockMetaAccessor(
		prefixBlockMeta, prefixSignedMeta, storage.GetInternalDb(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create block meta accessor: %w", err)
	}

	// load ABIs
	gatewayABI, erc20ABI, err := evm.GetContractABI(gatewayContractABI, erc20ContractABI)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract abi: %w", err)
	}

	// load token list
	allTokens := tokenlist.GetEVMTokenList(cfg.ChainID).Tokens
	var whitelistTokens []tokenlist.ERC20Token
	for _, addr := range cfg.WhitelistTokens {
		// find matching token in token list
		found := false
		for _, tok := range allTokens {
			if strings.EqualFold(addr, tok.Address) {
				whitelistTokens = append(whitelistTokens, tok)
				found = true
				break
			}
		}

		// all whitelisted tokens must be in the chain token list
		if !found {
			return nil, fmt.Errorf("whitelist token %s not found in token list", addr)
		}
	}

	// create token manager - storage is scoped to chain so assets should not collide
	tokenManager, err := evm.NewTokenManager(
		storage.GetInternalDb(),
		prefixTokenMeta,
		cfg.ChainID.GetGasAsset(),
		defaultDecimals,
		cfg.HTTPRequestTimeout,
		whitelistTokens,
		ethClient,
		gatewayContractABI,
		erc20ContractABI,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create token helper: %w", err)
	}

	// store the token metadata for the chain gas asset
	err = tokenManager.SaveTokenMeta(
		cfg.ChainID.GetGasAsset().Symbol.String(), evm.NativeTokenAddr, defaultDecimals,
	)
	if err != nil {
		return nil, err
	}

	// load whitelist contracts for the chain
	whitelistContracts := []common.Address{}
	for _, agg := range aggregators.DexAggregators(common.LatestVersion) {
		if agg.Chain.Equals(cfg.ChainID) {
			whitelistContracts = append(whitelistContracts, common.Address(agg.Address))
		}
	}

	fmt.Println("e.cfg.ObservationFlexibilityBlocks -------------- ", cfg.ObservationFlexibilityBlocks)

	return &EVMScanner{
		cfg:                  cfg,
		logger:               log.Logger.With().Stringer("chain", cfg.ChainID).Logger(),
		errCounter:           m.GetCounterVec(metrics.BlockScanError(cfg.ChainID)),
		ethRpc:               ethRpc,
		ethClient:            ethClient,
		db:                   storage,
		m:                    m,
		gasPrice:             big.NewInt(0),
		lastReportedGasPrice: 0,
		gasPriceChanged:      false,
		blockMetaAccessor:    blockMetaAccessor,
		bridge:               bridge,
		gatewayABI:           gatewayABI,
		erc20ABI:             erc20ABI,
		eipSigner:            etypes.NewLondonSigner(chainID),
		pubkeyMgr:            pubkeyMgr,
		gasCache:             make([]*big.Int, 0),
		solvencyReporter:     solvencyReporter,
		whitelistTokens:      whitelistTokens,
		whitelistContracts:   whitelistContracts,
		signerCacheManager:   signerCacheManager,
		tokenManager:         tokenManager,
	}, nil
}

// --------------------------------- exported ---------------------------------

// GetGasPrice returns the current gas price.
func (e *EVMScanner) GetGasPrice() *big.Int {
	if e.cfg.FixedGasRate > 0 {
		return big.NewInt(e.cfg.FixedGasRate)
	}
	return e.gasPrice
}

// GetNetworkFee returns current chain network fee according to Bifrost.
func (e *EVMScanner) GetNetworkFee() (transactionSize, txSwapGasLimit, transactionFeeRate uint64) {
	return e.cfg.MaxGasLimit, e.cfg.MaxSwapGasLimit, e.lastReportedGasPrice
}

// GetNonce returns the nonce (including pending) for the given address.
func (e *EVMScanner) GetNonce(addr string) (uint64, error) {
	return e.ethRpc.GetNonce(addr)
}

// GetNonceFinalized returns the nonce for the given address.
func (e *EVMScanner) GetNonceFinalized(addr string) (uint64, error) {
	return e.ethRpc.GetNonceFinalized(addr)
}

// FetchMemPool returns all transactions in the mempool.
func (e *EVMScanner) FetchMemPool(_ int64) (stypes.TxIn, error) {
	return stypes.TxIn{}, nil
}

// GetTokens returns all token meta data.
func (e *EVMScanner) GetTokens() ([]*evmtypes.TokenMeta, error) {
	return e.tokenManager.GetTokens()
}

// FetchTxs extracts all relevant transactions from the block at the provided height.
func (e *EVMScanner) FetchTxs(currentHeight, latestHeight int64) (stypes.TxIn, error) {
	logs, err := e.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(currentHeight),
		ToBlock:   big.NewInt(currentHeight),
		Addresses: []ecommon.Address{ecommon.HexToAddress(e.cfg.Mos)},
		Topics: [][]ecommon.Hash{{
			constants.EventOfBridgeOut.GetTopic(), // txIn -> voteTxIn
			constants.EventOfBridgeIn.GetTopic(),  // txOut -> voteTxOut
		}},
	})
	if err != nil {
		return stypes.TxIn{}, err
	}

	// process all transactions in the block
	block, err := e.ethRpc.GetBlock(currentHeight)
	if err != nil {
		return stypes.TxIn{}, err
	}
	txIn, err := e.processBlock(block, logs)
	if err != nil {
		e.logger.Error().Err(err).Int64("currentHeight", currentHeight).Msg("failed to search tx in block")
		return stypes.TxIn{}, fmt.Errorf("failed to process block: %d, err:%w", currentHeight, err)
	}

	e.currentBlockHeight = currentHeight
	// if reorgs are possible on this chain store block meta for handling
	if e.cfg.MaxReorgRescanBlocks > 0 {
		blockMeta := evmtypes.NewBlockMeta(block.Header(), txIn)
		if err = e.blockMetaAccessor.SaveBlockMeta(currentHeight, blockMeta); err != nil {
			e.logger.Err(err).Int64("currentHeight", currentHeight).Msg("fail to save block meta")
		}
		pruneHeight := currentHeight - e.cfg.MaxReorgRescanBlocks
		if pruneHeight > 0 {
			defer func() {
				if err = e.blockMetaAccessor.PruneBlockMeta(pruneHeight); err != nil {
					e.logger.Err(err).Int64("currentHeight", currentHeight).Msg("fail to prune block meta")
				}
			}()
		}
	}

	// skip reporting network fee and solvency if block more than flexibility blocks from tip
	if latestHeight-currentHeight > e.cfg.ObservationFlexibilityBlocks {
		return txIn, nil
	}

	// report network fee and solvency
	e.reportNetworkFee(currentHeight)
	if e.solvencyReporter != nil {
		if err = e.solvencyReporter(currentHeight); err != nil {
			e.logger.Err(err).Msg("failed to report Solvency info to THORNode")
		}
	}

	return txIn, nil
}

// --------------------------------- extraction ---------------------------------

func (e *EVMScanner) processBlock(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	txIn := stypes.TxIn{
		Chain:    e.cfg.ChainID,
		TxArray:  nil,
		Filtered: false,
		MemPool:  false,
	}

	// collect gas prices of txs in current block
	var txsGas []*big.Int
	for _, tx := range block.Transactions() {
		txsGas = append(txsGas, tx.GasPrice())
	}
	e.updateGasPrice(txsGas)

	// process reorg if possible on this chain
	if e.cfg.MaxReorgRescanBlocks > 0 {
		reorgedTxIns, err := e.processReorg(block.Header())
		if err != nil {
			e.logger.Error().Err(err).Msgf("fail to process reorg for block %d", block.NumberU64())
			return txIn, err
		}
		if len(reorgedTxIns) > 0 {
			for _, item := range reorgedTxIns {
				if len(item.TxArray) == 0 {
					continue
				}
				txIn.TxArray = append(txIn.TxArray, item.TxArray...)
			}
		}
	}

	// skip empty blocks
	if len(logs) == 0 {
		return txIn, nil
	}
	// collect all relevant transactions from the block
	txInBlock, err := e.getTxIn(block, logs)
	if err != nil {
		return txIn, err
	}
	if len(txInBlock.TxArray) > 0 {
		txIn.TxArray = append(txIn.TxArray, txInBlock.TxArray...)
	}
	return txIn, nil
}

func (e *EVMScanner) getTxInOptimized(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	txInbound := stypes.TxIn{
		Chain: e.cfg.ChainID,
	}

	for _, ll := range logs {
		if err := e.blockMetaAccessor.RemoveSignedTxItem(ll.TxHash.String()); err != nil {
			e.logger.Err(err).Str("tx hash", ll.TxHash.String()).Msg("Failed to remove signed tx item")
		}

		// extract the txInItem
		var (
			err      error
			txInItem *stypes.TxInItem
		)
		tmp := ll
		txInItem, err = e.getTxInFromSmartContract(&tmp)
		if err != nil {
			e.logger.Error().Err(err).Msg("failed to convert receipt to txInItem")
			continue
		}

		// skip invalid items
		if txInItem == nil {
			continue
		}
		// if len(txInItem.To) == 0 {
		// 	continue
		// }
		// add the txInItem to the txInbound
		txInItem.Height = block.Number()
		txInbound.TxArray = append(txInbound.TxArray, txInItem)
	}

	if len(txInbound.TxArray) == 0 {
		e.logger.Debug().Uint64("block", block.NumberU64()).Msg("no tx need to be processed in this block")
		return stypes.TxIn{}, nil
	}
	return txInbound, nil
}

func (e *EVMScanner) getTxIn(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	// CHANGEME: if an EVM chain supports some way of fetching all transaction receipts
	// within a block, register it here.
	switch e.cfg.ChainID {
	case common.BASEChain, common.BSCChain:
		return e.getTxInOptimized(block, logs)
	}

	txInbound := stypes.TxIn{
		Chain:    e.cfg.ChainID,
		Filtered: false,
		MemPool:  false,
	}

	// process all batches
	for _, ele := range logs {
		// extract the txInItem
		var (
			err      error
			txInItem *stypes.TxInItem
		)
		if err := e.blockMetaAccessor.RemoveSignedTxItem(ele.TxHash.String()); err != nil {
			e.logger.Err(err).Str("tx hash", ele.TxHash.String()).Msg("failed to remove signed tx item")
		}

		tmp := ele
		txInItem, err = e.getTxInFromSmartContract(&tmp)
		if err != nil {
			e.logger.Error().Err(err).Msg("Failed to convert receipt to txInItem")
			continue
		}

		// skip invalid items
		if txInItem == nil {
			continue
		}

		// add the txInItem to the txInbound
		txInItem.Height = block.Number()
		txInbound.TxArray = append(txInbound.TxArray, txInItem)
	}

	if len(txInbound.TxArray) == 0 {
		e.logger.Debug().Uint64("block", block.NumberU64()).Msg("no tx need to be processed in this block")
		return stypes.TxIn{}, nil
	}
	return txInbound, nil
}

// --------------------------------- reorg ---------------------------------

// processReorg compares the block's parent hash with the stored block hash. When a
// reorg is detected, it triggers a rescan of all cached blocks in the reorg window.
// The function returns observations from the rescanned blocks.
func (e *EVMScanner) processReorg(header *etypes.Header) ([]stypes.TxIn, error) {
	previousHeight := header.Number.Int64() - 1
	prevBlockMeta, err := e.blockMetaAccessor.GetBlockMeta(previousHeight)
	if err != nil {
		return nil, fmt.Errorf("fail to get block meta of height(%d) : %w", previousHeight, err)
	}

	// skip re-org processing if we did not store the previous block meta
	if prevBlockMeta == nil {
		return nil, nil
	}

	// no re-org if stored block hash at previous height is equal to current block parent
	if strings.EqualFold(prevBlockMeta.BlockHash, header.ParentHash.Hex()) {
		return nil, nil
	}
	e.logger.Info().
		Int64("height", previousHeight).
		Str("stored_hash", prevBlockMeta.BlockHash).
		Str("current_parent_hash", header.ParentHash.Hex()).
		Msg("reorg detected")

	// send erratas and determine the block heights to rescan
	heights, err := e.reprocessTxs()
	if err != nil {
		e.logger.Err(err).Msg("fail to reprocess all txs")
	}

	// rescan heights
	var txIns []stypes.TxIn
	for _, rescanHeight := range heights {
		e.logger.Info().Msgf("rescan block height: %d", rescanHeight)
		var block *etypes.Block
		block, err = e.ethRpc.GetBlock(rescanHeight)
		if err != nil {
			e.logger.Err(err).Int64("height", rescanHeight).Msg("fail to get block")
			continue
		}
		logs, err := e.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: big.NewInt(rescanHeight),
			ToBlock:   big.NewInt(rescanHeight),
			Addresses: []ecommon.Address{ecommon.HexToAddress(e.cfg.Mos)},
			Topics: [][]ecommon.Hash{{
				constants.EventOfBridgeOut.GetTopic(), // txIn -> voteTxIn
				constants.EventOfBridgeIn.GetTopic(),  // txOut -> voteTxOut
			}},
		})
		if err != nil {
			e.logger.Err(err).Int64("height", rescanHeight).Msg("fail to get logs")
			continue
		}

		if len(logs) == 0 {
			e.logger.Debug().Int64("height", rescanHeight).Msg("processReorg: no logs")
			continue
		}
		var txIn stypes.TxIn
		txIn, err = e.getTxIn(block, logs)
		if err != nil {
			e.logger.Err(err).Int64("height", rescanHeight).Msg("fail to extract txs from block")
			continue
		}
		if len(txIn.TxArray) > 0 {
			txIns = append(txIns, txIn)
		}
	}
	return txIns, nil
}

// reprocessTx is initiated when the chain client detects a reorg. It reads block
// metadata from local storage and processes all transactions, sending an RPC request to
// check each transaction's existence. If a transaction no longer exists, the chain
// client reports this to Thorchain.
//
// The []int64 return value represents the block heights to be rescanned.
func (e *EVMScanner) reprocessTxs() ([]int64, error) {
	blockMetas, err := e.blockMetaAccessor.GetBlockMetas()
	if err != nil {
		return nil, fmt.Errorf("fail to get block metas from local storage: %w", err)
	}
	var rescanBlockHeights []int64
	for _, blockMeta := range blockMetas {
		metaTxs := make([]types.TransactionMeta, 0)
		var errataTxs []stypes.ErrataTx
		for _, tx := range blockMeta.Transactions {
			if e.ethRpc.CheckTransaction(tx.Hash) {
				e.logger.Debug().Msgf("height: %d, tx: %s still exists", blockMeta.Height, tx.Hash)
				metaTxs = append(metaTxs, tx)
				continue
			}

			// send an errata if the transactino no longer exists on chain
			errataTxs = append(errataTxs, stypes.ErrataTx{
				TxID:  common.TxID(tx.Hash),
				Chain: e.cfg.ChainID,
			})
		}
		if len(errataTxs) > 0 {
			e.globalErrataQueue <- stypes.ErrataBlock{
				Height: blockMeta.Height,
				Txs:    errataTxs,
			}
		}

		// fetch the header header to determine if the hash has changed and requires rescan
		var header *etypes.Header
		header, err = e.ethRpc.GetHeader(blockMeta.Height)
		if err != nil {
			e.logger.Err(err).
				Int64("height", blockMeta.Height).
				Msg("fail to get block header to check for reorg")

			// err on the side of caution and rescan the block
			rescanBlockHeights = append(rescanBlockHeights, blockMeta.Height)
			continue
		}

		// if the block hash is different than previously recorded, rescan the block
		if !strings.EqualFold(blockMeta.BlockHash, header.Hash().Hex()) {
			rescanBlockHeights = append(rescanBlockHeights, blockMeta.Height)
		}

		// save the updated block meta
		blockMeta.PreviousHash = header.ParentHash.Hex()
		blockMeta.BlockHash = header.Hash().Hex()
		blockMeta.Transactions = metaTxs
		if err = e.blockMetaAccessor.SaveBlockMeta(blockMeta.Height, blockMeta); err != nil {
			e.logger.Err(err).Int64("height", blockMeta.Height).Msg("fail to save block meta")
		}
	}
	return rescanBlockHeights, nil
}

// --------------------------------- gas ---------------------------------

// updateGasPrice calculates and stores the current gas price to reported to thornode
func (e *EVMScanner) updateGasPrice(prices []*big.Int) {
	// skip empty blocks
	if len(prices) == 0 {
		return
	}

	// find the median gas price in the block
	sort.Slice(prices, func(i, j int) bool { return prices[i].Cmp(prices[j]) == -1 })
	gasPrice := prices[len(prices)/2]

	// add to the cache
	e.gasCache = append(e.gasCache, gasPrice)
	if len(e.gasCache) > e.cfg.GasCacheBlocks {
		e.gasCache = e.gasCache[(len(e.gasCache) - e.cfg.GasCacheBlocks):]
	}

	// skip update unless cache is full
	if len(e.gasCache) < e.cfg.GasCacheBlocks {
		return
	}

	// compute the median of the median prices in the cache
	medians := []*big.Int{}
	medians = append(medians, e.gasCache...)
	sort.Slice(medians, func(i, j int) bool { return medians[i].Cmp(medians[j]) == -1 })
	median := medians[len(medians)/2]

	// round the price up to nearest configured resolution
	resolution := big.NewInt(e.cfg.GasPriceResolution)
	median.Add(median, new(big.Int).Sub(resolution, big.NewInt(1)))
	median = median.Div(median, resolution)
	median = median.Mul(median, resolution)
	e.gasPrice = median

	// record metrics
	gasPriceFloat, _ := new(big.Float).SetInt64(e.gasPrice.Int64()).Float64()
	e.m.GetGauge(metrics.GasPrice(e.cfg.ChainID)).Set(gasPriceFloat)
	e.m.GetCounter(metrics.GasPriceChange(e.cfg.ChainID)).Inc()
}

// reportNetworkFee reports current network fee to map
func (e *EVMScanner) reportNetworkFee(height int64) {
	gasPrice := e.GetGasPrice()

	// skip posting if there is not yet a fee
	if gasPrice.Cmp(big.NewInt(0)) == 0 {
		return
	}

	// skip fee if less than 1 resolution away from the last
	feeDelta := new(big.Int).Sub(gasPrice, big.NewInt(int64(e.lastReportedGasPrice)))
	feeDelta.Abs(feeDelta)
	if e.lastReportedGasPrice != 0 && feeDelta.Cmp(big.NewInt(e.cfg.GasPriceResolution)) != 1 {
		skip := true

		// every 100 blocks send the fee if none is set
		if height%100 == 0 {
			hasNetworkFee, err := e.bridge.HasNetworkFee(e.cfg.ChainID)
			skip = err != nil || hasNetworkFee
		}

		if skip {
			return
		}
	}

	// post to map
	cId, _ := e.cfg.ChainID.ChainID()
	e.globalNetworkFeeQueue <- stypes.NetworkFee{
		ChainId:             cId,
		Height:              height,
		TransactionSize:     e.cfg.MaxGasLimit,
		TransactionSwapSize: e.cfg.MaxSwapGasLimit,
		TransactionRate:     gasPrice.Uint64(),
	}

	e.lastReportedGasPrice = gasPrice.Uint64()
}

// --------------------------------- parse transaction ---------------------------------

// getTxInFromSmartContract returns txInItem
func (e *EVMScanner) getTxInFromSmartContract(ll *etypes.Log) (*stypes.TxInItem, error) {
	txInItem := &stypes.TxInItem{
		Tx:       ll.TxHash.String()[2:], // drop the "0x" prefix
		LogIndex: ll.Index,
		Topic:    ll.Topics[0].Hex(),
	}
	cId, _ := e.cfg.ChainID.ChainID()
	txInItem.FromChain = cId
	p := evm.NewSmartContractLogParser(e.gatewayABI)

	// txInItem will be changed in p.getTxInItem function, so if the function return an
	// error txInItem should be abandoned
	err := p.GetTxInItem(ll, txInItem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs, err: %w", err)
	}

	e.logger.Debug().Msg("TxInItem parsed from smart contract")

	return txInItem, nil
}
