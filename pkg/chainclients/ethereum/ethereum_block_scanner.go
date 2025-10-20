package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"

	"github.com/mapprotocol/compass-tss/blockscanner"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	tokenlist "github.com/mapprotocol/compass-tss/common/tokenlist"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/signercache"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
)

// SolvencyReporter is to report solvency info to THORNode
type SolvencyReporter func(int64) error

const (
	ethToken        = "0x0000000000000000000000000000000000000000"
	symbolMethod    = "symbol"
	decimalMethod   = "decimals"
	defaultDecimals = 18 // on ETH , consolidate all decimals to 18, in Wei
	tenGwei         = 10000000000
	// prefixTokenMeta declares prefix to use in leveldb to avoid conflicts
	prefixTokenMeta    = `eth-tokenmeta-` // nolint gosec:G101 not a hardcoded credential
	prefixBlockMeta    = `eth-blockmeta-`
	prefixSignedTxItem = `signed-txitem-`
)

// ETHScanner is a scanner that understand how to interact with ETH chain ,and scan block , parse smart contract etc
type ETHScanner struct {
	cfg                                            config.BifrostBlockScannerConfiguration
	logger                                         zerolog.Logger
	db                                             blockscanner.ScannerStorage
	m                                              *metrics.Metrics
	errCounter                                     *prometheus.CounterVec
	gasPriceChanged                                bool
	gasPrice                                       *big.Int
	lastReportedGasPrice                           uint64
	client                                         *ethclient.Client
	blockMetaAccessor                              evm.BlockMetaAccessor
	globalErrataQueue                              chan<- stypes.ErrataBlock
	globalNetworkFeeQueue                          chan<- stypes.NetworkFee
	gatewayABI                                     *abi.ABI
	erc20ABI                                       *abi.ABI              // todo
	tokens                                         *evm.LevelDBTokenMeta // todo
	bridge                                         shareTypes.Bridge
	pubkeyMgr                                      pubkeymanager.PubKeyValidator
	eipSigner                                      etypes.Signer
	currentBlockHeight, reqTime, cacheLatestHeight int64
	gasCache                                       []*big.Int
	solvencyReporter                               SolvencyReporter
	whitelistTokens                                []tokenlist.ERC20Token
	signerCacheManager                             *signercache.CacheManager
}

// NewETHScanner create a new instance of ETHScanner
func NewETHScanner(cfg config.BifrostBlockScannerConfiguration,
	storage blockscanner.ScannerStorage,
	chainID *big.Int,
	client *ethclient.Client,
	bridge shareTypes.Bridge,
	m *metrics.Metrics,
	pubkeyMgr pubkeymanager.PubKeyValidator,
	solvencyReporter SolvencyReporter,
	signerCacheManager *signercache.CacheManager,
) (*ETHScanner, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	if m == nil {
		return nil, errors.New("metrics manager is nil")
	}
	if client == nil {
		return nil, errors.New("ETH client is nil")
	}
	if pubkeyMgr == nil {
		return nil, errors.New("pubkey manager is nil")
	}
	blockMetaAccessor, err := evm.NewLevelDBBlockMetaAccessor(prefixBlockMeta, prefixSignedTxItem, storage.GetInternalDb())
	if err != nil {
		return nil, fmt.Errorf("fail to create block meta accessor: %w", err)
	}
	tokens, err := evm.NewLevelDBTokenMeta(storage.GetInternalDb(), prefixTokenMeta)
	if err != nil {
		return nil, fmt.Errorf("fail to create token meta db: %w", err)
	}
	err = tokens.SaveTokenMeta("ETH", ethToken, defaultDecimals)
	if err != nil {
		return nil, err
	}
	gatewayABI, erc20ABI, err := evm.GetContractABI(gatewayContractABI, erc20ContractABI)
	if err != nil {
		return nil, fmt.Errorf("fail to create contract abi: %w", err)
	}

	return &ETHScanner{
		cfg:                  cfg,
		logger:               log.Logger.With().Str("module", "block_scanner").Str("chain", common.ETHChain.String()).Logger(),
		errCounter:           m.GetCounterVec(metrics.BlockScanError(common.ETHChain)),
		client:               client,
		db:                   storage,
		m:                    m,
		gasPrice:             big.NewInt(initialGasPrice),
		lastReportedGasPrice: 0,
		gasPriceChanged:      false,
		blockMetaAccessor:    blockMetaAccessor,
		tokens:               tokens,
		bridge:               bridge,
		gatewayABI:           gatewayABI,
		erc20ABI:             erc20ABI,
		eipSigner:            etypes.NewLondonSigner(chainID),
		pubkeyMgr:            pubkeyMgr,
		gasCache:             make([]*big.Int, 0),
		solvencyReporter:     solvencyReporter,
		whitelistTokens:      tokenlist.GetETHTokenList().Tokens,
		signerCacheManager:   signerCacheManager,
	}, nil
}

// GetGasPrice returns current gas price
func (e *ETHScanner) GetGasPrice() *big.Int {
	if e.cfg.FixedGasRate > 0 {
		return big.NewInt(e.cfg.FixedGasRate)
	}
	return e.gasPrice
}

func (e *ETHScanner) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), e.cfg.HTTPRequestTimeout)
}

// GetHeight return latest block height
func (e *ETHScanner) GetHeight() (int64, error) {
	// add cache
	if time.Now().Unix()-e.reqTime < 5 { // request every 5 seconds
		return e.cacheLatestHeight, nil
	}
	ctx, cancel := e.getContext()
	defer cancel()
	height, err := e.client.BlockNumber(ctx)
	if err != nil {
		return -1, fmt.Errorf("fail to get block height: %w", err)
	}
	e.cacheLatestHeight = int64(height)
	e.reqTime = time.Now().Unix()
	return int64(height), nil
}

// GetNetworkFee returns current chain network fee according to Bifrost.
func (e *ETHScanner) GetNetworkFee() (transactionSize, transactionSwapSize, transactionFeeRate uint64) {
	return e.cfg.MaxGasLimit, e.cfg.MaxSwapGasLimit, e.lastReportedGasPrice
}

// FetchMemPool get tx from mempool
func (e *ETHScanner) FetchMemPool(_ int64) (stypes.TxIn, error) {
	return stypes.TxIn{}, nil
}

// GetTokens return all the token meta data
func (e *ETHScanner) GetTokens() ([]*types.TokenMeta, error) {
	return e.tokens.GetTokens()
}

// FetchTxs query the ETH chain to get txs in the given block height
func (e *ETHScanner) FetchTxs(currentHeight, latestHeight int64) (stypes.TxIn, error) {
	block, err := e.getRPCBlock(currentHeight)
	if err != nil {
		return stypes.TxIn{}, err
	}
	logs, err := e.getRPCFilterLogs(ethereum.FilterQuery{
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

	txIn, err := e.processBlock(block, logs)
	if err != nil {
		e.logger.Error().Err(err).Int64("currentHeight", currentHeight).Msg("fail to search tx in block")
		return stypes.TxIn{}, fmt.Errorf("fail to process block: %d, err:%w", currentHeight, err)
	}
	// blockMeta need to be saved , even there is no transactions found on this block at the time of scan
	// because at the time of scan , so the block hash will be stored, and it can be used to detect re-org
	blockMeta := types.NewBlockMeta(block.Header(), txIn)
	if err = e.blockMetaAccessor.SaveBlockMeta(blockMeta.Height, blockMeta); err != nil {
		e.logger.Err(err).Msgf("fail to save block meta of currentHeight: %d ", blockMeta.Height)
	}

	e.currentBlockHeight = currentHeight
	pruneHeight := currentHeight - e.cfg.MaxReorgRescanBlocks
	if pruneHeight > 0 {
		defer func() {
			if err := e.blockMetaAccessor.PruneBlockMeta(pruneHeight); err != nil {
				e.logger.Err(err).Msgf("fail to prune block meta, currentHeight(%d)", pruneHeight)
			}
		}()
	}

	// skip reporting network fee and solvency if block more than flexibility blocks from tip
	if latestHeight-currentHeight > e.cfg.ObservationFlexibilityBlocks {
		return txIn, nil
	}

	// gas price to 1e8 from 1e18
	gasPrice := e.GetGasPrice()
	tcGasPrice := new(big.Int).Div(gasPrice, big.NewInt(1e10)).Uint64()
	if tcGasPrice == 0 {
		tcGasPrice = 1
	}

	// post to thorchain if there is a fee and it has changed
	if gasPrice.Cmp(big.NewInt(0)) != 0 && tcGasPrice != e.lastReportedGasPrice {
		cId, _ := common.ETHChain.ChainID()
		e.globalNetworkFeeQueue <- stypes.NetworkFee{
			ChainId:             cId,
			Height:              currentHeight,
			TransactionSize:     e.cfg.MaxGasLimit,
			TransactionSwapSize: e.cfg.MaxSwapGasLimit,
			TransactionRate:     tcGasPrice,
		}

		e.lastReportedGasPrice = tcGasPrice
	}

	if e.solvencyReporter != nil {
		if err = e.solvencyReporter(currentHeight); err != nil {
			e.logger.Err(err).Msg("fail to report Solvency info to THORNode")
		}
	}
	return txIn, nil
}

// updateGasPrice records base fee + 25th percentile priority fee, rounded up 10 gwei.
func (e *ETHScanner) updateGasPrice(baseFee *big.Int, priorityFees []*big.Int) {
	// skip empty blocks
	if len(priorityFees) == 0 {
		return
	}

	// find the 25th percentile priority fee in the block
	sort.Slice(priorityFees, func(i, j int) bool { return priorityFees[i].Cmp(priorityFees[j]) == -1 }) //
	priorityFee := priorityFees[len(priorityFees)/4]                                                    //

	// consider gas price as base fee + 25th percentile priority fee
	gasPriceWei := new(big.Int).Add(baseFee, priorityFee) // 20000

	// round the price up to nearest configured resolution
	resolution := big.NewInt(e.cfg.GasPriceResolution)                        // 10000
	gasPriceWei.Add(gasPriceWei, new(big.Int).Sub(resolution, big.NewInt(1))) // 20000 + 9999
	gasPriceWei = gasPriceWei.Div(gasPriceWei, resolution)                    // 29999 / 9999 = 3
	gasPriceWei = gasPriceWei.Mul(gasPriceWei, resolution)                    // 3 * 9999

	// add to the cache
	e.gasCache = append(e.gasCache, gasPriceWei)
	if len(e.gasCache) > e.cfg.GasCacheBlocks {
		e.gasCache = e.gasCache[(len(e.gasCache) - e.cfg.GasCacheBlocks):]
	}

	e.updateGasPriceFromCache() //
}

func (e *ETHScanner) updateGasPriceFromCache() {
	// skip update unless cache is full
	if len(e.gasCache) < e.cfg.GasCacheBlocks {
		return
	}

	// compute the mean of cache
	sum := new(big.Int)
	for _, fee := range e.gasCache {
		sum.Add(sum, fee)
	}
	// avg
	mean := new(big.Int).Quo(sum, big.NewInt(int64(e.cfg.GasCacheBlocks)))

	// compute the standard deviation of cache
	// 标准值
	std := new(big.Int)
	for _, fee := range e.gasCache {
		v := new(big.Int).Sub(fee, mean) // 每个值在减去平均值, 4 - 2
		v.Mul(v, v)                      // 2*2
		std.Add(std, v)                  // 4 +4 +4+ 4  16
	}
	std.Quo(std, big.NewInt(int64(e.cfg.GasCacheBlocks))) // 在除以缓存长度 16/4 = 4
	std.Sqrt(std)                                         // std开根号 2

	// mean + 3x standard deviation over cache blocks
	// 2 + 2*3 = 8
	e.gasPrice = mean.Add(mean, std.Mul(std, big.NewInt(3)))

	// record metrics
	gasPriceFloat, _ := new(big.Float).SetInt64(e.gasPrice.Int64()).Float64()
	e.m.GetGauge(metrics.GasPrice(common.ETHChain)).Set(gasPriceFloat)
	e.m.GetCounter(metrics.GasPriceChange(common.ETHChain)).Inc()
}

// processBlock extracts transactions from block
func (e *ETHScanner) processBlock(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	height := int64(block.NumberU64())
	txIn := stypes.TxIn{
		Chain:    common.ETHChain,
		TxArray:  nil,
		Filtered: false,
		MemPool:  false,
	}

	// update gas price
	var priorityFees []*big.Int
	for _, tx := range block.Transactions() {
		tipCap := tx.GasTipCap()
		if tipCap == nil {
			tipCap = big.NewInt(0)
		}
		priorityFees = append(priorityFees, tipCap)
	}
	e.updateGasPrice(block.BaseFee(), priorityFees)

	reorgedTxIns, err := e.processReorg(block.Header())
	if err != nil {
		e.logger.Error().Err(err).Msgf("fail to process reorg for block %d", height)
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

	if len(logs) == 0 {
		return txIn, nil
	}

	txInBlock, err := e.extractTxs(block, logs)
	if err != nil {
		return txIn, err
	}
	if len(txInBlock.TxArray) > 0 {
		txIn.TxArray = append(txIn.TxArray, txInBlock.TxArray...)
	}
	return txIn, nil
}

func (e *ETHScanner) extractTxs(block *etypes.Block, logs []etypes.Log) (stypes.TxIn, error) {
	txInbound := stypes.TxIn{
		Chain:    common.ETHChain,
		Filtered: false,
		MemPool:  false,
	}

	sem := semaphore.NewWeighted(e.cfg.Concurrency)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	processTx := func(ll *etypes.Log) {
		defer wg.Done()
		if err := sem.Acquire(context.Background(), 1); err != nil {
			e.logger.Err(err).Msg("fail to acquire semaphore")
			return
		}
		defer sem.Release(1)

		// just try to remove the transaction hash from key value store
		// it doesn't matter whether the transaction is ours or not , success or failure
		// as long as the transaction id matches
		if err := e.blockMetaAccessor.RemoveSignedTxItem(ll.TxHash.Hex()); err != nil { // todo more log in on e tx
			e.logger.Err(err).Msgf("fail to remove signed tx item, hash:%s", ll.TxHash.Hex())
		}

		txInItem, err := e.fromTxToTxInLog(ll)
		if err != nil {
			e.logger.Error().Err(err).Str("hash", ll.TxHash.Hex()).Msg("fail to get one tx from server")
			return
		}
		if txInItem == nil {
			return
		}
		// sometimes if a transaction failed due to gas problem , it will have no `to` address
		if len(txInItem.To) == 0 {
			return
		}

		//txInItem.BlockHeight = block.Number().Int64()
		mu.Lock()
		txInbound.TxArray = append(txInbound.TxArray, txInItem)
		mu.Unlock()
		e.logger.Debug().Str("hash", ll.TxHash.Hex()).Msgf("%s got %d tx", e.cfg.ChainID, 1)
	}

	for _, ll := range logs {
		wg.Add(1)
		tmp := ll
		go processTx(&tmp)
	}

	wg.Wait()

	count := len(txInbound.TxArray)
	if count == 0 {
		e.logger.Info().Int64("block", int64(block.NumberU64())).Msg("No tx need to be processed in this block")
		return stypes.TxIn{}, nil
	}
	e.logger.Debug().Int64("block", int64(block.NumberU64())).Msgf("There are %d tx in this block need to process", count)
	return txInbound, nil
}

func (e *ETHScanner) onObservedTxIn(txIn stypes.TxInItem, blockHeight int64) {
	blockMeta, err := e.blockMetaAccessor.GetBlockMeta(blockHeight)
	if err != nil {
		e.logger.Err(err).Msgf("fail to get block meta on block height(%d)", blockHeight)
		return
	}

	if blockMeta == nil {
		e.logger.Error().Msgf("block meta for height:%d is nil", blockHeight)
		return
	}
	for _, item := range blockMeta.Transactions {
		if item.Hash == txIn.Tx {
			return
		}
	}

	blockMeta.Transactions = append(blockMeta.Transactions, types.TransactionMeta{
		Hash:        txIn.Tx,
		BlockHeight: blockHeight,
	})
	if err = e.blockMetaAccessor.SaveBlockMeta(blockHeight, blockMeta); err != nil {
		e.logger.Err(err).Msgf("fail to save block meta to storage,block height(%d)", blockHeight)
	}
}

// processReorg will compare block's parent hash and the block hash we have in store
// when there is a reorg detected , it will return true, other false
func (e *ETHScanner) processReorg(block *etypes.Header) ([]stypes.TxIn, error) {
	previousHeight := block.Number.Int64() - 1
	prevBlockMeta, err := e.blockMetaAccessor.GetBlockMeta(previousHeight)
	if err != nil {
		return nil, fmt.Errorf("fail to get block meta of height(%d) : %w", previousHeight, err)
	}
	if prevBlockMeta == nil {
		return nil, nil
	}
	// the block's previous hash need to be the same as the block hash chain client recorded in block meta
	// blockMetas[PreviousHeight].BlockHash == Block.PreviousHash
	if strings.EqualFold(prevBlockMeta.BlockHash, block.ParentHash.Hex()) {
		return nil, nil
	}
	e.logger.Info().Msgf("re-org detected, current block height:%d ,previous block hash is : %s , however block meta at height: %d, block hash is %s", block.Number.Int64(), block.ParentHash.Hex(), prevBlockMeta.Height, prevBlockMeta.BlockHash)
	heights, err := e.reprocessTxs()
	if err != nil {
		e.logger.Err(err).Msg("fail to reprocess all txs")
	}
	var txIns []stypes.TxIn
	for _, item := range heights {
		e.logger.Info().Msgf("rescan block height: %d", item)
		var block *etypes.Block
		block, err = e.getRPCBlock(item)
		if err != nil {
			e.logger.Err(err).Msgf("fail to get block from RPC endpoint, height:%d", item)
			continue
		}
		if block.Transactions().Len() == 0 {
			continue
		}
		var txIn stypes.TxIn
		txIn, err = e.extractTxs(block, nil)
		if err != nil {
			e.logger.Err(err).Msgf("fail to extract txs from block (%d)", item)
			continue
		}
		if len(txIn.TxArray) > 0 {
			txIns = append(txIns, txIn)
		}
	}
	return txIns, nil
}

// reprocessTx will be kicked off only when chain client detected a re-org on ethereum chain
// it will read through all the block meta data from local storage, and go through all the txs.
// For each transaction, it will send a RPC request to ethereuem chain, double check whether the TX exist or not
// if the tx still exist, then it is all good, if a transaction previous we detected, however doesn't exist anymore, that means
// the transaction had been removed from chain, chain client should report to thorchain
// []int64 is the block heights that need to be rescanned
func (e *ETHScanner) reprocessTxs() ([]int64, error) {
	blockMetas, err := e.blockMetaAccessor.GetBlockMetas()
	if err != nil {
		return nil, fmt.Errorf("fail to get block metas from local storage: %w", err)
	}
	var rescanBlockHeights []int64
	for _, blockMeta := range blockMetas {
		metaTxs := make([]types.TransactionMeta, 0)
		var errataTxs []stypes.ErrataTx
		for _, tx := range blockMeta.Transactions {
			if e.checkTransaction(tx.Hash) {
				e.logger.Debug().Msgf("block height: %d, tx: %s still exist", blockMeta.Height, tx.Hash)
				metaTxs = append(metaTxs, tx)
				continue
			}
			// this means the tx doesn't exist in chain ,thus should errata it
			errataTxs = append(errataTxs, stypes.ErrataTx{
				TxID:  common.TxID(tx.Hash),
				Chain: common.ETHChain,
			})
		}
		if len(errataTxs) > 0 {
			e.globalErrataQueue <- stypes.ErrataBlock{
				Height: blockMeta.Height,
				Txs:    errataTxs,
			}
		}
		// Let's get the block again to fix the block hash
		var block *etypes.Header
		block, err = e.getHeader(blockMeta.Height)
		if err != nil {
			e.logger.Err(err).Msgf("fail to get block verbose tx result: %d", blockMeta.Height)
			rescanBlockHeights = append(rescanBlockHeights, blockMeta.Height)
			continue
		}
		if !strings.EqualFold(blockMeta.BlockHash, block.Hash().Hex()) {
			// if the block hash is different as previously recorded , then the block should be rescanned
			rescanBlockHeights = append(rescanBlockHeights, blockMeta.Height)
		}
		blockMeta.PreviousHash = block.ParentHash.Hex()
		blockMeta.BlockHash = block.Hash().Hex()
		blockMeta.Transactions = metaTxs
		if err = e.blockMetaAccessor.SaveBlockMeta(blockMeta.Height, blockMeta); err != nil {
			e.logger.Err(err).Msgf("fail to save block meta of height: %d ", blockMeta.Height)
		}
	}
	return rescanBlockHeights, nil
}

func (e *ETHScanner) checkTransaction(hash string) bool {
	ctx, cancel := e.getContext()
	defer cancel()
	tx, pending, err := e.client.TransactionByHash(ctx, ecommon.HexToHash(hash))
	if err != nil || tx == nil {
		return false
	}

	// pending transactions may fail, but we should only errata when there is certainty
	if pending {
		e.logger.Warn().Msgf("tx: %s is in pending status", hash)
		return true // unknown, prefer false positive
	}

	// ensure the tx was successful
	receipt, err := e.getReceipt(hash)
	if err != nil {
		e.logger.Warn().Err(err).Msgf("fail to get receipt for tx: %s", hash)
		return true // unknown, prefer false positive
	}
	return receipt.Status == etypes.ReceiptStatusSuccessful
}

func (e *ETHScanner) getReceipt(hash string) (*etypes.Receipt, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.TransactionReceipt(ctx, ecommon.HexToHash(hash))
}

func (e *ETHScanner) getHeader(height int64) (*etypes.Header, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.HeaderByNumber(ctx, big.NewInt(height))
}

func (e *ETHScanner) getBlock(height int64) (*etypes.Block, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.BlockByNumber(ctx, big.NewInt(height))
}

func (e *ETHScanner) getRPCBlock(height int64) (*etypes.Block, error) {
	block, err := e.getBlock(height)
	if err == ethereum.NotFound {
		return nil, btypes.ErrUnavailableBlock
	}
	if err != nil {
		return nil, fmt.Errorf("fail to fetch block: %w", err)
	}
	return block, nil
}

func (e *ETHScanner) getFilterLogs(query ethereum.FilterQuery) ([]etypes.Log, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.FilterLogs(ctx, query)
}

func (e *ETHScanner) getRPCFilterLogs(query ethereum.FilterQuery) ([]etypes.Log, error) {
	ret, err := e.getFilterLogs(query)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch logs: %w", err)
	}
	return ret, nil
}

func (e *ETHScanner) getDecimals(token string) (uint64, error) {
	if IsETH(token) {
		return defaultDecimals, nil
	}
	to := ecommon.HexToAddress(token)
	input, err := e.erc20ABI.Pack(decimalMethod)
	if err != nil {
		return defaultDecimals, fmt.Errorf("fail to pack decimal method: %w", err)
	}
	ctx, cancel := e.getContext()
	defer cancel()
	res, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &to,
		Data: input,
	}, nil)
	if err != nil {
		return defaultDecimals, fmt.Errorf("fail to call smart contract get decimals: %w", err)
	}
	output, err := e.erc20ABI.Unpack(decimalMethod, res)
	if err != nil {
		return defaultDecimals, fmt.Errorf("fail to unpack decimal method call result: %w", err)
	}
	switch output[0].(type) {
	case uint8:
		decimals, ok := abi.ConvertType(output[0], new(uint8)).(*uint8)
		if !ok {
			return defaultDecimals, fmt.Errorf("dev error: fail to cast uint8")
		}
		return uint64(*decimals), nil
	case *big.Int:
		decimals, ok := abi.ConvertType(output[0], new(*big.Int)).(*big.Int)
		if !ok {
			return defaultDecimals, fmt.Errorf("dev error: fail to cast big.Int")
		}
		return decimals.Uint64(), nil
	}
	return defaultDecimals, fmt.Errorf("%s is %T fail to parse it", output[0], output[0])
}

// replace the . in symbol to *, and replace the - in symbol to #
// because . and - had been reserved to use in THORChain symbol
var symbolReplacer = strings.NewReplacer(".", "*", "-", "#", `\u0000`, "", "\u0000", "")

func sanitiseSymbol(symbol string) string {
	return symbolReplacer.Replace(symbol)
}

func (e *ETHScanner) getSymbol(token string) (string, error) {
	if IsETH(token) {
		return "ETH", nil
	}
	to := ecommon.HexToAddress(token)
	input, err := e.erc20ABI.Pack(symbolMethod)
	if err != nil {
		return "", nil
	}
	ctx, cancel := e.getContext()
	defer cancel()
	res, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &to,
		Data: input,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("fail to call to smart contract and get symbol: %w", err)
	}
	var symbol string
	output, err := e.erc20ABI.Unpack(symbolMethod, res)
	if err != nil {
		symbol = string(res)
		e.logger.Err(err).Msgf("fail to unpack symbol method call,token address: %s , symbol: %s", token, symbol)
		return sanitiseSymbol(symbol), nil
	}
	// nolint
	symbol = *abi.ConvertType(output[0], new(string)).(*string)
	return sanitiseSymbol(symbol), nil
}

// isToValidContractAddress this method make sure the transaction to address is to THORChain router or a whitelist address
func (e *ETHScanner) isToValidContractAddress(addr *ecommon.Address, includeWhiteList bool) bool {
	// get the smart contract used by thornode
	contractAddresses := e.pubkeyMgr.GetContracts(common.ETHChain)
	if includeWhiteList {
		contractAddresses = append(contractAddresses, whitelistSmartContractAddress...)
	}
	// combine the whitelist smart contract address
	for _, item := range contractAddresses {
		if strings.EqualFold(item.String(), addr.String()) {
			return true
		}
	}
	return false
}

func (e *ETHScanner) getTokenMeta(token string) (types.TokenMeta, error) {
	tokenMeta, err := e.tokens.GetTokenMeta(token)
	if err != nil {
		return types.TokenMeta{}, fmt.Errorf("fail to get token meta: %w", err)
	}
	if tokenMeta.IsEmpty() {
		isWhiteListToken := false
		for _, item := range e.whitelistTokens {
			if strings.EqualFold(item.Address, token) {
				isWhiteListToken = true
				break
			}
		}

		if !isWhiteListToken {
			return types.TokenMeta{}, fmt.Errorf("token: %s is not whitelisted", token)
		}
		var symbol string
		symbol, err = e.getSymbol(token)
		if err != nil {
			return types.TokenMeta{}, fmt.Errorf("fail to get symbol: %w", err)
		}
		var decimals uint64
		decimals, err = e.getDecimals(token)
		if err != nil {
			e.logger.Err(err).Msgf("fail to get decimals from smart contract, default to: %d", defaultDecimals)
		}
		e.logger.Info().Msgf("token:%s, decimals: %d", token, decimals)
		tokenMeta = types.NewTokenMeta(symbol, token, decimals)
		if err = e.tokens.SaveTokenMeta(symbol, token, decimals); err != nil {
			return types.TokenMeta{}, fmt.Errorf("fail to save token meta: %w", err)
		}
	}
	return tokenMeta, nil
}

// convertAmount will convert the amount to 1e8 , the decimals used by THORChain
func (e *ETHScanner) convertAmount(token string, amt *big.Int) cosmos.Uint {
	if IsETH(token) {
		return cosmos.NewUintFromBigInt(amt).QuoUint64(common.One * 100)
	}
	decimals := uint64(defaultDecimals)
	tokenMeta, err := e.getTokenMeta(token)
	if err != nil {
		e.logger.Err(err).Msgf("fail to get token meta for token address: %s", token)
	}
	if !tokenMeta.IsEmpty() {
		decimals = tokenMeta.Decimal
	}
	if decimals != defaultDecimals {
		var value big.Int
		amt = amt.Mul(amt, value.Exp(big.NewInt(10), big.NewInt(defaultDecimals), nil))
		amt = amt.Div(amt, value.Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	}
	return cosmos.NewUintFromBigInt(amt).QuoUint64(common.One * 100)
}

// return value 0 means use the default value which is common.THORChainDecimals, use 1e8 as precision
func (e *ETHScanner) getTokenDecimalsForTHORChain(token string) int64 {
	if IsETH(token) {
		return 0
	}
	tokenMeta, err := e.getTokenMeta(token)
	if err != nil {
		e.logger.Err(err).Msgf("fail to get token meta for token address: %s", token)
	}
	if tokenMeta.IsEmpty() {
		return 0
	}
	// when the token's precision is more than THORChain , that's fine , just use THORChainDecimals
	if tokenMeta.Decimal >= common.THORChainDecimals {
		return 0
	}
	return int64(tokenMeta.Decimal)
}

func (e *ETHScanner) getAssetFromTokenAddress(token string) (common.Asset, error) {
	if IsETH(token) {
		return common.ETHAsset, nil
	}
	tokenMeta, err := e.getTokenMeta(token)
	if err != nil {
		return common.EmptyAsset, fmt.Errorf("fail to get token meta: %w", err)
	}
	if tokenMeta.IsEmpty() {
		return common.EmptyAsset, fmt.Errorf("token metadata is empty")
	}
	return common.NewAsset(fmt.Sprintf("ETH.%s-%s", tokenMeta.Symbol, strings.ToUpper(tokenMeta.Address)))
}

// getTxInFromSmartContract returns txInItem
func (e *ETHScanner) getTxInFromSmartContract(ll *etypes.Log, receipt *etypes.Receipt, maxLogs int64) (*stypes.TxInItem, error) {
	e.logger.Debug().Msg("Parse tx from smart contract")
	txInItem := &stypes.TxInItem{
		Tx:     ll.TxHash.Hex()[2:],
		Height: big.NewInt(0).SetUint64(ll.BlockNumber),
	}
	cId, _ := e.cfg.ChainID.ChainID()
	txInItem.FromChain = cId

	// 1 is Transaction success state
	if receipt.Status != etypes.ReceiptStatusSuccessful {
		e.logger.Info().Msgf("Find a Tx(%s) state: %d means failed , ignore", ll.TxHash.String(), receipt.Status)
		return nil, nil
	}
	p := evm.NewSmartContractLogParser(e.gatewayABI)
	// txInItem will be changed in p.GetTxInItem function, so if the function return an error
	// txInItem should be abandoned
	if err := p.GetTxInItem(ll, txInItem); err != nil {
		return nil, fmt.Errorf("fail to parse logs, err: %w", err)
	}
	// under no circumstance ETH gas price will be less than 1 Gwei , unless it is in dev environment
	txGasPrice := receipt.EffectiveGasPrice

	e.logger.Info().Msgf("Find tx: %s, gas price: %s, gas used: %d, logIndex:%d",
		txInItem.Tx, txGasPrice.String(), receipt.GasUsed, ll.Index)

	e.logger.Debug().Msgf("Tx in item: %+v", txInItem)
	return txInItem, nil
}

func (e *ETHScanner) fromTxToTxInLog(ll *etypes.Log) (*stypes.TxInItem, error) {
	if ll == nil {
		return nil, nil
	}
	receipt, err := e.getReceipt(ll.TxHash.Hex())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("fail to get transaction receipt: %w", err)
	}
	if receipt.Status != etypes.ReceiptStatusSuccessful {
		if e.signerCacheManager != nil {
			e.signerCacheManager.RemoveSigned(ll.TxHash.String())
		}
		e.logger.Debug().Msgf("tx(%s) state: %d means failed , ignore", ll.TxHash.String(), receipt.Status)
		// todo will next 100
		return nil, nil
		//return e.getTxInFromFailedTransaction(tx, receipt), nil
	}

	ret, err := e.getTxInFromSmartContract(ll, receipt, 0)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
