package mapo

import (
	"context"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	selfAbi "github.com/mapprotocol/compass-tss/pkg/abi"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
	"time"
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

	//registry := codectypes.NewInterfaceRegistry()
	//cryptocodec.RegisterInterfaces(registry)
	//cdc := codec.NewProtoCodec(registry)
	//kb := cKeys.NewInMemory(cdc)

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
