package evm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	btypes "github.com/mapprotocol/compass-tss/blockscanner/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// EthRPC is a struct that interacts with an ETH RPC compatible blockchain
type EthRPC struct {
	client  *ethclient.Client
	timeout time.Duration
	logger  zerolog.Logger
}

func NewEthRPC(client *ethclient.Client, timeout time.Duration, chain string) (*EthRPC, error) {
	return &EthRPC{
		client:  client,
		timeout: timeout,
		logger:  log.Logger.With().Str("module", "eth_rpc").Str("chain", chain).Logger(),
	}, nil
}

func (e *EthRPC) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), e.timeout)
}

func (e *EthRPC) EstimateGas(from string, tx *etypes.Transaction) (uint64, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.EstimateGas(ctx, ethereum.CallMsg{
		From:     ecommon.HexToAddress(from),
		To:       tx.To(),
		GasPrice: tx.GasPrice(),
		// Gas:      tx.Gas(),
		Value: tx.Value(),
		Data:  tx.Data(),
	})
}

func (e *EthRPC) GetReceipt(hash string) (*etypes.Receipt, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.TransactionReceipt(ctx, ecommon.HexToHash(hash))
}

func (e *EthRPC) GetHeader(height int64) (*etypes.Header, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.HeaderByNumber(ctx, big.NewInt(height))
}

func (e *EthRPC) GetBlockHeight() (int64, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	height, err := e.client.BlockNumber(ctx)
	if err != nil {
		e.logger.Info().Err(err).Msg("failed to get block height")
		return -1, fmt.Errorf("fail to get block height: %w", err)
	}
	return int64(height), nil
}

func (e *EthRPC) GetBlockHeightSafe() (int64, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	block, err := e.client.BlockByNumber(ctx, big.NewInt(rpc.SafeBlockNumber.Int64()))
	if err != nil {
		e.logger.Info().Err(err).Msg("failed to get block")
		return -1, fmt.Errorf("fail to get block: %w", err)
	}
	return block.Number().Int64(), nil
}

func (e *EthRPC) GetBlock(height int64) (*etypes.Block, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	return e.client.BlockByNumber(ctx, big.NewInt(height))
}

func (e *EthRPC) GetRPCBlock(height int64) (*etypes.Block, error) {
	block, err := e.GetBlock(height)
	if err == ethereum.NotFound {
		return nil, btypes.ErrUnavailableBlock
	}
	if err != nil {
		return nil, fmt.Errorf("fail to fetch block: %w", err)
	}
	return block, nil
}

// GetNonce gets nonce (including pending) of an address.
func (e *EthRPC) GetNonce(addr string) (uint64, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	nonce, err := e.client.PendingNonceAt(ctx, ecommon.HexToAddress(addr))
	if err != nil {
		return 0, fmt.Errorf("fail to get account nonce: %w", err)
	}
	return nonce, nil
}

// GetNonceFinalized gets the nonce excluding pending transactions.
func (e *EthRPC) GetNonceFinalized(addr string) (uint64, error) {
	ctx, cancel := e.getContext()
	defer cancel()
	nonce, err := e.client.NonceAt(ctx, ecommon.HexToAddress(addr), nil)
	if err != nil {
		return 0, fmt.Errorf("fail to get account nonce: %w", err)
	}
	return nonce, nil
}

// CheckTransaction returns true if a transaction is found and successful on chain. This
// is used to determine when a transaction has been dropped from the chain or failed on
// subsequent execution after reorgs. This can return false positives, but should not
// return false negatives - as we want to errata an observation only when we are certain
// it has later been dropped or failed.
func (e *EthRPC) CheckTransaction(hash string) bool {
	ctx, cancel := e.getContext()
	defer cancel()
	tx, pending, err := e.client.TransactionByHash(ctx, ecommon.HexToHash(hash))
	if err != nil || tx == nil {
		e.logger.Info().Str("txid", hash).Err(err).Msg("tx not found")
		return false
	}

	// pending transactions may fail, but we should only errata when there is certainty
	if pending {
		e.logger.Warn().Str("txid", hash).Msg("observed transaction is pending")
		return true // unknown, prefer false positive
	}

	// ensure the tx was successful
	receipt, err := e.GetReceipt(hash)
	if err != nil {
		e.logger.Warn().Str("txid", hash).Err(err).Msg("tx receipt not found")
		return true // unknown, prefer false positive
	}
	return receipt.Status == etypes.ReceiptStatusSuccessful
}

func (e *EthRPC) GetBlockSafe(number int64) (*Block, error) {
	ctx, cancel := e.getContext()
	defer cancel()

	block, err := e.getSafeBlock(ctx, "eth_getBlockByNumber", toBlockNumArg(big.NewInt(number)), true)
	if err != nil || block == nil {
		e.logger.Info().Any("number", number).Err(err).Msg("block not found")
		return nil, err
	}

	return block, nil
}

func (e *EthRPC) getSafeBlock(ctx context.Context, method string, args ...interface{}) (*Block, error) {
	var raw json.RawMessage
	err := e.client.Client().CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, ethereum.NotFound
	}
	var head *etypes.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}

	// Decode header and transactions.
	var body Block
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	body.Header = head
	if body.Header.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if body.Header.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if body.Header.TxHash == types.EmptyRootHash && len(body.Transactions) > 0 {
		return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if body.Header.TxHash != types.EmptyRootHash && len(body.Transactions) == 0 {
		return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}
	return &body, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}

type Transaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	ChainID              string        `json:"chainId"`
	GasFeeCap            string        `json:"gasFeeCap,omitempty"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
	AccessList           []interface{} `json:"accessList,omitempty"`
	YParity              string        `json:"yParity,omitempty"`
}

func (t *Transaction) GetGasPrice() (*big.Int, bool) {
	switch strings.TrimPrefix(t.Type, "0x") {
	case "0", "1", "120": // ArbitrumLegacyTx
		return big.NewInt(0).SetString(strings.TrimPrefix(t.GasPrice, "0x"), 16)
	case "2", "3", "4": // DynamicFeeTxType、 BlobTxType、 SetCodeTxType
		return big.NewInt(0).SetString(strings.TrimPrefix(t.MaxFeePerGas, "0x"), 16)
	case "101", "102", "104", "105": // ArbitrumUnsignedTxType、ArbitrumContractTxType、ArbitrumRetryTxType、ArbitrumSubmitRetryableTxType
		return big.NewInt(0).SetString(strings.TrimPrefix(t.MaxFeePerGas, "0x"), 16)
	case "126", "100", "106": // DepositTxType、ArbitrumDepositTxType、ArbitrumInternalTxType
		return big.NewInt(0), true
	}
	return big.NewInt(-1), false
}

type Block struct {
	Header       *etypes.Header
	Transactions []Transaction  `json:"transactions"`
	UncleHashes  []ecommon.Hash `json:"uncles"`
}

func (b *Block) GetHeader() *etypes.Header {
	return b.Header
}
