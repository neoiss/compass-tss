package evm

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	_ "embed"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	etypes "github.com/ethereum/go-ethereum/core/types"
	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

	vaultABI *abi.ABI
	erc20ABI *abi.ABI
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
	vaultABI, erc20ABI, err := evm.GetContractABI(gatewayContractABI, erc20ContractABI)
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
		vaultABI:             vaultABI,
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
	return e.cfg.MaxGasLimit, e.cfg.MaxSwapGasLimit, e.lastReportedGasPrice / 1e10 // 1e18 -> 1e8
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
func (e *EVMScanner) FetchTxs(height, chainHeight int64) (stypes.TxIn, error) {
	// log height every 100 blocks
	if height%100 == 0 {
		e.logger.Info().Int64("height", height).Msg("fetching txs for height")
	}

	logs, err := e.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(height),
		ToBlock:   big.NewInt(height),
		Addresses: []ecommon.Address{ecommon.HexToAddress(e.cfg.Mos)},
		Topics: [][]ecommon.Hash{{
			constants.EventOfDeposit.GetTopic(),
			constants.EventOfSwap.GetTopic(), // txIn -> voteTxIn
			constants.EventOfTransferOut.GetTopic(),
			constants.EventOfTransferAllowance.GetTopic(), // txOut -> voteTxOut
		}},
	})

	// process all transactions in the block
	block, err := e.ethRpc.GetBlock(height)
	if err != nil {
		return stypes.TxIn{}, err
	}
	txIn, err := e.processBlock(block, logs)
	if err != nil {
		e.logger.Error().Err(err).Int64("height", height).Msg("failed to search tx in block")
		return stypes.TxIn{}, fmt.Errorf("failed to process block: %d, err:%w", height, err)
	}

	e.currentBlockHeight = height
	// if reorgs are possible on this chain store block meta for handling
	if e.cfg.MaxReorgRescanBlocks > 0 {
		blockMeta := evmtypes.NewBlockMeta(block.Header(), txIn)
		if err = e.blockMetaAccessor.SaveBlockMeta(height, blockMeta); err != nil {
			e.logger.Err(err).Int64("height", height).Msg("fail to save block meta")
		}
		pruneHeight := height - e.cfg.MaxReorgRescanBlocks
		if pruneHeight > 0 {
			defer func() {
				if err = e.blockMetaAccessor.PruneBlockMeta(pruneHeight); err != nil {
					e.logger.Err(err).Int64("height", height).Msg("fail to prune block meta")
				}
			}()
		}
	}

	// skip reporting network fee and solvency if block more than flexibility blocks from tip
	if chainHeight-height > e.cfg.ObservationFlexibilityBlocks {
		return txIn, nil
	}

	// report network fee and solvency
	e.reportNetworkFee(height)
	if e.solvencyReporter != nil {
		if err = e.solvencyReporter(height); err != nil {
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

	// skip empty blocks
	if len(logs) == 0 {
		return txIn, nil
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

func (e *EVMScanner) getTxInOptimized(method string, block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	// Use custom method name as it varies between implementation
	// It should be akin to: "fetch all transaction receipts within a block", e.g. getBlockReceipts
	// This has shown to be more efficient than getTransactionReceipt in a batch call
	txInbound := stypes.TxIn{
		Chain:    e.cfg.ChainID,
		Filtered: false,
		MemPool:  false,
	}

	// tx lookup for compatibility with shared evm functions
	logByTxHash := make(map[string]*etypes.Log)
	for _, ele := range logs {
		// best effort remove the tx from the signed txs (ok if it does not exist)
		if err := e.blockMetaAccessor.RemoveSignedTxItem(ele.TxHash.String()); err != nil {
			e.logger.Err(err).Str("tx hash", ele.TxHash.String()).Msg("failed to remove signed tx item")
		}

		tmp := ele
		logByTxHash[ele.TxHash.String()] = &tmp
	}

	var receipts []*etypes.Receipt
	blockNumStr := hexutil.EncodeBig(block.Number())
	err := e.ethClient.Client().Call(
		&receipts,
		method,
		blockNumStr,
	)
	if err != nil {
		e.logger.Error().Err(err).Msg("failed to fetch block receipts")
		return stypes.TxIn{}, err
	}

	for _, receipt := range receipts {
		txForReceipt, ok := logByTxHash[receipt.TxHash.String()]
		if !ok {
			e.logger.Debug().
				Str("txHash", receipt.TxHash.String()).
				Uint64("blockNumber", block.NumberU64()).
				Msg("receipt not log, ignoring...")
			continue
		}

		// extract the txInItem
		var txInItem *stypes.TxInItem
		txInItem, err = e.receiptToTxInItem(txForReceipt, receipt)
		if err != nil {
			e.logger.Error().Err(err).Msg("failed to convert receipt to txInItem")
			continue
		}

		// skip invalid items
		if txInItem == nil {
			continue
		}
		if len(txInItem.To) == 0 {
			continue
		}
		if len([]byte(txInItem.Memo)) > constants.MaxMemoSize {
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

func (e *EVMScanner) getTxIn(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	// CHANGEME: if an EVM chain supports some way of fetching all transaction receipts
	// within a block, register it here.
	switch e.cfg.ChainID {
	case common.BASEChain, common.BSCChain:
		return e.getTxInOptimized("eth_getBlockReceipts", block, logs)
	}

	txInbound := stypes.TxIn{
		Chain:    e.cfg.ChainID,
		Filtered: false,
		MemPool:  false,
	}

	// collect all relevant transactions from the block into batches
	var rpcBatch []rpc.BatchElem
	for _, ele := range logs {
		// best effort remove the tx from the signed txs (ok if it does not exist)
		if err := e.blockMetaAccessor.RemoveSignedTxItem(ele.TxHash.String()); err != nil {
			e.logger.Err(err).Str("tx hash", ele.TxHash.String()).Msg("failed to remove signed tx item")
		}

		rpcBatch = append(rpcBatch, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{ele.TxHash.String()},
			Result: &etypes.Receipt{},
		})
	}

	// process all batches
	for idx, ele := range logs {
		// send the batch rpc request
		err := e.ethClient.Client().BatchCall(rpcBatch)
		if err != nil {
			e.logger.Error().Int("size", len(rpcBatch)).Err(err).Msg("failed to batch fetch transaction receipts")
			return stypes.TxIn{}, err
		}

		// process the batch rpc response
		batchEle := rpcBatch[idx]
		if batchEle.Error != nil {
			if !errors.Is(err, ethereum.NotFound) {
				e.logger.Error().Err(batchEle.Error).Msg("failed to fetch transaction receipt")
			}
			continue
		}

		// get the receipt
		receipt, ok := batchEle.Result.(*etypes.Receipt)
		if !ok {
			e.logger.Error().Msg("failed to cast to transaction receipt")
			continue
		}

		// extract the txInItem
		var txInItem *stypes.TxInItem
		tmp := ele
		txInItem, err = e.receiptToTxInItem(&tmp, receipt)
		if err != nil {
			e.logger.Error().Err(err).Msg("failed to convert receipt to txInItem")
			continue
		}

		// skip invalid items
		if txInItem == nil {
			continue
		}
		if len(txInItem.To) == 0 {
			continue
		}
		if len([]byte(txInItem.Memo)) > constants.MaxMemoSize {
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

// TODO: This is only used by unit tests now, but covers receiptToTxInItem internally -
// refactor so this logic is only in test code.
func (e *EVMScanner) getTxInItem(tx *etypes.Transaction) (*stypes.TxInItem, error) {
	//if tx == nil || tx.To() == nil {
	//	return nil, nil
	//}
	//
	//receipt, err := e.ethRpc.GetReceipt(tx.Hash().Hex())
	//if err != nil {
	//	if errors.Is(err, ethereum.NotFound) {
	//		return nil, nil
	//	}
	//	return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	//}
	//
	//return e.receiptToTxInItem(tx, receipt)
	return nil, nil
}

func (e *EVMScanner) receiptToTxInItem(ll *etypes.Log, receipt *etypes.Receipt) (*stypes.TxInItem, error) {
	if receipt.Status != etypes.ReceiptStatusSuccessful {
		e.logger.Debug().Stringer("txid", ll.TxHash).Uint64("status", receipt.Status).Msg("tx failed")

		// remove failed transactions from signer cache so they are retried
		if e.signerCacheManager != nil {
			e.signerCacheManager.RemoveSigned(ll.TxHash.String())
		}

		return nil, nil
		//return e.getTxInFromFailedTransaction(tx, receipt), nil
	}

	ret, err := e.getTxInFromSmartContract(ll, receipt, 0)
	if err != nil {
		return nil, err
	}
	return ret, nil

	//disableWhitelist, err := e.bridge.GetMimir(constants.EVMDisableContractWhitelist.String())
	//if err != nil {
	//	e.logger.Err(err).Msgf("fail to get %s", constants.EVMDisableContractWhitelist.String())
	//	disableWhitelist = 0
	//}

	//if disableWhitelist == 1 {
	//	// parse tx without whitelist
	//	destination := tx.To()
	//	isToVault, _ := e.pubkeyMgr.IsValidPoolAddress(destination.String(), e.cfg.ChainID)
	//
	//	switch {
	//	case isToVault:
	//		// Tx to a vault
	//		return e.getTxInFromTransaction(tx, receipt)
	//	case e.isToValidContractAddress(destination, true):
	//		// Deposit directly to router
	//		return e.getTxInFromSmartContract(tx, receipt, 0)
	//	case evm.IsSmartContractCall(tx, receipt):
	//		// Tx to a different contract, attempt to parse with max allowable logs
	//		return e.getTxInFromSmartContract(tx, receipt, int64(e.cfg.MaxContractTxLogs))
	//	default:
	//		// Tx to a non-contract or vault address
	//		return e.getTxInFromTransaction(tx, receipt)
	//	}
	//} else {
	//	// parse tx with whitelist
	//	if e.isToValidContractAddress(tx.To(), true) {
	//		return e.getTxInFromSmartContract(tx, receipt, 0)
	//	}
	//	return e.getTxInFromTransaction(tx, receipt)
	//}
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
				constants.EventOfDeposit.GetTopic(),
				constants.EventOfSwap.GetTopic(), // txIn -> voteTxIn
				constants.EventOfTransferOut.GetTopic(),
				constants.EventOfTransferAllowance.GetTopic(), // txOut -> voteTxOut
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

// reportNetworkFee reports current network fee to thornode
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

	// gas price to 1e8 from 1e18
	tcGasPrice := new(big.Int).Div(gasPrice, big.NewInt(1e10))

	// post to thorchain
	cId, _ := e.cfg.ChainID.ChainID()
	e.globalNetworkFeeQueue <- stypes.NetworkFee{
		ChainId:             cId,
		Height:              height,
		TransactionSize:     e.cfg.MaxGasLimit,
		TransactionSwapSize: e.cfg.MaxSwapGasLimit,
		TransactionRate:     tcGasPrice.Uint64(),
	}

	e.lastReportedGasPrice = gasPrice.Uint64()
}

// --------------------------------- parse transaction ---------------------------------

func (e *EVMScanner) getTxInFromTransaction(tx *etypes.Transaction, receipt *etypes.Receipt) (*stypes.TxInItem, error) {
	txInItem := &stypes.TxInItem{
		Tx: tx.Hash().Hex()[2:], // drop the "0x" prefix
	}

	sender, err := e.eipSigner.Sender(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}
	txInItem.Sender = strings.ToLower(sender.String())
	txInItem.To = tx.To().Bytes()

	// on native transactions the memo is hex encoded in the data field
	data := tx.Data()
	if len(data) > 0 {
		var memo []byte
		memo, err = hex.DecodeString(string(data))
		if err != nil {
			txInItem.Memo = string(data)
		} else {
			txInItem.Memo = string(memo)
		}
	}

	//nativeValue := e.tokenManager.ConvertAmount(evm.NativeTokenAddr, tx.Value())

	return txInItem, nil
}

// isToValidContractAddress this method make sure the transaction to address is to
// THORChain router or a whitelist address
func (e *EVMScanner) isToValidContractAddress(addr *ecommon.Address, includeWhiteList bool) bool {
	if addr == nil {
		return false
	}
	// get the smart contract used by thornode
	contractAddresses := e.pubkeyMgr.GetContracts(e.cfg.ChainID)
	if includeWhiteList {
		contractAddresses = append(contractAddresses, e.whitelistContracts...)
	}

	// combine the whitelist smart contract address
	for _, item := range contractAddresses {
		if strings.EqualFold(item.String(), addr.String()) {
			return true
		}
	}
	return false
}

// getTxInFromSmartContract returns txInItem
func (e *EVMScanner) getTxInFromSmartContract(ll *etypes.Log, receipt *etypes.Receipt, maxLogs int64) (*stypes.TxInItem, error) {
	txInItem := &stypes.TxInItem{
		Tx: ll.TxHash.String()[2:], // drop the "0x" prefix
	}
	//sender, err := e.eipSigner.Sender(tx)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get sender: %w", err)
	//}
	//txInItem.Sender = strings.ToLower(sender.String())
	// 1 is Transaction success state
	if receipt.Status != etypes.ReceiptStatusSuccessful {
		e.logger.Debug().Stringer("txid", ll.TxHash).Uint64("status", receipt.Status).Msg("tx failed")
		return nil, nil
	}
	p := evm.NewSmartContractLogParser(e.isToValidContractAddress,
		e.tokenManager.GetAssetFromTokenAddress,
		e.tokenManager.GetTokenDecimalsForTHORChain,
		e.tokenManager.ConvertAmount,
		e.vaultABI,
		e.cfg.ChainID.GetGasAsset(),
		maxLogs,
	)

	// txInItem will be changed in p.getTxInItem function, so if the function return an
	// error txInItem should be abandoned
	isVaultTransfer, err := p.GetTxInItem(ll, txInItem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logs, err: %w", err)
	}
	if isVaultTransfer {
		//contractAddresses := e.pubkeyMgr.GetContracts(e.cfg.ChainID)
		//isDirectlyToRouter := false
		//for _, item := range contractAddresses {
		//	if strings.EqualFold(item.String(), tx.To().String()) {
		//		isDirectlyToRouter = true
		//		break
		//	}
		//}
		//if isDirectlyToRouter {
		//	// it is important to keep this part outside the above loop, as when we do router
		//	// upgrade, which might generate multiple deposit event, along with tx that has
		//	// native value in it
		//	nativeValue := cosmos.NewUintFromBigInt(tx.Value())
		//	if !nativeValue.IsZero() {
		//		nativeValue = e.tokenManager.ConvertAmount(evm.NativeTokenAddr, tx.Value())
		//		if txInItem.Coins.GetCoin(e.cfg.ChainID.GetGasAsset()).IsEmpty() && !nativeValue.IsZero() {
		//			txInItem.Coins = append(txInItem.Coins, common.NewCoin(e.cfg.ChainID.GetGasAsset(), nativeValue))
		//		}
		//	}
		//}
	}
	e.logger.Debug().
		Str("tx hash", txInItem.Tx).
		Uint64("gas used", receipt.GasUsed).
		Uint64("tx status", receipt.Status).
		Msg("txInItem parsed from smart contract")

	//// under no circumstance EVM gas price will be less than 1 Gwei, unless it is in dev environment
	//txGasPrice := tx.GasPrice()
	//txInItem.Gas = common.MakeEVMGas(e.cfg.ChainID, txGasPrice, receipt.GasUsed, receipt.L1Fee)
	//if txInItem.Coins.IsEmpty() {
	//	return nil, nil
	//}
	return txInItem, nil
}

// getTxInFromFailedTransaction when a transaction failed due to out of gas, this method
// will check whether the transaction is an outbound it fake a txInItem if the failed
// transaction is an outbound , and report it back to thornode, thus the gas fee can be
// subsidised need to know that this will also cause the vault that send
// out the outbound to be slashed 1.5x gas it is for security purpose
func (e *EVMScanner) getTxInFromFailedTransaction(tx *etypes.Transaction, receipt *etypes.Receipt) *stypes.TxInItem {
	if receipt.Status == etypes.ReceiptStatusSuccessful {
		e.logger.Info().Str("hash", tx.Hash().String()).Msg("success transaction should not get into getTxInFromFailedTransaction")
		return nil
	}
	fromAddr, err := e.eipSigner.Sender(tx)
	if err != nil {
		e.logger.Err(err).Msg("failed to get from address")
		return nil
	}
	ok, cif := e.pubkeyMgr.IsValidPoolAddress(fromAddr.String(), e.cfg.ChainID)
	if !ok || cif.IsEmpty() {
		return nil
	}
	//txGasPrice := tx.GasPrice()
	txHash := tx.Hash().Hex()[2:]
	return &stypes.TxInItem{
		Tx: txHash,
		//Memo:   memo.NewOutboundMemo(common.TxID(txHash)).String(), todo memo
		Sender: strings.ToLower(fromAddr.String()),
		//To:     strings.ToLower(tx.To().String()),
		//Coins:  common.NewCoins(common.NewCoin(e.cfg.ChainID.GetGasAsset(), cosmos.NewUint(1))),
		//Gas:    common.MakeEVMGas(e.cfg.ChainID, txGasPrice, receipt.GasUsed, receipt.L1Fee),
	}
}
