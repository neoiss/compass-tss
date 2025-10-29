package signer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cenkalti/backoff"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Signer will pull the tx out from thorchain and then forward it to chain
type Signer struct {
	logger               zerolog.Logger
	cfg                  config.Bifrost
	wg                   *sync.WaitGroup
	mapBridge            shareTypes.Bridge
	stopChan             chan struct{}
	blockScanner         *blockscanner.BlockScanner
	mapChainBlockScanner *mapo.MapChainBlockScan
	chains               map[common.Chain]chainclients.ChainClient
	storage              SignerStorage
	oracleStorage        SignerStorage
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
	storage, err := NewSignerStore(cfg.Signer.SignerDbPath, cfg.Signer.LevelDB)
	if err != nil {
		return nil, fmt.Errorf("fail to create thorchain scan storage: %w", err)
	}

	oracleStorage, err := NewSignerStore(cfg.Signer.OracleDbPath, cfg.Signer.LevelDB)
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
		oracleStorage:        oracleStorage,
		errCounter:           m.GetCounterVec(metrics.SignerError),
		pubkeyMgr:            pubkeyMgr,
		mapBridge:            bridge,
		tssKeygen:            kg,
		tssServer:            tssServer,
		constantsProvider:    constantProvider,
		localPubKey:          selfKey,
		tssKeysignMetricMgr:  tssKeysignMetricMgr,
		observer:             obs,
	}, nil
}

func (s *Signer) getChain(chainID *big.Int) (chainclients.ChainClient, error) {
	chainName, ok := common.GetChainName(chainID)
	if !ok {
		s.logger.Debug().Str("chain", chainID.String()).Msg("Is not supported yet")
		return nil, errors.New("not supported")
	}
	chain, ok := s.chains[chainName]
	if !ok {
		s.logger.Debug().Str("chain", chainID.String()).Msg("Is not supported yet")
		return nil, errors.New("not supported")
	}
	return chain, nil
}

// Start signer process
func (s *Signer) Start() error {
	//  todo handler annotate
	s.wg.Add(1)
	go s.processTxnOut(s.mapChainBlockScanner.GetTxOutMessages()) // cache local

	s.wg.Add(1)
	go s.processKeygen(s.mapChainBlockScanner.GetKeygenMessages())

	s.wg.Add(1)
	go s.cacheOracle(s.mapChainBlockScanner.GetOracleMessages())

	s.wg.Add(1)
	go s.processOracle()

	s.wg.Add(1)
	go s.signTransactions()

	s.blockScanner.Start(nil, nil)
	return nil
}

// func (s *Signer) shouldSign(tx types.TxOutItem) bool {
// 	return s.pubkeyMgr.HasPubKey(tx.VaultPubKey)
// }

// signTransactions - looks for work to do by getting a list of all unsigned
// transactions stored in the storage
func (s *Signer) signTransactions() {
	s.logger.Info().Msg("Start to sign transactions")
	defer s.logger.Info().Msg("Stop to sign transactions")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		default:
			// When map relay chain is catching up , bifrost might get stale data from compass-tss , thus it shall pause signing
			catchingUp, err := s.mapBridge.IsSyncing()
			if err != nil {
				s.logger.Error().Err(err).Msg("Fail to get thorchain sync status")
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
	signerConcurrency, err := s.mapBridge.GetMimir(constants.SignerConcurrency.String())
	if err != nil {
		s.logger.Error().Err(err).Msg("Fail to get signer concurrency mimir")
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
			s.logger.Error().Err(err).Msg("Fail to create new pipeline")
			return
		}
	}

	// process transactions
	s.pipeline.SpawnSignings(s, s.mapBridge)
}

// processTxnOut processes outbound TxOuts and save them to storage
func (s *Signer) processTxnOut(ch <-chan types.TxOut) {
	s.logger.Info().Msg("Start to process tx out")
	defer s.logger.Info().Msg("Stop to process tx out")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		case txOut, more := <-ch:
			if !more {
				return
			}
			s.logger.Info().Msgf("Received a TxOut Array of %v from the MAP", txOut)
			items := make([]TxOutStoreItem, 0, len(txOut.TxArray))

			for i, tx := range txOut.TxArray {
				items = append(items, NewTxOutStoreItem(txOut.Height, tx.TxOutItem(txOut.Height), int64(i)))
			}
			if err := s.storage.Batch(items); err != nil {
				s.logger.Error().Err(err).Msg("Fail to save tx out items to storage")
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

			s.logger.Info().Interface("keygenBlock", keygenBlock).Msg("Received a keygen block from map relay")
			s.processKeygenBlock(keygenBlock)
		}
	}
}

func (s *Signer) cacheOracle(ch <-chan types.TxOut) {
	s.logger.Info().Msg("Start to cache tx oracle")
	defer s.logger.Info().Msg("Stop to cache tx oracle")
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		case txOut, ok := <-ch:
			if !ok {
				return
			}
			s.logger.Info().Msgf("Oracle Received a TxOut Array of %v from the MAPRelay", txOut)
			items := make([]TxOutStoreItem, 0, len(txOut.TxArray))

			for i, tx := range txOut.TxArray {
				items = append(items, NewTxOutStoreItem(txOut.Height, tx.TxOutItem(txOut.Height), int64(i)))
			}
			if err := s.oracleStorage.Batch(items); err != nil {
				s.logger.Error().Err(err).Msg("Fail to save tx out items to storage")
			}
		}
	}
}

func (s *Signer) processOracle() {
	s.logger.Info().Msg("Start to process tx oracle")
	defer s.logger.Info().Msg("Stop to process tx oracle")
	defer s.wg.Done()
	for {
		select {
		case <-s.stopChan:
			return
		default:
			list := s.oracleStorage.List() // this will trigger the storage to load all items
			for _, item := range list {
				txBytes, err := s.mapBridge.GetOracleStdTx(&item.TxOutItem)
				if err != nil {
					s.logger.Error().Err(err).Msg("Fail to get oracle std tx")
					continue
				}

				if len(txBytes) == 0 {
					continue
				}

				tmp := item
				bf := backoff.NewExponentialBackOff()
				bf.MaxElapsedTime = 5 * time.Second
				err = backoff.Retry(func() error {
					txID, err := s.mapBridge.Broadcast(txBytes)
					if err != nil {
						return fmt.Errorf("fail to send the tx to thorchain: %w", err)
					}
					s.oracleStorage.Remove(tmp)
					s.logger.Info().Str("mapHash", txID).Msg("Oracle tx sent successfully")
					return nil
				}, bf)
				if err != nil {
					s.logger.Error().Err(err).Msg("Fail to broadcast tx")
					continue
				}

				s.logger.Info().Interface("item", item).Msg("Processing oracle item")
			}
		}
	}
}

func (s *Signer) processKeygenBlock(keygenBlock *structure.KeyGen) {
	if len(keygenBlock.Ms) <= 0 {
		s.logger.Info().Msg("Processing keygen block, members is zero")
		return
	}
	members := make(common.PubKeys, 0, len(keygenBlock.Ms))
	memberAddrs := make([]ecommon.Address, 0, len(keygenBlock.Ms))
	for _, ele := range keygenBlock.Ms {
		if ele.Account.String() == "" {
			continue
		}
		members = append(members, common.PubKey(ecommon.Bytes2Hex(ele.Secp256Pubkey)))
		memberAddrs = append(memberAddrs, ele.Account)
	}
	// NOTE: in practice there is only one keygen in the keygen block
	keygenStart := time.Now()
	// debug blame
	pubKey, blame, err := s.tssKeygen.GenerateNewKey(keygenBlock.Epoch.Int64(), members)
	if !blame.IsEmpty() {
		s.logger.Error().Str("reason", blame.FailReason).
			Interface("nodes", blame.BlameNodes).Msg("Keygen blame")
	}
	keygenTime := time.Since(keygenStart).Milliseconds()
	if err != nil {
		s.errCounter.WithLabelValues("fail_to_keygen_pubkey", "").Inc()
		s.logger.Error().Err(err).Msg("Fail to generate new pubkey")
	}

	s.logger.Info().Int64("keygenTime", keygenTime).Msg("ProcessKeygenBlock keyGen time")
	secp256k1Sig := make([]byte, 0)
	if len(pubKey.Secp256k1.String()) > 0 {
		secp256k1Sig = s.secp256k1VerificationSignature(pubKey.Secp256k1)
	}
	blames := make([]ecommon.Address, 0)
	if len(blame.BlameNodes) > 0 {
		for _, node := range blame.BlameNodes {
			addr, err := keys.GetAddressByCompressPk(node.Pubkey)
			if err != nil {
				continue
			}
			blames = append(blames, addr)
		}
	}
	err = s.sendKeygenToMap(keygenBlock.Epoch, pubKey.Secp256k1, blames, memberAddrs, secp256k1Sig)
	if err != nil { // handler blame
		s.errCounter.WithLabelValues("fail_to_broadcast_keygen", "").Inc()
		s.logger.Error().Err(err).Msg("Fail to broadcast keygen")
	}

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
	ks, err := tss.NewKeySign(s.tssServer, s.mapBridge)
	if err != nil {
		s.logger.Error().Err(err).Msg("Fail to create keySign for secp256k1 check signing")
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
	if err != nil {
		// this is expected in some cases if we were not in the signing party
		s.logger.Info().Err(err).Msg("Fail secp256k1 check signing")
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
		s.logger.Error().Msg("Secp256k1 check signature verification failed")
		return nil
	} else {
		s.logger.Info().Msg("Secp256k1 check signature verified")
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
			s.logger.Error().Err(err).Msg("Fail to encrypt keyShares")
		}
	}

	txID, err := s.mapBridge.SendKeyGenStdTx(epoch, poolPubKey, signature, keyShares, blame, members)
	if err != nil {
		return fmt.Errorf("fail to get keygen id: %w", err)
	}

	s.logger.Info().Str("txId", txID).Int64("epoch", epoch.Int64()).Msg("Send keygen tx to relay")
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

	//  utxo
	pubKeys, _ := s.mapBridge.GetAsgardPubKeys()
	tx.VaultPubKey = pubKeys[0].PubKey

	blockHeight, err := s.mapBridge.GetBlockHeight()
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
		maxOutboundAttemptsMimir, err = s.mapBridge.GetMimir(mimirKey)
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

		// // determine if the round 7 retry is for an inactive vault
		// var vault ttypes.Vault
		// vault, err = s.mapBridge.GetVault(item.TxOutItem.VaultPubKey.String())
		// if err != nil {
		// 	log.Err(err).
		// 		Stringer("vault_pubkey", item.TxOutItem.VaultPubKey).
		// 		Msg("failed to get tx out item vault")
		// 	return nil, nil, err
		// }
		// inactiveVaultRound7Retry = vault.Status == ttypes.VaultStatus_InactiveVault
	}

	// if not in round 7 retry or the round 7 retry is on an inactive vault, discard
	// outbound if within configured blocks of reschedule
	if !item.Round7Retry || inactiveVaultRound7Retry {
		if blockHeight-signingTransactionPeriod > height-s.cfg.Signer.RescheduleBufferBlocks {
			s.logger.Error().Msgf("Tx was created at block height(%d), now it is (%d), it is older than (%d) blocks, skip it", height, blockHeight, signingTransactionPeriod)
			return nil, nil, nil
		}
	}

	chain, err := s.getChain(tx.Chain)
	if err != nil {
		s.logger.Error().Err(err).Msgf("Not supported %s", tx.Chain.String())
		return nil, nil, err
	}
	mimirKey := "HALTSIGNING"
	haltSigningGlobalMimir, err := s.mapBridge.GetMimir(mimirKey)
	if err != nil {
		s.logger.Err(err).Msgf("Fail to get %s", mimirKey)
		return nil, nil, err
	}
	if haltSigningGlobalMimir > 0 && haltSigningGlobalMimir < blockHeight {
		s.logger.Info().Msg("Signing has been halted globally")
		return nil, nil, nil
	}
	mimirKey = fmt.Sprintf("HALTSIGNING%s", tx.Chain)
	haltSigningMimir, err := s.mapBridge.GetMimir(mimirKey)
	if err != nil {
		s.logger.Err(err).Msgf("Fail to get %s", mimirKey)
		return nil, nil, err
	}
	if haltSigningMimir > 0 && haltSigningMimir < blockHeight {
		s.logger.Info().Msgf("Signing for %s is halted", tx.Chain)
		return nil, nil, nil
	}
	// if !s.shouldSign(tx) {
	// 	s.logger.Info().Str("signer_address", chain.GetAddress(tx.VaultPubKey)).Msg("different pool address, ignore")
	// 	return nil, nil, nil
	// }

	if len(tx.To) == 0 {
		s.logger.Info().Msg("To address is empty, map don't know where to send the fund , ignore")
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

	// If this is a UTXO chain, lock the vault around sign and broadcast to avoid
	// consolidate transactions from using the same UTXOs.
	if utxoClient, ok := chain.(*utxo.Client); ok {
		lock := utxoClient.GetVaultLock(string(tx.Vault)) // todo will next2
		// ensure vault rule
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
		s.logger.Warn().Msgf("Signed transaction is empty")
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
	s.logger.Info().Str("txId", hash).Str("memo", tx.Memo).Msg("broadcasted tx to chain")

	s.tssKeysignMetricMgr.SetTssKeysignMetric(hash, elapse.Milliseconds())

	return nil, observation, nil
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
	s.logger.Info().Int64("height", item.Height).Int("status", int(item.Status)).
		Interface("tx", item.TxOutItem).Msg("signing transaction")

	// a single keysign should not take longer than 5 minutes , regardless TSS or local
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	checkpoint, _, err := runWithContext(ctx, func() ([]byte, *types.TxInItem, error) {
		// todo will next 400
		return s.signAndBroadcast(item)
	})
	if err != nil {
		// mark the txout on round 7 failure to block other txs for the chain / pubkey
		ksErr := tss.KeysignError{}
		if errors.As(err, &ksErr) && ksErr.IsRound7() {
			s.logger.Error().Err(err).Interface("tx", item.TxOutItem).Msg("Round 7 signing error")
			item.Round7Retry = true
			item.Checkpoint = checkpoint
			if storeErr := s.storage.Set(item); storeErr != nil {
				s.logger.Error().Err(storeErr).Msg("Fail to update tx out store item with round 7 retry")
			}
		}

		s.logger.Error().Interface("tx", item.TxOutItem).Err(err).Msg("Fail to sign and broadcast tx out store item")
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
	// We have a successful broadcast! Remove the item from our store
	if err = s.storage.Remove(item); err != nil {
		s.logger.Error().Err(err).Msg("fail to update tx out store item")
	}
}
