package ethereum

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/stretchr/testify/assert"
)

func getConfigForTest() config.BifrostBlockScannerConfiguration {
	return config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1, // avoids querying eth for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second * 30,
		HTTPRequestWriteTimeout:    time.Second * 30,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		Concurrency:                1,
		GasCacheBlocks:             40,
		GasPriceResolution:         10_000_000_000,
		Mos:                        "0x297A54b40e48D1d3d8A2e91f77dE08a78a4ab10D",
	}
}

func GetMetricForTest() (*metrics.Metrics, error) {
	return metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.ETHChain},
	})
}

func Test_Scanner(t *testing.T) {
	cfg := getConfigForTest()

	storage, err := blockscanner.NewBlockScannerStorage("./db", config.LevelDBOptions{})
	assert.Nil(t, err)

	ethClient, err := ethclient.Dial("https://eth-sepolia.api.onfinality.io/public")
	assert.Nil(t, err)

	m, err := GetMetricForTest()
	assert.Nil(t, err)

	bridgeCfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "localhost",
		SignerPasswd:    "password",
		ChainHomeFolder: "./",
	}

	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := cKeys.NewInMemory(cdc)

	_, keyStore, err := keys.GetKeyringKeybase("", "eth-test")
	assert.Nil(t, err)
	k := keys.NewKeysWithKeybase(kb, "test-eth", "test-eth-password", keyStore)
	bridge, err := mapo.NewBridge(bridgeCfg, m, k)

	pkeyMgr, err := pubkeymanager.NewPubKeyManager(bridge, m)
	assert.Nil(t, err)

	solvencyReporter := func(height int64) error {
		return nil
	}

	scanner, err := NewETHScanner(cfg, storage, big.NewInt(11155111),
		ethClient, bridge, m, pkeyMgr, solvencyReporter, nil)
	assert.Nil(t, err)
	txIn, err := scanner.FetchTxs(8817015, 8817015)
	assert.Nil(t, err)
	assert.Equal(t, len(txIn.TxArray), 1)
	t.Log("txIn.Chain --------------- ", txIn.Chain)
	t.Log("txIn.Count --------------- ", txIn.Count)
	t.Log("txIn.Filtered --------------- ", txIn.Filtered)
	t.Log("txIn.MemPool --------------- ", txIn.MemPool)
	t.Log("txIn.ConfirmationRequired --------------- ", txIn.ConfirmationRequired)
	t.Log("txIn.AllowFutureObservation --------------- ", txIn.AllowFutureObservation)
	for idx, ele := range txIn.TxArray {
		t.Log("txArray idx=", idx, " TxHash=", ele.Tx)
		t.Log("txArray idx=", idx, " TxInType=", ele.TxOutType)
		t.Log("txArray idx=", idx, " Height=", ele.Height)
		t.Log("txArray idx=", idx, " Amount=", ele.Amount)
		t.Log("txArray idx=", idx, " OrderId=", ele.OrderId)
		t.Log("txArray idx=", idx, " Token=", ecommon.Bytes2Hex(ele.Token))
		t.Log("txArray idx=", idx, " Vault=", ecommon.Bytes2Hex(ele.Vault))
		t.Log("txArray idx=", idx, " To=", ecommon.Bytes2Hex(ele.To))
		t.Log("txArray idx=", idx, " Method=", ele.Method)
	}

	txIn, err = scanner.FetchTxs(8817025, 8817025)
	assert.Nil(t, err)
	assert.Equal(t, len(txIn.TxArray), 1)

	_ = storage.Close()
}
