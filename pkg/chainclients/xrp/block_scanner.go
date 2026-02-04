package xrp

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	xrplcommon "github.com/Peersyst/xrpl-go/xrpl/queries/common"
	"github.com/Peersyst/xrpl-go/xrpl/queries/ledger"
	requests "github.com/Peersyst/xrpl-go/xrpl/queries/transactions"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/blockscanner"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/utxo"
	mem "github.com/mapprotocol/compass-tss/x/memo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

// SolvencyReporter is to report solvency info to THORNode
type SolvencyReporter func(int64) error

const (
	// FeeUpdatePeriodBlocks is the block interval at which we report fee changes.
	FeeUpdatePeriodBlocks = 20

	// FeeCacheTransactions is the number of transactions over which we compute an average
	// (mean) fee price to use for outbound transactions. Note that only transactions
	// using the chain fee asset will be considered.
	FeeCacheTransactions = 200
)

var (
	ErrInvalidScanStorage = errors.New("scan storage is empty or nil")
	ErrInvalidMetrics     = errors.New("metrics is empty or nil")
	ErrEmptyTx            = errors.New("empty tx")
)

// XrpBlockScanner is to scan the blocks
type XrpBlockScanner struct {
	cfg              config.BifrostBlockScannerConfiguration
	logger           zerolog.Logger
	db               blockscanner.ScannerStorage
	bridge           shareTypes.Bridge
	solvencyReporter SolvencyReporter
	rpcClient        *rpc.Client

	globalNetworkFeeQueue chan types.NetworkFee

	// feeCache contains a rolling window of suggested fees.
	// Fees are stored at 100x the values on the observed chain due to compensate for the
	// difference in base chain decimals (relay:1e8, xrp:1e6).
	feeCache        []sdkmath.Uint
	lastFee         sdkmath.Uint
	lastAsgard      time.Time // cache vault address
	asgardAddresses []common.Address
}

// NewXrpBlockScanner create a new instance of BlockScan
func NewXrpBlockScanner(rpcHost string,
	cfg config.BifrostBlockScannerConfiguration,
	scanStorage blockscanner.ScannerStorage,
	bridge shareTypes.Bridge,
	m *metrics.Metrics,
	solvencyReporter SolvencyReporter,
) (*XrpBlockScanner, error) {
	if scanStorage == nil {
		return nil, errors.New("scanStorage is nil")
	}
	if m == nil {
		return nil, errors.New("metrics is nil")
	}

	logger := log.Logger.With().Str("module", "blockscanner").Str("chain", cfg.ChainID.String()).Logger()

	rpcConfig, err := rpc.NewClientConfig(rpcHost)
	if err != nil {
		return nil, fmt.Errorf("unable to create rpc config, %w", err)
	}
	rpcClient := rpc.NewClient(rpcConfig)
	return &XrpBlockScanner{
		cfg:              cfg,
		logger:           logger,
		db:               scanStorage,
		rpcClient:        rpcClient,
		feeCache:         make([]sdkmath.Uint, 0),
		lastFee:          sdkmath.NewUint(0),
		bridge:           bridge,
		solvencyReporter: solvencyReporter,
	}, nil
}

// GetHeight returns the index of the most recently validated ledger.
func (c *XrpBlockScanner) GetHeight() (int64, error) {
	ledgerIndex, err := c.rpcClient.GetLedgerIndex()
	if err != nil {
		return 0, err
	}

	return int64(ledgerIndex.Int()), nil
}

// FetchMemPool returns nothing since we are only concerned about finalized transactions in Xrp
func (c *XrpBlockScanner) FetchMemPool(height int64) (types.TxIn, error) {
	return types.TxIn{}, nil
}

// GetNetworkFee returns current chain network fee according to Bifrost.
func (c *XrpBlockScanner) GetNetworkFee() (transactionSize, transactionSwapSize, transactionFeeRate uint64) {
	return 1, 1, c.lastFee.Uint64()
}

func (c *XrpBlockScanner) updateFeeCache(fee sdkmath.Uint) {
	// sanity check to ensure fee is non-zero
	if fee.IsZero() {
		c.logger.Warn().Interface("fee", fee).Msg("transaction with zero fee")
		return
	}

	// add the fee to our cache
	c.feeCache = append(c.feeCache, fee)

	// truncate fee prices older than our max cached transactions
	if len(c.feeCache) > FeeCacheTransactions {
		c.feeCache = c.feeCache[(len(c.feeCache) - FeeCacheTransactions):]
	}
}

func (c *XrpBlockScanner) averageFee() sdkmath.Uint {
	// avoid divide by zero
	if len(c.feeCache) == 0 {
		return sdkmath.NewUint(0)
	}

	// compute mean
	sum := sdkmath.NewUint(0)
	for _, val := range c.feeCache {
		sum = sum.Add(val)
	}
	mean := sum.Quo(sdkmath.NewUint(uint64(len(c.feeCache))))

	return mean
}

func (c *XrpBlockScanner) updateFees(height int64) error {
	// post the gas fee over every cache period when we have a full gas cache
	if height%FeeUpdatePeriodBlocks == 0 && len(c.feeCache) == FeeCacheTransactions {
		avgFee := c.averageFee()

		// sanity check the fee is not zero
		if avgFee.IsZero() {
			return errors.New("suggested gas fee was zero")
		}

		// skip fee update if it has not changed
		if c.lastFee.Equal(avgFee) {
			return nil
		}

		// NOTE: We post the fee to the network instead of the transaction rate, and set the
		// transaction size 1 to ensure the MaxGas in the generated TxOut contains the
		// correct fee. We cannot pass the proper size and rate without a deeper change to
		// relay, as the rate on XRP chain is less than 1 and cannot be represented
		// by the uint.
		cId, _ := c.cfg.ChainID.ChainID()
		c.globalNetworkFeeQueue <- types.NetworkFee{
			ChainId:             cId,
			Height:              height,
			TransactionSize:     1,
			TransactionSwapSize: 1,
			TransactionRate:     avgFee.Uint64(),
		}

		c.lastFee = avgFee
		c.logger.Info().
			Uint64("fee", avgFee.Uint64()).
			Int64("height", height).
			Msg("sent network fee to relay")
	}

	return nil
}

func (c *XrpBlockScanner) processTxs(height int64, rawTxs []transaction.FlatTransaction) ([]*types.TxInItem, error) {
	var txIn []*types.TxInItem
	for _, rawTx := range rawTxs {
		// tx is nil, may not have been validated
		if rawTx == nil {
			continue
		}

		meta, err := c.decodeMetaBlobIfNecessary(rawTx)
		if err != nil {
			c.logger.Info().AnErr("error", err).Msg("fail to decode meta")
			continue
		}

		c.logger.Info().Any("meta", meta).Int64("height", height).
			Msg("get meta from tx")

		// Ignore failed transactions
		if meta["TransactionResult"] != "tesSUCCESS" {
			continue
		}

		ctxLog := c.logger.Info().Interface("tx", rawTx)
		ctxLog.Msg("get tx")
		flatTx, err := c.decodeTxBlobIfNecessary(rawTx)
		if err != nil {
			ctxLog.AnErr("error", err).Msg("fail to decode tx blob")
			continue
		}

		payment, err := c.processPayment(flatTx)
		if payment == nil && err == nil {
			// This was not a payment tx
			continue
		}
		if err != nil {
			ctxLog.AnErr("reason", err).Msg("skipping payment tx")
			continue
		}
		// update gas
		fee, err := convertFee(payment.Fee)
		if err != nil {
			return nil, fmt.Errorf("cannot convert xrp fee to relay fee: %w", err)
		}
		c.updateFeeCache(fee)
		// check if we are an asgard address
		if !c.isVaultAddress(payment.Destination.String()) && !c.isVaultAddress(payment.Account.String()) {
			continue
		}

		memo := ""
		if len(payment.Memos) == 1 {
			memo = payment.Memos[0].Memo.MemoData
		}
		// update the signer cache
		var parseMemo mem.Memo
		parseMemo, err = mem.ParseMemo(memo)
		if err != nil {
			// Debug log only as ParseMemo error is expected for THORName inbounds.
			ctxLog.Err(err).Msgf("fail to parse memo: %s", memo)
			continue
		}
		if !parseMemo.IsValid() {
			ctxLog.Str("memo", memo).Str("type", parseMemo.GetType().String()).Msg("invalid memo")
			continue
		}
		var (
			amount                                      int64
			topic                                       string
			invalidMemo                                 bool
			vaultAddress                                = payment.Account.String()
			callMethod                                  = constants.VoteTxOut
			toBytes, chainAndGasLimit, payload, toToken []byte
			orderId                                     ethcommon.Hash
			destChainID, gasUsed                        = big.NewInt(0), big.NewInt(0)
			txOutType                                   constants.TxInType
			mapChainID, _                               = common.MAPChain.ChainID()
			selfId, _                                   = c.cfg.ChainID.ChainID()
			nativeToken                                 = common.XRPAsset
		)
		hash, ok := rawTx["hash"].(string)
		if !ok {
			ctxLog.Msg("skipping tx, cannot cast hash to string")
			continue
		}
		toBytes, err = parseMemo.GetChain().DecodeAddress(parseMemo.GetDestination())
		if err != nil {
			ctxLog.Err(err).Str("txid", hash).Str("memo", memo).Msg("fail to decode memo")
			invalidMemo = true
		}

		// other2xrp
		if c.isVaultAddress(payment.Account.String()) { // sender
			switch parseMemo.GetType() {
			case mem.TxInbound:
				txOutType = constants.TRANSFER
			case mem.TxMigrate:
				txOutType = constants.MIGRATE
			case mem.TxRefund:
			default:
				ctxLog.Str("memo", memo).Str("type", parseMemo.GetType().String()).Msg("invalid memo")
				continue
			}
			chainAndGasLimit := make([]byte, 32)
			toChain := ethcommon.LeftPadBytes(selfId.Bytes(), 8)
			copy(chainAndGasLimit[8:16], toChain)
			orderId = ethcommon.HexToHash(parseMemo.GetOrderID())
			topic = constants.EventOfBridgeIn.GetTopic().String()
			gasUsed = big.NewInt(int64(payment.Fee.Uint64()))
		}
		// xrp2other
		if c.isVaultAddress(payment.Destination.String()) {
			// refund
			if invalidMemo {
				txOutType = constants.TRANSFER
				destChainID = mapChainID
			} else {
				switch parseMemo.GetType() {
				case mem.TxOutbound:
					destToken := parseMemo.GetToken()
					txOutType = constants.TRANSFER
					destChainID, err = c.bridge.GetChainID(parseMemo.GetChain().String())
					if err != nil {
						ctxLog.Msg(fmt.Sprint("fail to get destination chain id: %w, chain: %s", err, parseMemo.GetChain()))
						continue
					}
					payload, err = c.encodePayload(nativeToken.Symbol.String(),
						destToken, mapChainID, destChainID, toBytes, parseMemo)
					if err != nil {
						ctxLog.Msg(fmt.Sprint("fail to encode payload: %w", err))
						// todo refund
						continue
					}
				case mem.TxAdd:
					txOutType = constants.DEPOSIT
					destChainID = mapChainID
				default:
					ctxLog.Str("memo", memo).Str("type", parseMemo.GetType().String()).Msg("invalid memo")
					continue
				}
			}

			vaultAddress = payment.Destination.String()
			toToken, err = c.bridge.GetTokenAddress(selfId, nativeToken.Symbol.String())
			if err != nil {
				ctxLog.Msg(fmt.Sprint("fail to get token address: %w, chainID: %s, token: %s", err, selfId, nativeToken))
				continue
			}

			fromChain := ethcommon.LeftPadBytes(selfId.Bytes(), 8)
			toChain := ethcommon.LeftPadBytes(destChainID.Bytes(), 8)

			if destChainID.Cmp(mapChainID) != 0 {
				destChainID = mapChainID
				toChain = ethcommon.LeftPadBytes(mapChainID.Bytes(), 8)
				toBytes = c.bridge.GetFusionReceiver().Bytes()
			}
			copy(chainAndGasLimit[0:8], fromChain)
			copy(chainAndGasLimit[8:16], toChain)

			orderId = evm.GenerateOrderID(selfId.String(), hash)
			topic = constants.EventOfBridgeOut.GetTopic().String()
			callMethod = constants.VoteTxIn
		}

		vaultPbuKey, err := utxo.GetAsgardPubKeyByAddress(c.cfg.ChainID, c.bridge, common.Address(vaultAddress))
		if err != nil {
			ctxLog.Msg(fmt.Sprint("fail to get vault pub key by address: %w", err))
			continue
		}

		txIn = append(txIn, &types.TxInItem{
			Tx:               hash,                          // done
			Memo:             memo,                          // done
			Height:           new(big.Int).SetInt64(height), // done
			Amount:           big.NewInt(amount),
			OrderId:          orderId,     // done
			GasUsed:          gasUsed,     // done
			Token:            toToken,     // done
			Vault:            vaultPbuKey, // done
			From:             nil,
			To:               toBytes,                                 // done
			Payload:          payload,                                 // done
			Method:           callMethod,                              // done
			LogIndex:         0,                                       // done
			ChainAndGasLimit: new(big.Int).SetBytes(chainAndGasLimit), // done
			TxOutType:        uint8(txOutType),                        // done
			Sequence:         big.NewInt(0),                           // done
			Topic:            topic,                                   // done
			Timestamp:        time.Now().Unix(),                       // done
		})
	}

	return txIn, nil
}

// The expected response from the ledger method.
type LedgerResponseWithTxHashes struct {
	Ledger struct {
		Transactions []string `json:"transactions,omitempty"`
	} `json:"ledger"`
	LedgerHash  string                 `json:"ledger_hash"`
	LedgerIndex xrplcommon.LedgerIndex `json:"ledger_index"`
	Validated   bool                   `json:"validated,omitempty"`
}

func (c *XrpBlockScanner) FetchTxs(height, chainHeight int64) (types.TxIn, error) {
	// First retrieve all transaction hashes in block
	res, err := c.rpcClient.Request(&ledger.Request{
		LedgerIndex:  xrplcommon.LedgerIndex(height),
		Transactions: true,
		Expand:       false,
	})
	if err != nil {
		return types.TxIn{}, err
	}
	var ledgerTxHashes LedgerResponseWithTxHashes
	err = res.GetResult(&ledgerTxHashes)
	if err != nil {
		return types.TxIn{}, err
	}

	// Verify ledger has been validated, it should be
	if !ledgerTxHashes.Validated {
		return types.TxIn{}, btypes.ErrUnavailableBlock
	}

	// Next, get all transactions in block
	// Set binary to false, xrp client unfortunately does not fully support decoding all transactions
	ledger, err := c.rpcClient.GetLedger(&ledger.Request{
		LedgerIndex:  xrplcommon.LedgerIndex(height),
		Transactions: true,
		Expand:       true,
		Binary:       false,
	})
	if err != nil {
		return types.TxIn{}, err
	}

	// Verify ledger has been validated, it should be
	if !ledger.Validated {
		return types.TxIn{}, btypes.ErrUnavailableBlock
	}

	// Verify that we could fetch all transactions, if not, fetch remaining individually.
	flatTransactions := ledger.Ledger.Transactions
	if len(ledgerTxHashes.Ledger.Transactions) > len(flatTransactions) {
		c.logger.Info().Int("from ledger with just hashes", len(ledgerTxHashes.Ledger.Transactions)).
			Int("from ledger with expanded txs", len(flatTransactions)).Msg("number of transactions don't match")
		// We did not get all transactions, most likely due to the 256 tx limit when requesting tx in json.
		// Average txs/ledger is currently well below 256.
		// In the future, with better xrp client support, we could request an unlimited amount of txs when binary=true,
		// but this code should still remain to address any other potential server side limitation,
		// i.e.: max response size setting, new tx amount limitation for binary==true, remote/centralized server.
		// Retrieve remaining txs individually, enough load could get us behind, but better behind than missing txs.
		flatTransactions, err = c.fetchTxsByHash(ledgerTxHashes.Ledger.Transactions, flatTransactions)
		if err != nil {
			return types.TxIn{}, err
		}
	}

	txs, err := c.processTxs(height, flatTransactions)
	if err != nil {
		return types.TxIn{}, err
	}

	txIn := types.TxIn{
		Chain:    c.cfg.ChainID,
		TxArray:  txs,
		Filtered: false,
		MemPool:  false,
	}

	// skip reporting network fee and solvency if block more than flexibility blocks from tip
	if chainHeight-height > c.cfg.ObservationFlexibilityBlocks {
		return txIn, nil
	}

	err = c.updateFees(height)
	if err != nil {
		c.logger.Err(err).Int64("height", height).Msg("unable to update network fee")
	}

	if err = c.solvencyReporter(height); err != nil {
		c.logger.Err(err).Msg("fail to send solvency to relay")
	}

	return txIn, nil
}

const (
	maxConcurrentRequests = 10
)

func (c *XrpBlockScanner) fetchTxsByHash(txHashes []string, txs []transaction.FlatTransaction) ([]transaction.FlatTransaction, error) {
	// Create txMap for easy lookup of known txs
	txMap := make(map[string]transaction.FlatTransaction)

	// Populate with transactions that we have and don't want to query again
	for _, tx := range txs {
		// We verified we got this tx from a validated ledger, set per tx
		tx["validated"] = true
		hash, ok := tx["hash"].(string)
		if !ok {
			continue
		}
		txMap[hash] = tx
	}

	ctx := context.Background()
	transactions := make([]transaction.FlatTransaction, len(txHashes))
	sem := semaphore.NewWeighted(int64(maxConcurrentRequests))
	errChan := make(chan error, maxConcurrentRequests)
	var wg sync.WaitGroup

	for i, hash := range txHashes {
		// Check for any errors before starting new goroutines
		select {
		case err := <-errChan:
			return nil, err
		default:
			// No errors yet, continue
		}

		// Check whether we already have this tx, if so, add it to maintain order and don't request it
		if txMap[hash] != nil {
			transactions[i] = txMap[hash]
			continue
		}

		// Acquire semaphore (blocks if maxConcurrentRequests already running)
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(idx int, txHash string) {
			defer sem.Release(1)
			defer wg.Done()

			// Request the transaction from the server
			res, err := c.rpcClient.Request(&requests.TxRequest{
				Transaction: txHash,
				Binary:      false,
			})
			if err != nil {
				c.logger.Err(err).Str("hash", txHash).Msg("error requesting tx")
				select {
				case errChan <- fmt.Errorf("failed to request transaction %s: %w", txHash, err):
				default:
					// Error channel is full, another error is already being processed
				}
				return
			}

			var txResponse transaction.FlatTransaction
			err = res.GetResult(&txResponse)
			if err != nil {
				c.logger.Err(err).Str("hash", txHash).Msg("error parsing tx results")
				select {
				case errChan <- fmt.Errorf("failed to parse tx results %s: %w", txHash, err):
				default:
					// Error channel is full, another error is already being processed
				}
				return
			}

			transactions[idx] = txResponse
		}(i, hash)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if any errors occurred (non-blocking)
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	// Check if context was canceled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return transactions, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Address Checks
////////////////////////////////////////////////////////////////////////////////////////

func (c *XrpBlockScanner) isVaultAddress(addressToCheck string) bool {
	asgards, err := c.getVaultAddress()
	if err != nil {
		c.logger.Err(err).Msg("fail to get vault addresses")
		return false
	}
	for _, addr := range asgards {
		if strings.EqualFold(addr.String(), addressToCheck) {
			return true
		}
	}
	return false
}

func (c *XrpBlockScanner) getVaultAddress() ([]common.Address, error) {
	if time.Since(c.lastAsgard) < constants.MAPRelayChainBlockTime && c.asgardAddresses != nil {
		return c.asgardAddresses, nil
	}
	newAddresses, err := utxo.GetAsgardAddress(c.cfg.ChainID, c.bridge)
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards: %w", err)
	}
	if len(newAddresses) > 0 { // ensure we don't overwrite with empty list
		c.asgardAddresses = newAddresses
	}
	c.lastAsgard = time.Now()
	return c.asgardAddresses, nil
}

func (c *XrpBlockScanner) encodePayload(nativeToken, destToken string, mapChainID, destChainID *big.Int, to []byte, parsedMemo mem.Memo) (payload []byte, err error) {
	var relayData []byte
	// if the dest token is not native token, we need to build the relay data
	if !strings.EqualFold(destToken, nativeToken) {
		destTokenAddress, err := c.bridge.GetTokenAddress(mapChainID, destToken)
		if err != nil {
			return nil, fmt.Errorf("fail to get token address: %w, chainID: %s, token: %s", err, mapChainID, destToken)
		}
		decimals, err := c.bridge.GetTokenDecimals(mapChainID, destTokenAddress)
		if err != nil {
			return nil, fmt.Errorf("fail to get token decimals: %w, chainID: %s, token: %s", err, mapChainID, destToken)
		}
		if decimals.Cmp(big.NewInt(0)) == 0 {
			decimals = big.NewInt(constants.DefaultTokenDecimals)
		}
		minAmount := utxo.ConvertDecimal(parsedMemo.GetAmount().BigInt(), 6, decimals.Uint64())
		relayData, err = utxo.EncodeRelayData(ethcommon.BytesToAddress(destTokenAddress), minAmount)
		if err != nil {
			return nil, fmt.Errorf("fail to encode relay data: %w, token: %s, minAmount: %d", err, ethcommon.BytesToAddress(destTokenAddress), minAmount)
		}
	}

	var targetData []byte
	// if the dest chain is not the map chain, we need to build the target data.
	if destChainID.Cmp(mapChainID) != 0 {
		targetData, err = utxo.EncodeTargetData(to, destChainID)
		if err != nil {
			return nil, fmt.Errorf("fail to encode target data: %w, to: %s, chainID: %d", err, parsedMemo.GetDestination(), destChainID)
		}
	}

	affiliateData, err := c.encodeAffiliates(parsedMemo.GetAffiliates())
	if err != nil {
		return nil, fmt.Errorf("fail to encode affiliates: %w", err)
	}

	payload, err = utxo.EncodePayload(affiliateData, relayData, targetData)
	if err != nil {
		return nil, fmt.Errorf("fail to encode payload: %w", err)
	}
	return payload, nil
}

func (c *XrpBlockScanner) encodeAffiliates(affiliates mem.Affiliates) ([]byte, error) {
	as := make([]*utxo.Affiliate, 0, len(affiliates))
	for _, aff := range affiliates {
		if aff.Compressed {
			id, err := c.bridge.GetAffiliateIDByAlias(aff.Name)
			if err != nil {
				return []byte{}, fmt.Errorf("fail to get affiliate id by alias: %w", err)
			}
			as = append(as, &utxo.Affiliate{
				ID:  id,
				Bps: uint16(aff.Bps.Uint64()),
			})
			continue
		}

		id, err := c.bridge.GetAffiliateIDByName(aff.Name)
		if err != nil {
			return []byte{}, fmt.Errorf("fail to get affiliate id by name: %w", err)
		}

		as = append(as, &utxo.Affiliate{
			ID:  id,
			Bps: uint16(aff.Bps.Uint64()),
		})
	}
	return utxo.EncodeAffiliateData(as)
}
