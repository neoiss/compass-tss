package signer

import (
	"errors"
	"fmt"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/blockscanner"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/mapprotocol/compass-tss/config"

	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	ttypes "github.com/mapprotocol/compass-tss/x/types"
)

type MapChainBlockScan struct {
	logger         zerolog.Logger
	wg             *sync.WaitGroup
	stopChan       chan struct{}
	txOutChan      chan types.TxOut
	keygenChan     chan ttypes.KeygenBlock
	cfg            config.BifrostBlockScannerConfiguration
	scannerStorage blockscanner.ScannerStorage
	thorchain      mapo.ThorchainBridge
	errCounter     *prometheus.CounterVec
	pubkeyMgr      pubkeymanager.PubKeyValidator
}

// NewMapChainBlockScan create a new instance of map block scanner
func NewMapChainBlockScan(cfg config.BifrostBlockScannerConfiguration, scanStorage blockscanner.ScannerStorage, bridge mapo.ThorchainBridge, m *metrics.Metrics, pubkeyMgr pubkeymanager.PubKeyValidator) (*MapChainBlockScan, error) {
	if scanStorage == nil {
		return nil, errors.New("scanStorage is nil")
	}
	if m == nil {
		return nil, errors.New("metric is nil")
	}
	return &MapChainBlockScan{
		logger:         log.With().Str("module", "blockscanner").Str("chain", "map").Logger(),
		wg:             &sync.WaitGroup{},
		stopChan:       make(chan struct{}),
		txOutChan:      make(chan types.TxOut),
		keygenChan:     make(chan ttypes.KeygenBlock),
		cfg:            cfg,
		scannerStorage: scanStorage,
		thorchain:      bridge,
		errCounter:     m.GetCounterVec(metrics.MapChainBlockScannerError),
		pubkeyMgr:      pubkeyMgr,
	}, nil
}

func (b *MapChainBlockScan) GetTxOutMessages() <-chan types.TxOut {
	return b.txOutChan
}

func (b *MapChainBlockScan) GetKeygenMessages() <-chan ttypes.KeygenBlock {
	return b.keygenChan
}

func (b *MapChainBlockScan) GetHeight() (int64, error) {
	return b.thorchain.GetBlockHeight()
}

// GetNetworkFee : MapChainBlockScan's only exists to satisfy the BlockScannerFetcher interface
// and should never be called, since broadcast network fees are for external chains' observed fees.
func (b *MapChainBlockScan) GetNetworkFee() (transactionSize, transactionFeeRate uint64) {
	b.logger.Error().Msg("MapChainBlockScan GetNetworkFee was called (which should never happen)")
	return 0, 0
}

func (b *MapChainBlockScan) FetchMemPool(height int64) (types.TxIn, error) {
	return types.TxIn{}, nil
}

func (b *MapChainBlockScan) FetchTxs(height, _ int64) (types.TxIn, error) {
	if err := b.processTxOutBlock(height); err != nil {
		return types.TxIn{}, err
	}
	if err := b.processKeygenBlock(height); err != nil {
		return types.TxIn{}, err
	}
	return types.TxIn{}, nil
}

func (b *MapChainBlockScan) processKeygenBlock(blockHeight int64) error {
	pk := b.pubkeyMgr.GetNodePubKey()
	keygen, err := b.thorchain.GetKeygenBlock(blockHeight, pk.String())
	if err != nil {
		return fmt.Errorf("fail to get keygen from thorchain: %w", err)
	}

	// custom error (to be dropped and not logged) because the block is
	// available yet
	if keygen.Height == 0 {
		return btypes.ErrUnavailableBlock
	}

	if len(keygen.Keygens) > 0 {
		b.keygenChan <- keygen
	}
	return nil
}

func (b *MapChainBlockScan) processTxOutBlock(blockHeight int64) error {
	for _, pk := range b.pubkeyMgr.GetSignPubKeys() {
		if len(pk.String()) == 0 {
			continue
		}
		tx, err := b.thorchain.GetKeysign(blockHeight, pk.String())
		if err != nil {
			if errors.Is(err, btypes.ErrUnavailableBlock) {
				// custom error (to be dropped and not logged) because the block is
				// available yet
				return btypes.ErrUnavailableBlock
			}
			return fmt.Errorf("fail to get keysign from block scanner: %w", err)
		}

		if len(tx.TxArray) == 0 {
			b.logger.Debug().Int64("block", blockHeight).Msg("nothing to process")
			continue
		}
		b.txOutChan <- tx
	}
	return nil
}
