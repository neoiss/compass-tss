package mapo

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/mapprotocol/compass-tss/internal/structure"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/blockscanner"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/mapprotocol/compass-tss/config"

	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
)

type MapChainBlockScan struct {
	logger            zerolog.Logger
	wg                *sync.WaitGroup
	stopChan          chan struct{}
	txOutChan         chan types.TxOut
	keygenChan        chan *structure.KeyGen
	cfg               config.BifrostBlockScannerConfiguration
	scannerStorage    blockscanner.ScannerStorage
	mapBridge         shareTypes.Bridge
	errCounter        *prometheus.CounterVec
	pubkeyMgr         pubkeymanager.PubKeyValidator
	blockMetaAccessor evm.BlockMetaAccessor
}

// NewBlockScan create a new instance of map block scanner
func NewBlockScan(cfg config.BifrostBlockScannerConfiguration, scanStorage blockscanner.ScannerStorage, bridge shareTypes.Bridge,
	m *metrics.Metrics, pubkeyMgr pubkeymanager.PubKeyValidator) (*MapChainBlockScan, error) {
	if scanStorage == nil {
		return nil, errors.New("scanStorage is nil")
	}
	if m == nil {
		return nil, errors.New("metric is nil")
	}
	// set storage prefixes
	prefixBlockMeta := fmt.Sprintf("%s-blockmeta-", strings.ToLower(cfg.ChainID.String()))
	prefixSignedMeta := fmt.Sprintf("%s-signedtx-", strings.ToLower(cfg.ChainID.String()))
	//prefixTokenMeta := fmt.Sprintf("%s-tokenmeta-", strings.ToLower(cfg.ChainID.String()))

	// create block meta accessor
	blockMetaAccessor, err := evm.NewLevelDBBlockMetaAccessor(
		prefixBlockMeta, prefixSignedMeta, scanStorage.GetInternalDb(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create block meta accessor: %w", err)
	}

	return &MapChainBlockScan{
		logger:            log.With().Str("module", "blockscanner").Str("chain", "map").Logger(),
		wg:                &sync.WaitGroup{},
		stopChan:          make(chan struct{}),
		txOutChan:         make(chan types.TxOut),
		keygenChan:        make(chan *structure.KeyGen),
		cfg:               cfg,
		scannerStorage:    scanStorage,
		mapBridge:         bridge,
		errCounter:        m.GetCounterVec(metrics.MapChainBlockScannerError),
		pubkeyMgr:         pubkeyMgr,
		blockMetaAccessor: blockMetaAccessor,
	}, nil
}

func (b *MapChainBlockScan) GetTxOutMessages() <-chan types.TxOut {
	return b.txOutChan
}

func (b *MapChainBlockScan) GetKeygenMessages() <-chan *structure.KeyGen {
	return b.keygenChan
}

func (b *MapChainBlockScan) GetHeight() (int64, error) {
	return b.mapBridge.GetBlockHeight()
}

// GetNetworkFee : MapChainBlockScan's only exists to satisfy the Fetcher interface
// and should never be called, since broadcast network fees are for external chains' observed fees.
func (b *MapChainBlockScan) GetNetworkFee() (transactionSize, transactionSwapSize, transactionFeeRate uint64) {
	b.logger.Error().Msg("MapChainBlockScan GetNetworkFee was called (which should never happen)")
	return 0, 0, 0
}

func (b *MapChainBlockScan) FetchMemPool(height int64) (types.TxIn, error) {
	return types.TxIn{}, nil
}

func (b *MapChainBlockScan) FetchTxs(height, _ int64) (types.TxIn, error) {
	if err := b.processTxOutBlock(height); err != nil {
		return types.TxIn{}, err
	}
	if err := b.processKeygenBlock(); err != nil {
		return types.TxIn{}, err
	}
	return types.TxIn{}, nil
}

func (b *MapChainBlockScan) processKeygenBlock() error {
	// done
	keygen, err := b.mapBridge.GetKeygenBlock()
	if err != nil {
		return fmt.Errorf("fail to get keygen from mapBridge: %w", err)
	}
	if keygen == nil {
		return nil
	}

	b.keygenChan <- keygen
	return nil
}

// todo handler
func (b *MapChainBlockScan) processTxOutBlock(blockHeight int64) error {
	tx, err := b.mapBridge.GetTxByBlockNumber(blockHeight, b.cfg.Mos)
	if err != nil {
		if errors.Is(err, btypes.ErrUnavailableBlock) {
			// custom error (to be dropped and not logged) because the block is
			// available yet
			return btypes.ErrUnavailableBlock
		}
		return fmt.Errorf("fail to get keysign from block scanner: %w", err)
	}

	if len(tx.TxArray) == 0 {
		b.logger.Debug().Int64("block", blockHeight).Msg("Nothing to process")
		return nil
	}
	b.txOutChan <- tx
	return nil
}
