# Compass-TSS (CrossX)

**A multi-party secure validation node for cross-chain transactions on mainnet.**

Compass-TSS (CrossX) is forked from [THORChain's Bifrost](https://gitlab.com/thorchain/thornode/-/tree/develop/bifrost), adapted for the MAP Protocol ecosystem. It serves as a core component of cross-chain infrastructure, participating in cross-chain transaction validation through Threshold Signature Scheme (TSS) technology. Multiple CrossX nodes form a validation network that jointly signs and confirms cross-chain messages, ensuring secure and trustworthy asset transfers between different blockchains while eliminating single points of failure and private key exposure risks.

## Supported Chains

| Chain | Type | Status      |
|-------|------|-------------|
| BTC   | UTXO | Production  |
| ETH   | EVM  | Production  |
| BSC   | EVM  | Production  |
| BASE  | EVM  | Production  |
| ARB   | EVM  | In Progress |
| XRP   | -    | In Progress |
| DOGE  | UTXO | In Progress |




## Prerequisites

- Go 1.23.4+
- Access to chain RPC endpoints for the chains you want to monitor
- TSS bootstrap peers for joining the validator network

## Configuration

Compass-TSS uses a default configuration file at [`config/default.yaml`](config/default.yaml). Most settings can be overridden via environment variables.

See [docs/configuration.md](docs/configuration.md) for the full environment variable reference.

### Required Environment Variables

Before running, you **must** set the following environment variables:

```bash
# MAP Relay Chain
export SIGNER_NAME=<your_signer_name>
export SIGNER_PASSWD=<your_signer_password>
export CHAIN_ID=<map_relay_chain_id>
export CHAIN_API=<map_relay_rpc_address>

# TSS
export PEER=<bootstrap_peer_addresses>
export EXTERNAL_IP=<your_node_external_ip>

# Chain RPC endpoints (at least one chain is required)
export BTC_HOST=<btc_rpc_host>
export ETH_HOST=<eth_rpc_host>
# ... see docs/configuration.md for all chains
```

### Optional Environment Variables

```bash
# Block scanner retry backoff
export BLOCK_SCANNER_BACKOFF=<seconds>

# Start block heights
export BTC_START_BLOCK_HEIGHT=<height>
export ETH_START_BLOCK_HEIGHT=<height>

# Disable specific chains
export DOGE_DISABLED=true
export AVAX_DISABLED=true
```

## Build & Run

```bash
git clone https://github.com/mapprotocol/compass-tss.git
cd compass-tss/cmd/compass
go mod tidy
go build -o compass .
./compass
```

### Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--log-level, -l` | Log level (debug, info, warn, error) | `info` |
| `--pretty-log, -p` | Enable pretty console logging | `false` |
| `--version` | Show version | - |

## Project Structure

```
compass-tss/
├── cmd/compass/     # Application entry point
├── blockscanner/    # Multi-chain block scanning
├── observer/        # Chain event observation
├── signer/          # Transaction signing
├── tss/             # TSS key management
├── p2p/             # P2P network communication
├── mapclient/       # MAP Relay Chain client
├── config/          # Default configuration (default.yaml)
├── common/          # Shared types and utilities
├── txscript/        # Transaction script handling (BTC/DOGE/LTC/BCH)
└── docs/            # Documentation
```

## License

[MIT](LICENSE)
