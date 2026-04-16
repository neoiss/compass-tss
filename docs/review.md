# Compass-TSS Code Review

## 1. 项目概述

Compass-TSS (CrossX) 是 MAP Protocol 生态的**多方安全验证节点**，fork 自 THORChain Bifrost，实现 TSS (Threshold Signature Scheme) 技术，多个验证节点联合签名跨链交易，不暴露私钥，消除单点故障。

| 项目 | 信息 |
|------|------|
| Module | `github.com/mapprotocol/compass-tss` |
| Go 版本 | 1.23.4 |
| 入口 | `cmd/compass/main.go` |
| 配置 | YAML (`config/default.yaml`) + 环境变量 |
| 存储 | LevelDB |
| 日志 | zerolog |
| 监控 | Prometheus (9000), Health (6040), Cross (6041) |
| 测试 | 182 个测试文件, ~58K 行测试代码 |

### 支持链

| 类型 | 链 |
|------|-----|
| UTXO | BTC, DOGE, LTC, BCH |
| EVM | ETH, BSC, BASE, ARB, POL, XLAYER, UNI |
| 其他 | XRP, TRON |
| 中继 | MAP Relay Chain |

---

## 2. 核心架构

```
┌─────────────────────────────────────────────────────┐
│                   cmd/compass/main.go                │
├──────────┬──────────┬───────────┬───────────────────┤
│ Observer │  Signer  │ TSS Server│   P2P Network     │
├──────────┴──────────┴───────────┴───────────────────┤
│              pkg/chainclients/                       │
│  ┌────────┐ ┌─────┐ ┌──────┐ ┌─────┐ ┌──────┐      │
│  │ethereum│ │ evm │ │ utxo │ │ xrp │ │ tron │      │
│  └────────┘ └─────┘ └──────┘ └─────┘ └──────┘      │
├─────────────────────────────────────────────────────┤
│  blockscanner  │  mapclient  │  pubkeymanager       │
├─────────────────────────────────────────────────────┤
│  LevelDB Storage  │  Prometheus Metrics             │
└─────────────────────────────────────────────────────┘
```

### 关键包职责

| 包 | 职责 |
|----|------|
| `observer/` | 监听源链交易 (TxIn)，提交观察结果到 MAP |
| `signer/` | 处理出站交易 (TxOut)，通过 TSS 签名后广播到目标链 |
| `blockscanner/` | 通用区块扫描，高度追踪，持久化 |
| `tss/` | TSS 密钥管理，keysign/keygen 编排 |
| `p2p/` | libp2p 验证节点间通信 |
| `pkg/chainclients/` | 各链客户端实现 |
| `pkg/chainclients/mapo/` | MAP 中继链桥接客户端 |
| `pubkeymanager/` | 验证者公钥管理 |
| `config/` | 配置管理 |
| `internal/cross/` | 跨链交易追踪和持久化 |

---

## 3. 核心流程

### 入站流程 (Observer)

```
源链区块 → BlockScanner.FetchTxs()
    → Observer.processTxIns()
    → 确认检查 (每5秒)
    → Observer.deck() (每10秒批量发送)
    → signAndSendToMapRelay()
    → MAP Bridge: Broadcast(签名的 TxIn 观察)
```

### 出站流程 (Signer)

```
MAP Bridge: GetTxByBlockNumber() → TxOut
    → Signer.processTxnOut()
    → LevelDB 持久化
    → Pipeline.SpawnSignings()
    → [Vault/Chain 信号量 + 互斥锁]
    → chainClient.SignTx() → TSS RemoteSign()
    → chainClient.BroadcastTx()
    → MAP Bridge: Broadcast(签名确认)
```

### TSS 签名流程

```
签名请求 → KeySign.RemoteSign()
    → 任务队列 (按 PoolPubKey 分组)
    → 每秒批量处理 (最多15条/批)
    → TSS Server KeySign()
    → P2P 多方协调
    → 返回 R, S, RecoveryID
```

---

## 4. 关键依赖

| 类别 | 依赖 | 版本 |
|------|------|------|
| TSS | binance-chain/tss-lib | v0.0.0-20201118 |
| EVM | ethereum/go-ethereum (op-geth fork) | v1.14.11 |
| Cosmos | cosmos/cosmos-sdk | v0.50.11 |
| BTC | btcsuite/btcd | v0.24.2 |
| P2P | libp2p/go-libp2p | v0.11.0 |
| 存储 | syndtr/goleveldb | v1.0.1 |
| 配置 | spf13/viper | v1.19.0 |
| 监控 | prometheus/client_golang | v1.20.1 |

---

## 5. 代码质量问题

### 5.1 Critical - 未处理的 Panic

生产代码中使用 `panic` 会导致整个服务崩溃：

| 文件 | 行 | 问题 |
|------|-----|------|
| `tss/keygen.go` | 116 | `panic("tss keygen timeout")` — 应返回 error |
| `pkg/chainclients/mapo/bridge.go` | 178, 195 | `panic(err)` in GetContext() — key 获取失败直接崩溃 |
| `internal/keys/keys.go` | 79 | `panic(err)` in GetSignerInfo() |
| `common/type_convert.go` | 42, 47, 51 | `panic(fmt.Errorf(...))` — 输入格式错误导致崩溃 |
| `common/tokenlist/*.go` | init() | `panic(err)` — 初始化失败全进程崩溃 |

**建议**: 所有 panic 改为返回 error，由调用方决定处理策略。

### 5.2 High - Goroutine 泄漏风险

| 文件 | 问题 |
|------|------|
| `observer/observe.go:361` | goroutine 无 context 取消支持，graceful shutdown 时可能泄漏 |
| `cmd/compass/main.go:186-205` | health/cross server goroutine 未加入 WaitGroup 追踪 |
| `tss/keygen.go:107-110` | keygen goroutine 写 channel 失败时泄漏，fallback 是 panic |
| `p2p/communication.go` | 多个通信 goroutine 错误未传播 |

### 5.3 High - Context 误用

多处 RPC 调用使用 `context.Background()` 而非带超时的 context：

| 文件 | 问题 |
|------|------|
| `pkg/chainclients/ethereum/ethereum_block_scanner.go:387` | sem.Acquire 用 Background() |
| `pkg/chainclients/evm/evm_block_scanner.go:168, 373` | FilterLogs 用 Background() |
| `pkg/chainclients/mapo/bridge.go:259` | `context.TODO()` 在生产代码中 |

**建议**: 所有外部调用添加 `context.WithTimeout()`。

### 5.4 Medium - 废弃 API 使用

10+ 个文件使用废弃的 `io/ioutil` 包：
- `tss/go-tss/common/tss_helper_test.go`
- `tss/go-tss/keysign/notifier_test.go`
- `tss/go-tss/cmd/tss-recovery/tss-recovery.go`
- `txscript/txscript/reference_test.go`

**建议**: 替换为 `os.ReadFile` / `io.ReadAll`。

### 5.5 Medium - 已知技术债务 (TODO/FIXME)

| 文件 | 内容 |
|------|------|
| `tss/keygen.go:144` | "TODO later on need to have both secp256k1 and ed25519" |
| `config/config.go:83, 869` | "can be cleaned once all deployments are updated" |
| `constants/constant_values.go:76, 176` | "TODO: remove on hard fork" |
| `pkg/chainclients/utxo/signer_internal.go:225` | "TODO: Cleanup magic numbers" |
| `pkg/chainclients/utxo/client_internal.go:1021` | "TODO check what we do if get multiple addresses" |
| `internal/keys/keys.go:33` | "TODO this is a bad way, need to fix it" (密码字段) |

### 5.6 Low - 日志不一致

混合使用 `.Err(err)` 和 `.Error(err)`，部分错误路径日志级别不正确（用 Info 而非 Error）。`cmd/compass/main.go:156` TSS 启动失败仅 log.Err 而非 Fatal。

---

## 6. 安全关注点

| 级别 | 问题 | 位置 |
|------|------|------|
| Medium | 使用废弃的 RIPEMD160 | `pkg/chainclients/xrp/keymanager/codec.go` |
| Low | 私钥错误消息可能泄露信息 | `p2p/conversion/key_provider.go:82, 86` |
| Low | 测试代码硬编码密码 "password" | `ethereum_block_scanner_test.go`, `evm_block_scanner_test.go` |
| Info | 测试用 math/rand 而非 crypto/rand | 多个 test 文件 |

---

## 7. 并发模型

### Pipeline 并发控制 (Signer)

```
信号量: 每种 Vault 状态 (Active: N, Retiring: N, Inactive: 1)
互斥锁: 每个 vault+chain 组合只允许 1 个签名
```

这确保了确定性排序，防止竞态条件。

### Observer 并发

- 每条链独立 BlockScanner goroutine
- 全局 TxIn 队列 (channel)
- LRU 缓存 (10K) 防重复观察
- onDeck map 通过 mutex 保护

---

## 8. 数据持久化

| 存储 | 用途 |
|------|------|
| BlockScanner DB | 各链区块高度、最后处理的区块 |
| Signer DB | 待签名的 TxOut |
| Oracle DB | Oracle/Fee 交易缓存 |
| Observer DB | 已确认的 onDeck TxIn |
| Cross Storage DB | 跨链交易追踪 (orderId → 状态) |
| 内存 LRU | 已签名 TxOut 去重 (10K 条目) |

---

## 9. 改进建议汇总

| 优先级 | 建议 | 影响 |
|--------|------|------|
| P0 | 生产代码中的 panic 全部改为 error 返回 | 防止服务崩溃 |
| P0 | keygen timeout 改为 graceful error handling | TSS 稳定性 |
| P1 | 所有外部 RPC 调用添加 context timeout | 防止 goroutine 挂起 |
| P1 | goroutine 添加 context 取消和 WaitGroup 追踪 | 优雅停机 |
| P1 | bridge.go GetContext() panic 改为 error | 服务可用性 |
| P2 | 替换废弃的 ioutil 包 | 代码质量 |
| P2 | 统一日志级别和错误处理风格 | 可维护性 |
| P2 | 清理 TODO/FIXME 技术债务 | 代码健康 |
| P3 | 添加 Dockerfile | 部署便利性 |
| P3 | 添加 Makefile 到项目根目录 | 构建标准化 |
