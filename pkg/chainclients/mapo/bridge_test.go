package mapo

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	selfAbi "github.com/mapprotocol/compass-tss/pkg/abi"
	"github.com/mapprotocol/compass-tss/pkg/contract"
	"github.com/rs/zerolog/log"
)

func GetMetricForTest() (*metrics.Metrics, error) {
	return metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.ETHChain},
	})
}

func getBridgeForTest(t *testing.T) shareTypes.Bridge {
	m, err := GetMetricForTest()
	assert.Nil(t, err)

	bridgeCfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "https://testnet-rpc.maplabs.io",
		SignerPasswd:    "password",
		ChainHomeFolder: "./",
		Maintainer:      "0x0EdA5e4015448A2283662174DD7def3C3d262D38",
		ViewController:  "0x7Ea4dFBa2fA7de4C18395aCD391D9E67bECA47A6",
	}

	name := "test-eth"
	//  dont push
	priStr := os.Getenv("pri")
	kb, _, err := keys.GetKeyringKeybase(priStr, name)
	assert.Nil(t, err)

	k := keys.NewKeysWithKeybase(kb, name, "123456", priStr)
	bridge, err := NewBridge(bridgeCfg, m, k)
	assert.Nil(t, err)
	return bridge
}

func Test_Bridge_GetNetworkFee(t *testing.T) {
	bri := getBridgeForTest(t)
	size, swapSize, rate, err := bri.GetNetworkFee(common.ETHChain)
	assert.Nil(t, err)
	t.Log("ETH GAS size: ", size)
	t.Log("ETH GAS swapSize: ", swapSize)
	t.Log("ETH GAS rate: ", rate)

	exist, err := bri.HasNetworkFee(common.ETHChain)
	assert.Nil(t, err)
	assert.Equal(t, true, exist, "check eth gas")

	exist, err = bri.HasNetworkFee(common.BSCChain)
	assert.Nil(t, err)
	assert.Equal(t, true, exist, "check bsc gas")

	exist, err = bri.HasNetworkFee(common.DOGEChain)
	assert.NotNil(t, err)
	assert.Equal(t, false, exist, "check DOGE gas")
}

func Test_Bridge_PostNetworkFee(t *testing.T) {
	ethClient, err := ethclient.Dial("https://testnet-rpc.maplabs.io")
	assert.Nil(t, err)
	pkStr := os.Getenv("pri_key")
	priKey, err := ecrypto.HexToECDSA(pkStr)
	addr := ecommon.HexToAddress("0xad76db9c043fB5386D8D5C4634F55bbAda559B29")
	assert.Nil(t, err)

	ai, err := selfAbi.New(maintainerAbi)
	assert.Nil(t, err)

	to := ecommon.HexToAddress("0x0EdA5e4015448A2283662174DD7def3C3d262D38")

	input, err := ai.PackInput(constants.VoteNetworkFee,
		big.NewInt(1),
		big.NewInt(1360095883558914),
		big.NewInt(882082),
		big.NewInt(100000000), // gasPrice
		big.NewInt(1000000),   // gasLimit
		big.NewInt(1500000))   // swapGasLimit
	assert.Nil(t, err)

	head, err := ethClient.HeaderByNumber(context.Background(), nil)
	assert.Nil(t, err)

	gasFeeCap := head.BaseFee

	createdTx := ethereum.CallMsg{
		From:     addr,
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	}

	t.Log("input ", ecommon.Bytes2Hex(input))
	t.Log("gasFeeCap ", gasFeeCap)

	gasLimit, err := ethClient.EstimateGas(context.Background(), createdTx)
	assert.Nil(t, err)

	nonce, err := ethClient.NonceAt(context.Background(), addr, nil)
	assert.Nil(t, err)

	// create tx
	tipCap := new(big.Int).Mul(gasFeeCap, big.NewInt(10))
	tipCap.Div(tipCap, big.NewInt(100))
	td := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Value:     nil,
		To:        &to,
		Gas:       gasLimit,
		GasTipCap: tipCap,
		GasFeeCap: gasFeeCap,
		Data:      input,
	})

	signedTx, err := types.SignTx(td, types.NewLondonSigner(big.NewInt(212)), priKey)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(context.Background(), signedTx)
	assert.Nil(t, err)

	t.Log("postGasFee tx successfully, tx ================= ", signedTx.Hash().Hex())

}

func Test(t *testing.T) {
	TestingT(t)
}

type BridgeSuite struct {
	b *Bridge
}

var _ = Suite(&BridgeSuite{})

func (s *BridgeSuite) SetUpSuite(c *C) {
	c.Log("SetUpSuite -------- ")
} // run once on tests running

func (s *BridgeSuite) TearDownSuite(c *C) {
	c.Log("TearDownSuite -------- ")
}

func (s *BridgeSuite) SetUpTest(c *C) {
	c.Log("SetUpTest ------- ")
	m, err := GetMetricForTest()
	c.Assert(err, IsNil)

	bridgeCfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "https://testnet-rpc.maplabs.io",
		SignerPasswd:    "password",
		ChainHomeFolder: "./",
		Maintainer:      "0x0EdA5e4015448A2283662174DD7def3C3d262D38",
		ViewController:  "0x7Ea4dFBa2fA7de4C18395aCD391D9E67bECA47A6",
	}

	name := "test-eth"
	//  dont push
	priStr := os.Getenv("pri")
	kb, _, err := keys.GetKeyringKeybase(priStr, name)
	c.Assert(err, IsNil)

	k := keys.NewKeysWithKeybase(kb, name, "123456", priStr)

	// main module logger
	logger := log.With().Str("module", "mapo_client").Logger()

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	ethClient, err := ethclient.Dial(bridgeCfg.ChainHost)
	c.Assert(err, IsNil)

	chainID, err := getChainID(ethClient, time.Second*5)
	c.Assert(err, IsNil)

	priv, err := k.GetPrivateKey()
	c.Assert(err, IsNil)
	temp, err := codec.ToCmtPubKeyInterface(priv.PubKey())
	c.Assert(err, IsNil)
	pk, err := common.NewPubKeyFromCrypto(temp)
	c.Assert(err, IsNil)
	ethPrivateKey, err := evm.GetPrivateKey(priv)
	c.Assert(err, IsNil)
	mainAbi, err := newMaintainerABi()
	c.Assert(err, IsNil)
	tokenRegistry, err := NewTokenRegistry()
	c.Assert(err, IsNil)

	ai, err := selfAbi.New(maintainerAbi)
	c.Assert(err, IsNil)

	viewAi, err := selfAbi.New(viewABI)
	c.Assert(err, IsNil)

	mainCall := contract.New(ethClient, []ecommon.Address{ecommon.HexToAddress(bridgeCfg.Maintainer)}, ai)
	viewCall := contract.New(ethClient, []ecommon.Address{ecommon.HexToAddress(bridgeCfg.ViewController)}, viewAi)
	keySignWrapper, err := evm.NewKeySignWrapper(ethPrivateKey, pk, nil, chainID, "MAP")
	c.Assert(err, IsNil)

	rpcClient, err := evm.NewEthRPC(
		ethClient,
		time.Second*5,
		bridgeCfg.ChainID.String(),
	)

	c.Assert(err, IsNil)
	s.b = &Bridge{
		logger:        logger,
		cfg:           bridgeCfg,
		keys:          k,
		errCounter:    m.GetCounterVec(metrics.MapChainClientError),
		httpClient:    httpClient,
		m:             m,
		chainID:       chainID,
		broadcastLock: &sync.RWMutex{},
		ethClient:     ethClient,
		stopChan:      make(chan struct{}),
		wg:            &sync.WaitGroup{},
		ethPriKey:     ethPrivateKey,
		kw:            keySignWrapper,
		ethRpc:        rpcClient,
		mainAbi:       mainAbi,
		tokenRegistry: tokenRegistry,
		mainCall:      mainCall,
		viewCall:      viewCall,
		epoch:         big.NewInt(1),
		gasPrice:      big.NewInt(0),
	}

}

func (s *BridgeSuite) TearDownTest(c *C) {
	c.Log("TearDownTest --------- ")
}

func (s *BridgeSuite) Benchmark_GetNetworkFee(c *C) {
	for i := 0; i < c.N; i++ {
		// logic to
	}
	c.Log("Benchmark_GetNetworkFee -------- ")
}

func (s *BridgeSuite) Test_CheckOrderId(c *C) {
	var exist bool
	err := s.b.mainCall.Call(constants.IsOrderExecuted, &exist, 0,
		ecommon.HexToHash("fac1402b5656301d2b7682d05e364f90f0a3609217a26747ea22c4ad652b3512"), true)
	c.Assert(err, IsNil)
	c.Log("CheckOrderId -------- ", exist)
}
