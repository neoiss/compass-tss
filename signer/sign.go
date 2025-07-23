package signer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec"
	ecommon "github.com/ethereum/go-ethereum/common"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/internal/structure"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/observer"
	"github.com/mapprotocol/compass-tss/pkg/chainclients"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/utxo"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/mapprotocol/compass-tss/tss"
	tssp "github.com/mapprotocol/compass-tss/tss/go-tss/tss"
	ttypes "github.com/mapprotocol/compass-tss/x/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Signer will pull the tx out from thorchain and then forward it to chain
type Signer struct {
	logger               zerolog.Logger
	cfg                  config.Bifrost
	wg                   *sync.WaitGroup
	thorchainBridge      shareTypes.Bridge
	stopChan             chan struct{}
	blockScanner         *blockscanner.BlockScanner
	mapChainBlockScanner *mapo.MapChainBlockScan
	chains               map[common.Chain]chainclients.ChainClient
	storage              SignerStorage
	m                    *metrics.Metrics
	errCounter           *prometheus.CounterVec
	tssKeygen            *tss.KeyGen
	tssServer            *tssp.TssServer
	pubkeyMgr            pubkeymanager.PubKeyValidator
	constantsProvider    *ConstantsProvider
	localPubKey          common.PubKey
	tssKeysignMetricMgr  *metrics.TssKeysignMetricMgr
	observer             *observer.Observer
	pipeline             *pipeline
	isKeyGen             bool
}

// NewSigner create a new instance of signer
func NewSigner(cfg config.Bifrost,
	bridge shareTypes.Bridge,
	thorKeys *keys.Keys,
	pubkeyMgr pubkeymanager.PubKeyValidator,
	tssServer *tssp.TssServer,
	chains map[common.Chain]chainclients.ChainClient,
	m *metrics.Metrics,
	tssKeysignMetricMgr *metrics.TssKeysignMetricMgr,
	obs *observer.Observer,
) (*Signer, error) {
	storage, err := NewSignerStore(cfg.Signer.SignerDbPath, cfg.Signer.LevelDB, bridge.GetConfig().SignerPasswd)
	if err != nil {
		return nil, fmt.Errorf("fail to create thorchain scan storage: %w", err)
	}
	if tssKeysignMetricMgr == nil {
		return nil, fmt.Errorf("fail to create signer , tss keysign metric manager is nil")
	}
	var na *structure.MaintainerInfo
	for i := 0; i < 300; i++ { // wait for 5 min before timing out
		signerAddr, err := thorKeys.GetEthAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get address from thorKeys signer: %w", err)
		}
		na, err = bridge.GetNodeAccount(signerAddr.String())
		if err != nil {
			return nil, fmt.Errorf("fail to get node account from thorchain,err:%w", err)
		}
		if na == nil {
			continue
		}

		if len(na.Secp256Pubkey) != 0 {
			break
		}
		time.Sleep(constants.MAPRelayChainBlockTime)
		log.Info().Msg("Waiting for node account to be registered...")
	}

	if na == nil {
		return nil, fmt.Errorf("fail to get node account from map relay chain")
	}
	if len(na.Secp256Pubkey) == 0 {
		return nil, fmt.Errorf("unable to find pubkey for this node account. exiting... ")
	}
	selfKey := common.PubKey(ecommon.Bytes2Hex(na.Secp256Pubkey))
	pubkeyMgr.AddNodePubKey(selfKey)

	cfg.Signer.BlockScanner.ChainID = common.MAPChain // hard code to map
	// Create pubkey manager and add our private key
	mapChainBlockScanner, err := mapo.NewBlockScan(cfg.Signer.BlockScanner, storage, bridge, m, pubkeyMgr)
	if err != nil {
		return nil, fmt.Errorf("fail to create thorchain block scan: %w", err)
	}

	blockScanner, err := blockscanner.NewBlockScanner(cfg.Signer.BlockScanner, storage, m, bridge, mapChainBlockScanner)
	if err != nil {
		return nil, fmt.Errorf("fail to create block scanner: %w", err)
	}
	ops := []shareTypes.BridgeOption{mapo.WithBlockScanner(mapChainBlockScanner)}
	err = bridge.InitBlockScanner(ops...)
	if err != nil {
		return nil, fmt.Errorf("fail to init block scanner: %w", err)
	}

	kg, err := tss.NewTssKeyGen(thorKeys, tssServer, bridge)
	if err != nil {
		return nil, fmt.Errorf("fail to create Tss Key gen,err:%w", err)
	}
	constantProvider := NewConstantsProvider(bridge)
	return &Signer{
		logger:               log.With().Str("module", "signer").Logger(),
		cfg:                  cfg,
		wg:                   &sync.WaitGroup{},
		stopChan:             make(chan struct{}),
		blockScanner:         blockScanner,
		mapChainBlockScanner: mapChainBlockScanner,
		chains:               chains,
		m:                    m,
		storage:              storage,
		errCounter:           m.GetCounterVec(metrics.SignerError),
		pubkeyMgr:            pubkeyMgr,
		thorchainBridge:      bridge,
		tssKeygen:            kg,
		tssServer:            tssServer,
		constantsProvider:    constantProvider,
		localPubKey:          selfKey,
		tssKeysignMetricMgr:  tssKeysignMetricMgr,
		observer:             obs,
		isKeyGen:             false,
	}, nil
}

func (s *Signer) getChain(chainID common.Chain) (chainclients.ChainClient, error) {
	chain, ok := s.chains[chainID]
	if !ok {
		s.logger.Debug().Str("chain", chainID.String()).Msg("is not supported yet")
		return nil, errors.New("not supported")
	}
	return chain, nil
}

// Start signer process
func (s *Signer) Start() error {
	//  todo handler annotate
	//s.wg.Add(1)
	//go s.processTxnOut(s.mapChainBlockScanner.GetTxOutMessages(), 1) // cache local

	s.wg.Add(1)
	go s.processKeygen(s.mapChainBlockScanner.GetKeygenMessages())

	//s.wg.Add(1)
	//go s.signTransactions()

	s.blockScanner.Start(nil, nil)
	return nil
}

func (s *Signer) shouldSign(tx types.TxOutItem) bool {
	return s.pubkeyMgr.HasPubKey(tx.VaultPubKey)
}

// signTransactions - looks for work to do by getting a list of all unsigned
// transactions stored in the storage
func (s *Signer) signTransactions() {
	s.logger.Info().Msg("start to sign transactions")
	defer s.logger.Info().Msg("stop to sign transactions")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		default:
			// When map relay chain is catching up , bifrost might get stale data from compass-tss , thus it shall pause signing
			catchingUp, err := s.thorchainBridge.IsSyncing()
			if err != nil {
				s.logger.Error().Err(err).Msg("fail to get thorchain sync status")
				time.Sleep(constants.MAPRelayChainBlockTime)
				break // this will break select
			}
			if !catchingUp {
				s.processTransactions()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func runWithContext(ctx context.Context, fn func() ([]byte, *types.TxInItem, error)) ([]byte, *types.TxInItem, error) {
	ch := make(chan error, 1)
	var checkpoint []byte
	var txIn *types.TxInItem
	go func() {
		var err error
		checkpoint, txIn, err = fn()
		ch <- err
	}()
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case err := <-ch:
		return checkpoint, txIn, err
	}
}

func (s *Signer) processTransactions() {
	signerConcurrency, err := s.thorchainBridge.GetMimir(constants.SignerConcurrency.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get signer concurrency mimir")
		return
	}

	// default to 10 if unset
	if signerConcurrency <= 0 {
		signerConcurrency = 10
	}

	// if previously set to different concurrency, drain existing signings
	if s.pipeline != nil && s.pipeline.concurrency != signerConcurrency {
		s.pipeline.Wait()
		s.pipeline = nil
	}

	// if not set, or set to different concurrency, create new pipeline
	if s.pipeline == nil {
		s.pipeline, err = newPipeline(signerConcurrency)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to create new pipeline")
			return
		}
	}

	// process transactions
	s.pipeline.SpawnSignings(s, s.thorchainBridge)
}

// processTxnOut processes outbound TxOuts and save them to storage
func (s *Signer) processTxnOut(ch <-chan types.TxOut, idx int) {
	s.logger.Info().Int("idx", idx).Msg("start to process tx out")
	defer s.logger.Info().Int("idx", idx).Msg("stop to process tx out")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		case txOut, more := <-ch:
			if !more {
				return
			}
			s.logger.Info().Msgf("Received a TxOut Array of %v from the Thorchain", txOut)
			items := make([]TxOutStoreItem, 0, len(txOut.TxArray))

			for i, tx := range txOut.TxArray {
				items = append(items, NewTxOutStoreItem(txOut.Height, tx.TxOutItem(txOut.Height), int64(i)))
			}
			if err := s.storage.Batch(items); err != nil {
				s.logger.Error().Err(err).Msg("fail to save tx out items to storage")
			}
		}
	}
}

func (s *Signer) processKeygen(ch <-chan *structure.KeyGen) {
	s.logger.Info().Msg("start to process keygen")
	defer s.logger.Info().Msg("stop to process keygen")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		case keygenBlock, more := <-ch:
			if !more {
				return
			}
			if s.isKeyGen {
				fmt.Println("ignore keyGen msg, because it is already a keygen, epoch=", keygenBlock.Epoch)
				continue
			}
			s.isKeyGen = true
			s.logger.Info().Interface("keygenBlock", keygenBlock).Msg("received a keygen block from map relay")
			s.processKeygenBlock(keygenBlock)
			s.isKeyGen = false
		}
	}
}

func (s *Signer) scheduleKeygenRetry(keygenBlock *structure.KeyGen) bool {
	//// todo handler
	//churnRetryInterval, err := s.thorchainBridge.GetMimir(constants.ChurnRetryInterval.String())
	//if err != nil {
	//	s.logger.Error().Err(err).Msg("fail to get churn retry mimir")
	//	return false
	//}
	//if churnRetryInterval <= 0 {
	//	churnRetryInterval = constants.NewConstantValue().GetInt64Value(constants.ChurnRetryInterval)
	//}
	//keygenRetryInterval, err := s.thorchainBridge.GetMimir(constants.KeygenRetryInterval.String())
	//if err != nil {
	//	s.logger.Error().Err(err).Msg("fail to get keygen retries mimir")
	//	return false
	//}
	//if keygenRetryInterval <= 0 {
	//	return false
	//}
	//
	//// sanity check the retry interval is at least 1.5x the timeout
	//retryIntervalDuration := time.Duration(keygenRetryInterval) * constants.MAPRelayChainBlockTime
	//if retryIntervalDuration <= s.cfg.Signer.KeygenTimeout*3/2 {
	//	s.logger.Error().
	//		Stringer("retryInterval", retryIntervalDuration).
	//		Stringer("keygenTimeout", s.cfg.Signer.KeygenTimeout).
	//		Msg("retry interval too short")
	//	return false
	//}
	//

	////height, err := s.thorchainBridge.GetBlockHeight()
	////if err != nil {
	////	s.logger.Error().Err(err).Msg("fail to get last chain height")
	////	return false
	////}
	////
	////// target retry height is the next keygen retry interval over the keygen block height
	////targetRetryHeight := (keygenRetryInterval - ((height - keygenBlock.Height) % keygenRetryInterval)) + height
	////
	////// skip trying close to churn retry
	////if targetRetryHeight > keygenBlock.Height+churnRetryInterval-keygenRetryInterval {
	////	return false
	////}
	//
	//go func() {
	//	// every block, try to start processing again
	//	for {
	//		time.Sleep(constants.MAPRelayChainBlockTime)
	//		// trunk-ignore(golangci-lint/govet): shadow
	//		height, err := s.thorchainBridge.GetBlockHeight()
	//		if err != nil {
	//			s.logger.Error().Err(err).Msg("fail to get last chain height")
	//		}
	//		if height >= targetRetryHeight {
	//			s.logger.Info().
	//				Interface("keygenBlock", keygenBlock).
	//				Int64("currentHeight", height).
	//				Msg("retrying keygen")
	//			s.processKeygenBlock(keygenBlock)
	//			return
	//		}
	//	}
	//}()
	//
	//s.logger.Info().
	//	Interface("keygenBlock", keygenBlock).
	//	Int64("retryHeight", targetRetryHeight).
	//	Msg("scheduled keygen retry")

	return true
}

func (s *Signer) processKeygenBlock(keygenBlock *structure.KeyGen) {
	s.logger.Debug().Interface("keygenBlock", keygenBlock).Msg("processing keygen block")
	members := make(common.PubKeys, 0, len(keygenBlock.Ms))
	memberAddrs := make([]ecommon.Address, 0, len(keygenBlock.Ms))
	for _, ele := range keygenBlock.Ms {
		if ele.Addr.String() == "" {
			continue
		}
		members = append(members, common.PubKey(ecommon.Bytes2Hex(ele.Secp256Pubkey)))
		memberAddrs = append(memberAddrs, ele.Addr)
	}
	// NOTE: in practice there is only one keygen in the keygen block
	//for _, keygenReq := range keygenBlock.Keygens {
	keygenStart := time.Now()
	// todo debug blame
	pubKey, blame, err := s.tssKeygen.GenerateNewKey(keygenBlock.Epoch.Int64(), members)
	if !blame.IsEmpty() {
		s.logger.Error().
			Str("reason", blame.FailReason).
			Interface("nodes", blame.BlameNodes).
			Msg("keygen blame")
	}
	keygenTime := time.Since(keygenStart).Milliseconds()
	if err != nil {
		s.errCounter.WithLabelValues("fail_to_keygen_pubkey", "").Inc()
		s.logger.Error().Err(err).Msg("fail to generate new pubkey")
		return
	}

	s.logger.Info().Int64("keygenTime", keygenTime).Msg("processKeygenBlock keyGen time")
	// generate a verification signature to ensure we can sign with the new key
	if len(pubKey.Secp256k1.String()) == 0 {
		fmt.Println("pk is empty ------------------ ")
		return
	}
	secp256k1Sig := s.secp256k1VerificationSignature(pubKey.Secp256k1)
	if len(secp256k1Sig) == 0 {
		fmt.Println("secp256k1Sig ------------------ ", len(secp256k1Sig))
		time.Sleep(time.Minute)
		return
	}
	if err = s.sendKeygenToMap(keygenBlock.Epoch, pubKey.Secp256k1, nil, memberAddrs, secp256k1Sig); err != nil { // todo handler blame
		s.errCounter.WithLabelValues("fail_to_broadcast_keygen", "").Inc()
		s.logger.Error().Err(err).Msg("fail to broadcast keygen")
	}

	fmt.Println("processKeygenBlock  GenerateNewKey 111111 -------------- ")
	// monitor the new pubkey and any new members
	if !pubKey.Secp256k1.IsEmpty() {
		s.pubkeyMgr.AddPubKey(pubKey.Secp256k1, true)
	}
	for _, pk := range members {
		s.pubkeyMgr.AddPubKey(pk, false)
	}
}

// secp256k1VerificationSignature will make a best effort to sign the public key with
// its own private key as a sanity check to ensure parties are able to sign. The
// signature will be included in the TssPool message if successful, and verified by
// THORNode before the keygen is accepted.
func (s *Signer) secp256k1VerificationSignature(pk common.PubKey) []byte {
	// create keysign instance
	ks, err := tss.NewKeySign(s.tssServer, s.thorchainBridge)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to create keySign for secp256k1 check signing")
		return nil
	}
	ks.Start()
	defer ks.Stop()

	// sign the public key with its own private key
	ethPubKey, err := ecrypto.DecompressPubkey(ecommon.Hex2Bytes(pk.String()))
	if err != nil {
		return nil
	}
	pubBytes := ecrypto.FromECDSAPub(ethPubKey)

	data := pubBytes[1:]
	dataHash := ecrypto.Keccak256(data)
	sigBytes, v, err := ks.RemoteSign(dataHash[:], pk.String())
	fmt.Println("v ------------------ ", v)
	if err != nil {
		// this is expected in some cases if we were not in the signing party
		s.logger.Info().Err(err).Msg("fail secp256k1 check signing")
		return nil

	} else if sigBytes == nil {
		return nil
	}

	// build the signature
	r := new(big.Int).SetBytes(sigBytes[:32])
	ss := new(big.Int).SetBytes(sigBytes[32:])
	signature := &btcec.Signature{R: r, S: ss}

	spk := &btcec.PublicKey{
		Curve: ethPubKey.Curve,
		X:     ethPubKey.X,
		Y:     ethPubKey.Y,
	}

	if !signature.Verify(dataHash[:], spk) {
		s.logger.Error().Msg("secp256k1 check signature verification failed")
		return nil
	} else {
		fmt.Println("secp256k1VerificationSignature pk.String() ----------------- ", dataHash[:], "hex",
			hex.EncodeToString(dataHash))
		s.logger.Info().Msg("secp256k1 check signature verified")
	}

	return append(sigBytes, v[0]+27)
}

func (s *Signer) sendKeygenToMap(epoch *big.Int, poolPubKey common.PubKey, blame, members []ecommon.Address, signature []byte) error {
	var keyShares []byte
	var err error
	if s.cfg.Signer.BackupKeyshares && !poolPubKey.IsEmpty() {
		keyShares, err = tss.EncryptKeyShares(
			filepath.Join(constants.DefaultHome, fmt.Sprintf("localstate-%s.json", poolPubKey.String())), // todo handler
			os.Getenv("SIGNER_SEED_PHRASE"),
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to encrypt keyShares")
		}
	}

	txID, err := s.thorchainBridge.SendKeyGenStdTx(epoch, poolPubKey, signature, keyShares, blame, members)
	if err != nil {
		return fmt.Errorf("fail to get keygen id: %w", err)
	}

	s.logger.Info().Str("txId", txID).Int64("epoch", epoch.Int64()).Msg("sent keygen tx to thorchain")
	return nil
}

// signAndBroadcast will sign the tx and broadcast it to the corresponding chain. On
// SignTx error for the chain client, if we receive checkpoint bytes we also return them
// with the error so they can be set on the TxOutStoreItem and re-used on a subsequent
// retry to avoid double spend. The second returned value is an optional observation
// that should be submitted to THORChain.
func (s *Signer) signAndBroadcast(item TxOutStoreItem) ([]byte, *types.TxInItem, error) {
	height := item.Height
	tx := item.TxOutItem

	// set the checkpoint on the tx out item if it was stored
	if item.Checkpoint != nil {
		tx.Checkpoint = item.Checkpoint
	}

	blockHeight, err := s.thorchainBridge.GetBlockHeight()
	if err != nil {
		s.logger.Error().Err(err).Msgf("fail to get block height")
		return nil, nil, err
	}
	signingTransactionPeriod, err := s.constantsProvider.GetInt64Value(blockHeight, constants.SigningTransactionPeriod)
	s.logger.Debug().Msgf("signing transaction period:%d", signingTransactionPeriod)
	if err != nil {
		s.logger.Error().Err(err).Msgf("fail to get constant value for(%s)", constants.SigningTransactionPeriod)
		return nil, nil, err
	}

	// if in round 7 retry, discard outbound if over the max outbound attempts
	inactiveVaultRound7Retry := false
	if item.Round7Retry {
		mimirKey := "MAXOUTBOUNDATTEMPTS"
		var maxOutboundAttemptsMimir int64
		maxOutboundAttemptsMimir, err = s.thorchainBridge.GetMimir(mimirKey)
		if err != nil {
			s.logger.Err(err).Msgf("fail to get %s", mimirKey)
			return nil, nil, err
		}
		attempt := (blockHeight - height) / signingTransactionPeriod
		if attempt > maxOutboundAttemptsMimir {
			s.logger.Warn().
				Int64("outbound_height", height).
				Int64("current_height", blockHeight).
				Int64("attempt", attempt).
				Msg("round 7 retry outbound tx has reached max outbound attempts")
			return nil, nil, nil
		}

		// determine if the round 7 retry is for an inactive vault
		var vault ttypes.Vault
		vault, err = s.thorchainBridge.GetVault(item.TxOutItem.VaultPubKey.String())
		if err != nil {
			log.Err(err).
				Stringer("vault_pubkey", item.TxOutItem.VaultPubKey).
				Msg("failed to get tx out item vault")
			return nil, nil, err
		}
		inactiveVaultRound7Retry = vault.Status == ttypes.VaultStatus_InactiveVault
	}

	// if not in round 7 retry or the round 7 retry is on an inactive vault, discard
	// outbound if within configured blocks of reschedule
	if !item.Round7Retry || inactiveVaultRound7Retry {
		if blockHeight-signingTransactionPeriod > height-s.cfg.Signer.RescheduleBufferBlocks {
			s.logger.Error().Msgf("tx was created at block height(%d), now it is (%d), it is older than (%d) blocks, skip it", height, blockHeight, signingTransactionPeriod)
			return nil, nil, nil
		}
	}

	chain, err := s.getChain(tx.Chain)
	if err != nil {
		s.logger.Error().Err(err).Msgf("not supported %s", tx.Chain.String())
		return nil, nil, err
	}
	mimirKey := "HALTSIGNING"
	haltSigningGlobalMimir, err := s.thorchainBridge.GetMimir(mimirKey)
	if err != nil {
		s.logger.Err(err).Msgf("fail to get %s", mimirKey)
		return nil, nil, err
	}
	if haltSigningGlobalMimir > 0 && haltSigningGlobalMimir < blockHeight {
		s.logger.Info().Msg("signing has been halted globally")
		return nil, nil, nil
	}
	mimirKey = fmt.Sprintf("HALTSIGNING%s", tx.Chain)
	haltSigningMimir, err := s.thorchainBridge.GetMimir(mimirKey)
	if err != nil {
		s.logger.Err(err).Msgf("fail to get %s", mimirKey)
		return nil, nil, err
	}
	if haltSigningMimir > 0 && haltSigningMimir < blockHeight {
		s.logger.Info().Msgf("signing for %s is halted", tx.Chain)
		return nil, nil, nil
	}
	if !s.shouldSign(tx) {
		s.logger.Info().Str("signer_address", chain.GetAddress(tx.VaultPubKey)).Msg("different pool address, ignore")
		return nil, nil, nil
	}

	if len(tx.ToAddress) == 0 {
		s.logger.Info().Msg("To address is empty, THORNode don't know where to send the fund , ignore")
		return nil, nil, nil // return nil and discard item
	}

	// don't sign if the block scanner is unhealthy. This is because the
	// network may not be able to detect the outbound transaction, and
	// therefore reschedule the transaction to another signer. In a disaster
	// scenario, the network could broadcast a transaction several times,
	// bleeding funds.
	if !chain.IsBlockScannerHealthy() {
		return nil, nil, fmt.Errorf("the block scanner for chain %s is unhealthy, not signing transactions due to it", chain.GetChain())
	}

	start := time.Now()
	defer func() {
		s.m.GetHistograms(metrics.SignAndBroadcastDuration(chain.GetChain())).Observe(time.Since(start).Seconds())
	}()

	if !tx.OutHash.IsEmpty() {
		s.logger.Info().Str("OutHash", tx.OutHash.String()).Msg("tx had been sent out before")
		return nil, nil, nil // return nil and discard item
	}

	// We get the keysign object from thorchain again to ensure it hasn't
	// been signed already, and we can skip. This helps us not get stuck on
	// a task that we'll never sign, because 2/3rds already has and will
	// never be available to sign again.
	txOut, err := s.thorchainBridge.GetTxByBlockNumber(height, tx.VaultPubKey.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get keysign items")
		return nil, nil, err
	}
	for _, txArray := range txOut.TxArray {
		if txArray.TxOutItem(item.TxOutItem.Height).Equals(tx) && !txArray.OutHash.IsEmpty() {
			// already been signed, we can skip it
			s.logger.Info().Str("tx_id", tx.OutHash.String()).Msgf("already signed. skipping...")
			return nil, nil, nil
		}
	}

	// If this is a UTXO chain, lock the vault around sign and broadcast to avoid
	// consolidate transactions from using the same UTXOs.
	if utxoClient, ok := chain.(*utxo.Client); ok {
		lock := utxoClient.GetVaultLock(tx.VaultPubKey.String())
		lock.Lock()
		defer lock.Unlock()
	}

	// If SignedTx is set, we already signed and should only retry broadcast.
	var signedTx, checkpoint []byte
	var elapse time.Duration
	var observation *types.TxInItem
	if len(item.SignedTx) > 0 {
		s.logger.Info().Str("memo", tx.Memo).Msg("retrying broadcast of already signed tx")
		signedTx = item.SignedTx
		observation = item.Observation
	} else {
		startKeySign := time.Now()
		signedTx, checkpoint, observation, err = chain.SignTx(tx, height)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to sign tx")
			return checkpoint, nil, err
		}
		elapse = time.Since(startKeySign)
	}

	// looks like the transaction is already signed
	if len(signedTx) == 0 {
		s.logger.Warn().Msgf("signed transaction is empty")
		return nil, nil, nil
	}

	// broadcast the transaction
	hash, err := chain.BroadcastTx(tx, signedTx)
	if err != nil {
		s.logger.Error().Err(err).Str("memo", tx.Memo).Msg("fail to broadcast tx to chain")

		// store the signed tx for the next retry
		item.SignedTx = signedTx
		item.Observation = observation
		if storeErr := s.storage.Set(item); storeErr != nil {
			s.logger.Error().Err(storeErr).Msg("fail to update tx out store item with signed tx")
		}

		return nil, observation, err
	}
	s.logger.Info().Str("txid", hash).Str("memo", tx.Memo).Msg("broadcasted tx to chain")

	if s.isTssKeysign(tx.VaultPubKey) {
		s.tssKeysignMetricMgr.SetTssKeysignMetric(hash, elapse.Milliseconds())
	}

	return nil, observation, nil
}

func (s *Signer) isTssKeysign(pubKey common.PubKey) bool {
	return !s.localPubKey.Equals(pubKey)
}

// Stop the signer process
func (s *Signer) Stop() error {
	s.logger.Info().Msg("receive request to stop signer")
	defer s.logger.Info().Msg("signer stopped successfully")
	close(s.stopChan)
	s.wg.Wait()
	if err := s.m.Stop(); err != nil {
		s.logger.Error().Err(err).Msg("fail to stop metric server")
	}
	s.blockScanner.Stop()
	return s.storage.Close()
}

////////////////////////////////////////////////////////////////////////////////////////
// pipelineSigner Interface
////////////////////////////////////////////////////////////////////////////////////////

func (s *Signer) isStopped() bool {
	select {
	case <-s.stopChan:
		return true
	default:
		return false
	}
}

func (s *Signer) storageList() []TxOutStoreItem {
	return s.storage.List()
}

func (s *Signer) processTransaction(item TxOutStoreItem) {
	s.logger.Info().
		Int64("height", item.Height).
		Int("status", int(item.Status)).
		Interface("tx", item.TxOutItem).
		Msg("Signing transaction")

	// a single keysign should not take longer than 5 minutes , regardless TSS or local
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	checkpoint, obs, err := runWithContext(ctx, func() ([]byte, *types.TxInItem, error) {
		// todo will next 400
		return s.signAndBroadcast(item)
	})
	if err != nil {
		// mark the txout on round 7 failure to block other txs for the chain / pubkey
		ksErr := tss.KeysignError{}
		if errors.As(err, &ksErr) && ksErr.IsRound7() {
			s.logger.Error().Err(err).Interface("tx", item.TxOutItem).Msg("round 7 signing error")
			item.Round7Retry = true
			item.Checkpoint = checkpoint
			if storeErr := s.storage.Set(item); storeErr != nil {
				s.logger.Error().Err(storeErr).Msg("fail to update tx out store item with round 7 retry")
			}
		}

		s.logger.Error().Interface("tx", item.TxOutItem).Err(err).Msg("fail to sign and broadcast tx out store item")
		cancel()
		return
		// The 'item' for loop should not be items[0],
		// because problems which return 'nil, nil' should be skipped over instead of blocking others.
		// When signAndBroadcast returns an error (such as from a keysign timeout),
		// a 'return' and not a 'continue' should be used so that nodes can all restart the list,
		// for when the keysign failure was from a loss of list synchrony.
		// Otherwise, out-of-sync lists would cycle one timeout at a time, maybe never resynchronising.
	}
	cancel()

	// if enabled and the observation is non-nil, instant observe the outbound
	if s.cfg.Signer.AutoObserve && obs != nil {
		s.observer.ObserveSigned(types.TxIn{
			Chain:                item.TxOutItem.Chain,
			TxArray:              []*types.TxInItem{obs},
			MemPool:              true,
			Filtered:             true,
			ConfirmationRequired: 0,

			// Instant EVM observations have wrong gas and need future correct observations
			AllowFutureObservation: item.TxOutItem.Chain.IsEVM(),
		})
	}

	// We have a successful broadcast! Remove the item from our store
	if err = s.storage.Remove(item); err != nil {
		s.logger.Error().Err(err).Msg("fail to update tx out store item")
	}
}
