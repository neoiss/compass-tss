# Configuration and Environment Variables

Default values are defined in [`config/default.yaml`](../config/default.yaml). Most settings can be overridden via environment variables.

## 1. MAP Relay Chain

| Configuration Item           | Environment Variable | Default                    | Description                 |
|------------------------------|----------------------|----------------------------|-----------------------------|
| `bifrost.mapo.signer_name`  | `SIGNER_NAME`       | `first_cosmos_name`        | Signer name                 |
| `bifrost.mapo.chain_id`     | `CHAIN_ID`          | `Map`                      | MAP Relay Chain ID          |
| `bifrost.mapo.chain_host`   | `CHAIN_API`         | `https://rpc.maplabs.io`  | MAP Relay Chain RPC Address |

---

## 2. TSS

| Configuration Item            | Environment Variable | Default                                                      | Description         |
|-------------------------------|----------------------|--------------------------------------------------------------|---------------------|
| `bifrost.tss.bootstrap_peers` | `PEER`              | `8.219.235.57,8.219.58.134,8.219.253.63,8.219.243.18`       | TSS Bootstrap Peers |
| `bifrost.tss.external_ip`     | `EXTERNAL_IP`       |                                                              | Node External IP    |

---

## 3. Block Scanner Backoff

Wait time after block/transaction fetch failures. Applies to all chains.

| Configuration Item                                                 | Environment Variable    | Default |
|--------------------------------------------------------------------|-------------------------|---------|
| `bifrost.chains.*.block_scanner.block_height_discover_back_off`    | `BLOCK_SCANNER_BACKOFF` | `5s`    |

---

## 4. Chain Configuration

### 4.1 BTC

| Configuration Item                                    | Environment Variable       | Default                     |
|-------------------------------------------------------|----------------------------|-----------------------------|
| `bifrost.chains.BTC.disabled`                         | `BTC_DISABLED`             | `false`                     |
| `bifrost.chains.BTC.rpc_host`                         | `BTC_HOST`                 | `https://bitcoin-rpc.publicnode.com`  |
| `bifrost.chains.BTC.username`                         | `BTC_USERNAME`             |                             |
| `bifrost.chains.BTC.password`                         | `BTC_PASSWORD`             |                             |
| `bifrost.chains.BTC.authorization_bearer`             | `BTC_AUTHORIZATION_BEARER` |                             |
| `bifrost.chains.BTC.block_scanner.start_block_height` | `BTC_START_BLOCK_HEIGHT`  |                             |

### 4.2 ETH

| Configuration Item                                    | Environment Variable       | Default                                              |
|-------------------------------------------------------|----------------------------|------------------------------------------------------|
| `bifrost.chains.ETH.disabled`                         | `ETH_DISABLED`             | `false`                                              |
| `bifrost.chains.ETH.rpc_host`                         | `ETH_HOST`                 | `https://eth.drpc.org`     |
| `bifrost.chains.ETH.username`                         | `ETH_USERNAME`             |                                                      |
| `bifrost.chains.ETH.password`                         | `ETH_PASSWORD`             |                                                      |
| `bifrost.chains.ETH.authorization_bearer`             | `ETH_AUTHORIZATION_BEARER` |                                                      |
| `bifrost.chains.ETH.block_scanner.start_block_height` | `ETH_START_BLOCK_HEIGHT`  | `9763169`                                            |

### 4.3 BSC

| Configuration Item                                    | Environment Variable       | Default                          |
|-------------------------------------------------------|----------------------------|----------------------------------|
| `bifrost.chains.BSC.disabled`                         | `BSC_DISABLED`             | `false`                          |
| `bifrost.chains.BSC.rpc_host`                         | `BSC_HOST`                 | `https://bsc-dataseed.bnbchain.org`   |
| `bifrost.chains.BSC.username`                         | `BSC_USERNAME`             |                                  |
| `bifrost.chains.BSC.password`                         | `BSC_PASSWORD`             |                                  |
| `bifrost.chains.BSC.authorization_bearer`             | `BSC_AUTHORIZATION_BEARER` |                                  |
| `bifrost.chains.BSC.block_scanner.start_block_height` | `BSC_START_BLOCK_HEIGHT`  |                                  |

### 4.4 BASE

| Configuration Item                                     | Environment Variable        | Default                    |
|--------------------------------------------------------|-----------------------------|----------------------------|
| `bifrost.chains.BASE.disabled`                         | `BASE_DISABLED`             | `false`                    |
| `bifrost.chains.BASE.rpc_host`                         | `BASE_HOST`                 | `https://mainnet.base.org`    |
| `bifrost.chains.BASE.username`                         | `BASE_USERNAME`             |                            |
| `bifrost.chains.BASE.password`                         | `BASE_PASSWORD`             |                            |
| `bifrost.chains.BASE.authorization_bearer`             | `BASE_AUTHORIZATION_BEARER` |                            |
| `bifrost.chains.BASE.block_scanner.start_block_height` | `BASE_START_BLOCK_HEIGHT`  |                            |

### 4.5 ARB

| Configuration Item                                    | Environment Variable       | Default                                    |
|-------------------------------------------------------|----------------------------|--------------------------------------------|
| `bifrost.chains.ARB.disabled`                         | `ARB_DISABLED`             | `false`                                    |
| `bifrost.chains.ARB.rpc_host`                         | `ARB_HOST`                 | `https://arb1.arbitrum.io/rpc`     |
| `bifrost.chains.ARB.username`                         | `ARB_USERNAME`             |                                            |
| `bifrost.chains.ARB.password`                         | `ARB_PASSWORD`             |                                            |
| `bifrost.chains.ARB.authorization_bearer`             | `ARB_AUTHORIZATION_BEARER` |                                            |
| `bifrost.chains.ARB.block_scanner.start_block_height` | `ARB_START_BLOCK_HEIGHT`  |                                            |

### 4.6 DOGE

| Configuration Item                                     | Environment Variable        | Default |
|--------------------------------------------------------|-----------------------------|---------|
| `bifrost.chains.DOGE.disabled`                         | `DOGE_DISABLED`             | `false` |
| `bifrost.chains.DOGE.rpc_host`                         | `DOGE_HOST`                 |         |
| `bifrost.chains.DOGE.username`                         | `DOGE_USERNAME`             |         |
| `bifrost.chains.DOGE.password`                         | `DOGE_PASSWORD`             |         |
| `bifrost.chains.DOGE.authorization_bearer`             | `DOGE_AUTHORIZATION_BEARER` |         |
| `bifrost.chains.DOGE.block_scanner.start_block_height` | `DOGE_START_BLOCK_HEIGHT`  |         |

### 4.7 LTC

| Configuration Item                                    | Environment Variable       | Default |
|-------------------------------------------------------|----------------------------|---------|
| `bifrost.chains.LTC.disabled`                         | `LTC_DISABLED`             | `false` |
| `bifrost.chains.LTC.rpc_host`                         | `LTC_HOST`                 |         |
| `bifrost.chains.LTC.username`                         | `LTC_USERNAME`             |         |
| `bifrost.chains.LTC.password`                         | `LTC_PASSWORD`             |         |
| `bifrost.chains.LTC.authorization_bearer`             | `LTC_AUTHORIZATION_BEARER` |         |
| `bifrost.chains.LTC.block_scanner.start_block_height` | `LTC_START_BLOCK_HEIGHT`  |         |

### 4.8 BCH

| Configuration Item                                    | Environment Variable       | Default |
|-------------------------------------------------------|----------------------------|---------|
| `bifrost.chains.BCH.disabled`                         | `BCH_DISABLED`             | `false` |
| `bifrost.chains.BCH.rpc_host`                         | `BCH_HOST`                 |         |
| `bifrost.chains.BCH.username`                         | `BCH_USERNAME`             |         |
| `bifrost.chains.BCH.password`                         | `BCH_PASSWORD`             |         |
| `bifrost.chains.BCH.authorization_bearer`             | `BCH_AUTHORIZATION_BEARER` |         |
| `bifrost.chains.BCH.block_scanner.start_block_height` | `BCH_START_BLOCK_HEIGHT`  |         |

### 4.9 AVAX

| Configuration Item                                     | Environment Variable        | Default |
|--------------------------------------------------------|-----------------------------|---------|
| `bifrost.chains.AVAX.disabled`                         | `AVAX_DISABLED`             | `false` |
| `bifrost.chains.AVAX.rpc_host`                         | `AVAX_HOST`                 |         |
| `bifrost.chains.AVAX.username`                         | `AVAX_USERNAME`             |         |
| `bifrost.chains.AVAX.password`                         | `AVAX_PASSWORD`             |         |
| `bifrost.chains.AVAX.authorization_bearer`             | `AVAX_AUTHORIZATION_BEARER` |         |
| `bifrost.chains.AVAX.block_scanner.start_block_height` | `AVAX_START_BLOCK_HEIGHT`  |         |