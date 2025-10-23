package evm

import (
	_ "embed"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/mapprotocol/compass-tss/constants"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/stretchr/testify/assert"

	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
)

func getConfigForNativeTest() config.BifrostBlockScannerConfiguration {
	return config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1,
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
		Mos:                        "0xB449B28f0f3f569E90A65a1caeB612E3A9ca2051",
	}
}

func GetMetricForNativeTest() (*metrics.Metrics, error) {
	return metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.ETHChain},
	})
}

func Test_Scanner(t *testing.T) {
	cfg := getConfigForNativeTest()
	os.Setenv("KEYSTORE_PASSWORD", "123456")

	storage, err := blockscanner.NewBlockScannerStorage("./db", config.LevelDBOptions{})
	assert.Nil(t, err)

	ethClient, err := ethclient.Dial("https://bsc-prebsc-dataseed.bnbchain.org")
	assert.Nil(t, err)

	m, err := GetMetricForNativeTest()
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

	_, keyStore, err := keys.GetKeyringKeybase("/Users/zmm/Library/Ethereum/keystore/UTC--2025-09-23T07-18-09.804272000Z--69a99844d11bea5c6b73c84166e3c6b62cd870f5",
		"eth-test")
	assert.Nil(t, err)
	k := keys.NewKeysWithKeybase(kb, "test-eth", "test-eth-password", keyStore)
	bridge, err := mapo.NewBridge(bridgeCfg, m, k)

	pkeyMgr, err := pubkeymanager.NewPubKeyManager(bridge, m)
	assert.Nil(t, err)

	solvencyReporter := func(height int64) error {
		return nil
	}

	rpcClient, err := evm.NewEthRPC(
		ethClient,
		time.Second*5,
		cfg.ChainID.String(),
	)
	assert.Nil(t, err)

	scanner, err := NewEVMScanner(cfg, storage, big.NewInt(97),
		ethClient, rpcClient, bridge, m, pkeyMgr, solvencyReporter, nil)
	assert.Nil(t, err)
	txIn, err := scanner.FetchTxs(68901869, 68901869)
	assert.Nil(t, err)
	t.Log("bridgeOut", constants.EventOfBridgeOut.GetTopic().String())
	t.Log("bridgeIn", constants.EventOfBridgeIn.GetTopic().String())
	assert.Equal(t, len(txIn.TxArray), 1)
	t.Log("txIn.Chain --------------- ", txIn.Chain)
	t.Log("txIn.Count --------------- ", txIn.Count)
	t.Log("txIn.Filtered --------------- ", txIn.Filtered)
	t.Log("txIn.MemPool --------------- ", txIn.MemPool)
	t.Log("txIn.ConfirmationRequired --------------- ", txIn.ConfirmationRequired)
	t.Log("txIn.AllowFutureObservation --------------- ", txIn.AllowFutureObservation)
	for idx, ele := range txIn.TxArray {
		t.Log("txArray idx=", idx, " Height=", ele.Height)
		t.Log("txArray idx=", idx, " OrderId=", ele.OrderId)
		t.Log("txArray idx=", idx, " RefundAddr=", ele.RefundAddr)
		t.Log("txArray idx=", idx, " ChainAndGasLimit=", ele.ChainAndGasLimit)
		t.Log("txArray idx=", idx, " Vault=", ele.Vault)
		t.Log("txArray idx=", idx, " TxInType=", ele.TxOutType)
		t.Log("txArray idx=", idx, " TxHash=", ele.Tx)
		t.Log("txArray idx=", idx, " Sequence=", ele.Sequence)
		t.Log("txArray idx=", idx, " Token=", ecommon.Bytes2Hex(ele.Token))
		t.Log("txArray idx=", idx, " Amount=", ele.Amount)
		t.Log("txArray idx=", idx, " From=", ecommon.Bytes2Hex(ele.From))
		t.Log("txArray idx=", idx, " To=", ecommon.Bytes2Hex(ele.To))
		t.Log("txArray idx=", idx, " Payload=", ecommon.Bytes2Hex(ele.Payload))
		t.Log("txArray idx=", idx, " Method=", ele.Method)
		t.Log("txArray idx=", idx, " FromChain=", ele.FromChain)
	}

	//txIn, err = scanner.FetchTxs(68901869, 68901869)
	//assert.Nil(t, err)
	//assert.Equal(t, len(txIn.TxArray), 1)

	_ = storage.Close()
}
