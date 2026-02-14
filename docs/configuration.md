# Configuration Documentation

## Overview

This document outlines the configuration options available in the Compass TSS project. The configuration is structured
hierarchically and includes settings for various blockchain integrations, LevelDB options, metrics, and TSS-related
parameters.

Default values are defined in [`config/default.yaml`](../config/default.yaml).

---

## Global Settings

### Observer LevelDB (`observer_leveldb`)

Default LevelDB configuration shared across components:

- `filter_bits_per_key`: Bloom filter bits per key.
- `compaction_table_size_multiplier`: Multiplier for table size during compaction.
- `write_buffer`: Size of write buffer in bytes.
- `block_cache_capacity`: Size of block cache in bytes.
- `compact_on_init`: Trigger full compaction at initialization.

### Observer Workers (`observer_workers`)

Number of goroutines handling cross-chain transaction storage.

---

## Metrics Configuration (`metrics`)

Settings for monitoring and profiling:

- `enabled`: Enable metrics collection.
- `pprof_enabled`: Enable pprof profiling.
- `listen_port`: Port for metrics server.
- `read_timeout`: Read timeout duration, A zero or negative value means there will be no timeout.
- `write_timeout`: Write timeout duration, No timeout is applied if Timeout is 0 or negative.
- `chains`: List of supported chains for metrics.

---

## MAP Relay Configuration (`mapo`)

Configuration for interacting with the MAP chain:

- `chain_id`: Chain identifier.
- `chain_host`: RPC endpoint for the chain.
- `signer_name`: Name of the signer.
- `keystore_path`: Path to keystore file.
- `maintainer`: Maintainer contract address.
- `tss_manager`: TSS manager contract address.
- `relay`: Relay contract address.
- `gas_service`: Gas service contract address.
- `token_registry`: Token registry contract address.
- `view_controller`: View controller contract address.
- `affiliate_fee_manager`: Affiliate fee manager contract address.
- `fusion_receiver`: Fusion receiver contract address.
- `cross_data_path`: Path for cross-chain data storage.
- `cross_data_address`: Address for cross-chain data service.
- `increase_gas_limit`: Additional gas limit for transactions, If increase_gas_limit is not 0, the final gas limit is
  the estimated gas limit plus increase_gas_limit.

---

## Signer Configuration (`signer`)

Settings for managing keyshares and signing operations:

- `backup_keyshares`: Backup key shares locally
- `signer_db_path`: Path to signer database
- `retry_interval`: Interval between retries
- `reschedule_buffer_blocks`: Buffer blocks before rescheduling
- `block_scanner`: [Block scanner](#block-scanner-block_scanner)) settings
- `leveldb`: [LevelDB](#observer-leveldb-observer_leveldb) settings
- `keygen_timeout`: Timeout for key generation
- `keysign_timeout`: Timeout for key signing
- `party_timeout`: Timeout for party coordination
- `pre_param_timeout`: Timeout for pre-parameter setup

#### Block Scanner (`block_scanner`)

- `max_reorg_rescan_blocks`: Maximum number of blocks to rescan during a chain reorganization
- `chain_id`: Identifier for the blockchain.
- `block_height_discover_back_off`: Backoff time when discovering block heights
- `observation_flexibility_blocks`: Number of blocks behind the tip to submit network fee and solvency observations
- `http_request_timeout`: Timeout for HTTP requests
- `db_path`: Path to the database storage
- `scan_mempool`: Whether to scan mempool transactions
- `start_block_height`: Starting block height
- `gas_cache_blocks`: Number of blocks to cache gas price data
- `concurrency`: Number of concurrent goroutines used for RPC requests on data within a block - e.g. transactions,
  receipts, logs, etc. Blocks are processed sequentially.
- `max_gas_limit`: Maximum gas limit for outbound transactions.
- `max_swap_gas_limit`: Maximum gas limit for swap transactions.
- `fixed_gas_rate`: Fixed gas rate to enforce.
- `gas_price_resolution`: Resolution of gas price per unit (in wei, satoshi).
- `max_resume_block_lag`: Maximum lag behind the latest consensus inbound height upon startup.
- `max_healthy_lag`: Maximum lag before the scanner is considered unhealthy.
- `transaction_batch_size`: Number of transactions to batch in a single request.
- `gateway`: Gateway contract address.

---

## TSS Configuration (`tss`)

Settings for Threshold Signature Scheme (TSS):

- `rendezvous`: Rendezvous string for peer discovery, used to isolate different networks.
- `p2p_port`: P2P communication port.
- `info_address`: Address for TSS info service.
- `bootstrap_peers`: List of bootstrap peers.
- `external_ip`: External IP address.

---

## Chain-Specific Configurations (`chains`)

Each chain supports the following general settings:

- `disabled`: Disable the chain integration.
- `chain_id`: Unique identifier for the chain.
- `username`: Username for RPC authentication.
- `password`: Password for RPC authentication.
- `rpc_host`: RPC endpoint for the chain.
- `mempool_tx_id_cache_size`: Number of transaction ids to cache in memory.
- `scanner_leveldb`: [LevelDB](#observer-leveldb-observer_leveldb) options for the scanner.
- `min_confirmations`: Minimum confirmations required.
- `max_rpc_retries`: Maximum number of retries for RPC requests.
- `max_pending_nonces`: Maximum number of pending nonces to allow before aborting new signing attempts.
- `authorization_bearer`: Authorization token for RPC.
- `limit_multiplier`: Multiplier for gas limits.
- `max_gas_limit`: Maximum gas limit.
- `utxo`: [UTXO-specific](#utxo-configurations-utxo) settings (Bitcoin-like chains).
- `block_scanner`: [Block scanner](#block-scanner-block_scanner)) settings.

---

## UTXO Configurations (`utxo`)

UTXO chain specific configuration.

- `block_cache_count`: The number of blocks to cache in storage.
- `transaction_batch_size`: The number of transactions to batch in a single request.
- `max_mempool_batches`: The maximum number of mempool batches to fetch from the mempool in a single scanning pass.
- `estimated_average_tx_size`: The estimated average size (in bytes) of a transaction.
- `estimated_average_tx_swap_size`: The estimated average size (in bytes) of a swap transaction.
- `default_min_relay_fee_sats`: The default minimum relay fee (in satoshis).
- `default_sats_per_vbyte`: The default fee rate in sats per vbyte. It is only used as a fallback when local fee
  information is unavailable.
- `max_sats_per_vbyte`: The maximum fee rate in sats per vbyte. It is used to cap the fee rate, since overpaid fees may
  be rejected by the chain daemon.
- `min_sats_per_vbyte`: The minimum fee rate in satoshis per virtual byte to ensure transactions are not stuck in the
  mempool.
- `min_utxo_confirmations`: The minimum number of confirmations required for a UTXO to be considered spendable.
- `max_utxos_to_spend`: The maximum number of UTXOs that can be spent in a single transaction.

---