package signer

import (
	"fmt"
	"time"

	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	"github.com/mapprotocol/compass-tss/common"

	"github.com/mapprotocol/compass-tss/x/types"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////////////
// Internal Types
////////////////////////////////////////////////////////////////////////////////////////

// vaultChain is a public key and chain used as a key for the vault/chain lock.
type vaultChain struct {
	Vault common.PubKey
	Chain string
}

type semaphore chan struct{}

// acquire will asynchronously acquire all available capacity from the semaphore.
func (s semaphore) acquire() int {
	count := 0
	for {
		select {
		case s <- struct{}{}:
			count++
		default:
			return count
		}
	}
}

// release will release the provided count to the semaphore.
func (s semaphore) release(count int) {
	for i := 0; i < count; i++ {
		<-s
	}
}

// pipelineSigner is the signer interface required for the pipeline.
type pipelineSigner interface {
	isStopped() bool
	storageList() []TxOutStoreItem
	processTransaction(item TxOutStoreItem)
}

////////////////////////////////////////////////////////////////////////////////////////
// pipeline
////////////////////////////////////////////////////////////////////////////////////////

type pipeline struct {
	// concurrency is the number of concurrent signing routines to allow.
	concurrency int64

	// vaultStatusConcurrency maps vault status to a semaphore for concurrent signings.
	vaultStatusConcurrency map[types.VaultStatus]semaphore

	// vaultChainLock maps a vault/chain combination to a lock. The lock is represented as
	// a channel instead of a mutex so we can check if it is taken without blocking.
	vaultChainLock map[vaultChain]chan struct{}
}

// NewPipeline creates a new pipeline instance using the provided concurrency for active
// and retiring vault status semaphores. The inactive vault status semaphore will always
// be 1 - allowing only 1 concurrent signing routine for inactive vault refunds.
func newPipeline(concurrency int64) (*pipeline, error) {
	log.Info().Int64("concurrency", concurrency).Msg("creating new signer pipeline")

	if concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be greater than 0")
	}

	return &pipeline{
		concurrency: concurrency,
		vaultStatusConcurrency: map[types.VaultStatus]semaphore{
			types.VaultStatus_ActiveVault:   make(semaphore, int(concurrency)),
			types.VaultStatus_RetiringVault: make(semaphore, int(concurrency)),
			types.VaultStatus_InactiveVault: make(semaphore, 1),
		},
		vaultChainLock: make(map[vaultChain]chan struct{}),
	}, nil
}

// SpawnSiginings will fetch all transactions from the provided Signer's storage, and
// start signing routines for any transactions that have:
//  1. Sufficient capacity in the vault status semaphore for the source vault's status.
//  2. An available lock on the vault/chain combination (only 1 can run at a time).
//
// The signing routines will be spawned in a goroutine, and this function will not
// block on their completion. The spawned routines will release the corresponding vault
// status semaphore and vault/chain lock when they are complete.
func (p *pipeline) SpawnSignings(s pipelineSigner, bridge shareTypes.Bridge) {
	allItems := s.storageList()

	// gather all vault/chain combinations with an out item in retry
	retryItems := make(map[vaultChain][]TxOutStoreItem)
	for _, item := range allItems {
		if item.Round7Retry || len(item.SignedTx) > 0 {
			vc := vaultChain{item.TxOutItem.VaultPubKey, item.TxOutItem.Chain.String()}
			retryItems[vc] = append(retryItems[vc], item)
		}
	}
	log.Info().Msgf("SpawnSignings found %d all items", len(allItems))

	var itemsToSign []TxOutStoreItem

	// add retry items to our items to sign
	for _, items := range retryItems {
		// there should be no vault/chain with more than 1 item in retry
		if len(items) > 1 {
			for i := range items { // sanitize signed tx for log
				items[i].SignedTx = nil
			}
			log.Error().
				Interface("items", items).
				Msg("found multiple retry items for vault/chain")
		} else {
			itemsToSign = append(itemsToSign, items[0])
			log.Warn().
				Interface("items", items).
				Msg("found retry items")
		}
	}

	// add all items from vault/chains with no items in retry
	for _, item := range allItems {
		vc := vaultChain{item.TxOutItem.VaultPubKey, item.TxOutItem.Chain.String()}
		if _, ok := retryItems[vc]; !ok {
			itemsToSign = append(itemsToSign, item)
		}
	}

	// get the available capacities for each vault status
	availableCapacities := make(map[types.VaultStatus]int)
	for status, semaphore := range p.vaultStatusConcurrency {
		availableCapacities[status] = semaphore.acquire()
	}

	// release remaining capacity for each vault status on return
	defer func() {
		for status, capacity := range availableCapacities {
			p.vaultStatusConcurrency[status].release(capacity)
		}
	}()

	// get all locked vault/chains - otherwise races if a vault/chain unlocks mid-iteration
	lockedVaultChains := make(map[vaultChain]bool)
	for vc, lock := range p.vaultChainLock {
		if len(lock) > 0 {
			lockedVaultChains[vc] = true
		}
	}
	log.Info().Msgf("SpawnSignings will handle tx locking")

	// spawn signing routines for each item
	for _, item := range itemsToSign {
		log.Info().Msgf("SpawnSignings will handler %s a tx", item.TxOutItem.TxHash)
		// return if the signer is stopped
		if s.isStopped() {
			return
		}

		vc := vaultChain{item.TxOutItem.VaultPubKey, item.TxOutItem.Chain.String()}

		// check if the vault/chain is locked
		if lockedVaultChains[vc] {
			continue
		}

		// if no lock exists, create one
		if _, ok := p.vaultChainLock[vc]; !ok {
			p.vaultChainLock[vc] = make(chan struct{}, 1)
		}

		// // get vault to determine vault status
		// vault, err := bridge.GetVault(item.TxOutItem.VaultPubKey.String())
		// if err != nil {
		// 	log.Err(err).
		// 		Stringer("vault_pubkey", item.TxOutItem.VaultPubKey).
		// 		Msg("Failed to get tx out item vault")
		// 	return
		// }

		// check if the vault status semaphore has capacity
		// todo next
		if availableCapacities[types.VaultStatus_ActiveVault] == 0 {
			continue
		}

		// acquire the vault status semaphore and vault/chain lock
		availableCapacities[types.VaultStatus_ActiveVault]--
		p.vaultChainLock[vc] <- struct{}{}
		lockedVaultChains[vc] = true

		// spawn signing routine
		go func(item TxOutStoreItem, vaultStatus types.VaultStatus) {
			// release the vault status semaphore and vault/chain lock when complete
			defer func() {
				vc2 := vaultChain{item.TxOutItem.VaultPubKey, item.TxOutItem.Chain.String()}
				<-p.vaultChainLock[vc2]
				p.vaultStatusConcurrency[vaultStatus].release(1)
			}()

			// process the transaction
			s.processTransaction(item)
		}(item, types.VaultStatus_ActiveVault)
	}
}

// Wait will block until all pipeline signing routines have completed.
func (p *pipeline) Wait() {
	log.Info().Msg("waiting for signer pipeline routines to complete")
	for {
		running := false
		for _, semaphore := range p.vaultStatusConcurrency {
			if len(semaphore) > 0 {
				running = true
				break
			}
		}
		if !running {
			log.Info().Msg("signer pipeline routines complete")
			return
		}
		time.Sleep(time.Second)
	}
}
