package evm

import (
	_ "embed"
	"encoding/json"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/blockscanner"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	. "gopkg.in/check.v1"
)

const TestGasPriceResolution = 100_000_000

const Mainnet = 97

var (
	//go:embed test/deposit_evm_transaction.json
	depositEVMTx []byte
	//go:embed test/deposit_evm_receipt.json
	depositEVMReceipt []byte
	//go:embed test/transfer_out_transaction.json
	transferOutTx []byte
	//go:embed test/transfer_out_receipt.json
	transferOutReceipt []byte
	//go:embed test/deposit_tkn_transaction.json
	depositTknTx []byte
	//go:embed test/deposit_tkn_receipt.json
	depositTknReceipt []byte
	//go:embed test/block_by_number.json
	blockByNumberResp []byte
)

type BlockScannerTestSuite struct {
	m      *metrics.Metrics
	bridge shareTypes.Bridge
	keys   *keys.Keys
}

var _ = Suite(&BlockScannerTestSuite{})

func (s *BlockScannerTestSuite) SetUpTest(c *C) {
	s.m = GetMetricForTest(c)

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
	bridge, err := mapo.NewBridge(bridgeCfg, m, k)
	c.Assert(err, IsNil)
	s.bridge = bridge
}

func getConfigForTest() config.BifrostBlockScannerConfiguration {
	return config.BifrostBlockScannerConfiguration{
		ChainID:                    common.BSCChain,
		StartBlockHeight:           1, // avoids querying map for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second * 30,
		HTTPRequestWriteTimeout:    time.Second * 30,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		GasCacheBlocks:             100,
		Concurrency:                1,
		GasPriceResolution:         TestGasPriceResolution,
		TransactionBatchSize:       500,
		Mos:                        "0xCD9Fc97860755cD6630E5172dbad9454a1e029a9",
	}
}

func (s *BlockScannerTestSuite) TestNewBlockScanner(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		c.Assert(err, IsNil)
		type RPCRequest struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      interface{}     `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		var rpcRequest RPCRequest
		err = json.Unmarshal(body, &rpcRequest)
		c.Assert(err, IsNil)
		if rpcRequest.Method == "eth_chainId" {
			_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x61"}`))
			c.Assert(err, IsNil)
		}
		if rpcRequest.Method == "eth_gasPrice" {
			_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
			c.Assert(err, IsNil)
		}
	}))
	storage, err := blockscanner.NewBlockScannerStorage("", config.LevelDBOptions{})
	c.Assert(err, IsNil)
	ethClient, err := ethclient.Dial(server.URL)
	c.Assert(err, IsNil)
	rpcClient, err := evm.NewEthRPC(ethClient, time.Second, "BSC")
	c.Assert(err, IsNil)
	pubKeyManager, err := pubkeymanager.NewPubKeyManager(s.bridge, s.m)
	c.Assert(err, IsNil)
	solvencyReporter := func(height int64) error {
		return nil
	}
	bs, err := NewEVMScanner(getConfigForTest(), nil, big.NewInt(int64(Mainnet)), ethClient, rpcClient, s.bridge, s.m, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, NotNil)
	c.Assert(bs, IsNil)

	bs, err = NewEVMScanner(getConfigForTest(), storage, big.NewInt(int64(Mainnet)), ethClient, rpcClient, s.bridge, nil, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, NotNil)
	c.Assert(bs, IsNil)

	bs, err = NewEVMScanner(getConfigForTest(), storage, big.NewInt(int64(Mainnet)), nil, rpcClient, s.bridge, s.m, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, NotNil)
	c.Assert(bs, IsNil)

	bs, err = NewEVMScanner(getConfigForTest(), storage, big.NewInt(int64(Mainnet)), ethClient, rpcClient, s.bridge, s.m, nil, solvencyReporter, nil)
	c.Assert(err, NotNil)
	c.Assert(bs, IsNil)

	bs, err = NewEVMScanner(getConfigForTest(), storage, big.NewInt(int64(Mainnet)), ethClient, rpcClient, s.bridge, s.m, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, IsNil)
	c.Assert(bs, NotNil)
}

func (s *BlockScannerTestSuite) TestProcessBlock(c *C) {
	storage, err := blockscanner.NewBlockScannerStorage("./test", config.LevelDBOptions{})
	c.Assert(err, IsNil)
	ethClient, err := ethclient.Dial("https://rpc.baseservice.workers.dev/main?p=filter&chain=bsctest")
	c.Assert(err, IsNil)
	rpcClient, err := evm.NewEthRPC(ethClient, time.Second, "BSC")
	c.Assert(err, IsNil)
	pubKeyManager, err := pubkeymanager.NewPubKeyManager(s.bridge, s.m)
	c.Assert(err, IsNil)
	solvencyReporter := func(height int64) error {
		return nil
	}
	bs, err := NewEVMScanner(getConfigForTest(), storage, big.NewInt(int64(Mainnet)), ethClient, rpcClient,
		s.bridge, s.m, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, IsNil)
	txIn, err := bs.FetchTxs(59119141, 59119141)
	c.Assert(err, IsNil)
	c.Check(len(txIn.TxArray), Equals, 1)

	c.Log("txIn.Chain --------------- ", txIn.Chain)
	c.Log("txIn.Count --------------- ", txIn.Count)
	c.Log("txIn.Filtered --------------- ", txIn.Filtered)
	c.Log("txIn.MemPool --------------- ", txIn.MemPool)
	c.Log("txIn.ConfirmationRequired --------------- ", txIn.ConfirmationRequired)
	c.Log("txIn.AllowFutureObservation --------------- ", txIn.AllowFutureObservation)
	for idx, ele := range txIn.TxArray {
		c.Log("txArray idx=", idx, " TxHash=", ele.Tx)
		c.Log("txArray idx=", idx, " TxInType=", ele.TxInType.String())
		c.Log("txArray idx=", idx, " Height=", ele.Height)
		c.Log("txArray idx=", idx, " Amount=", ele.Amount)
		c.Log("txArray idx=", idx, " OrderId=", ele.OrderId)
		c.Log("txArray idx=", idx, " Token=", ecommon.Bytes2Hex(ele.Token))
		c.Log("txArray idx=", idx, " Vault=", ecommon.Bytes2Hex(ele.Vault))
		c.Log("txArray idx=", idx, " To=", ecommon.Bytes2Hex(ele.To))
		c.Log("txArray idx=", idx, " Method=", ele.Method)
	}
}

func httpTestHandler(c *C, rw http.ResponseWriter, fixture string) {
	var content []byte
	var err error

	switch fixture {
	case "500":
		rw.WriteHeader(http.StatusInternalServerError)
	default:
		content, err = os.ReadFile(fixture)
		if err != nil {
			c.Fatal(err)
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err = rw.Write(content); err != nil {
		c.Fatal(err)
	}
}

func (s *BlockScannerTestSuite) TestGetTxInItem(c *C) {
}

func (s *BlockScannerTestSuite) TestProcessReOrg(c *C) {

}

// -------------------------------------------------------------------------------------
// GasPriceV2
// -------------------------------------------------------------------------------------

func (s *BlockScannerTestSuite) TestUpdateGasPrice(c *C) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		c.Assert(err, IsNil)
		type RPCRequest struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      interface{}     `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		var rpcRequest RPCRequest
		err = json.Unmarshal(body, &rpcRequest)
		c.Assert(err, IsNil)
		if rpcRequest.Method == "eth_chainId" {
			_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x539"}`))
			c.Assert(err, IsNil)
		}
		if rpcRequest.Method == "eth_gasPrice" {
			_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
			c.Assert(err, IsNil)
		}
	}))
	storage, err := blockscanner.NewBlockScannerStorage("", config.LevelDBOptions{})
	c.Assert(err, IsNil)
	ethClient, err := ethclient.Dial(server.URL)
	c.Assert(err, IsNil)
	rpcClient, err := evm.NewEthRPC(ethClient, time.Second, "BSC")
	c.Assert(err, IsNil)
	pubKeyManager, err := pubkeymanager.NewPubKeyManager(s.bridge, s.m)
	c.Assert(err, IsNil)
	solvencyReporter := func(height int64) error {
		return nil
	}
	conf := getConfigForTest()
	bs, err := NewEVMScanner(conf, storage, big.NewInt(int64(Mainnet)), ethClient, rpcClient, s.bridge, s.m, pubKeyManager, solvencyReporter, nil)
	c.Assert(err, IsNil)
	c.Assert(bs, NotNil)

	// almost fill gas cache
	for i := 0; i < 99; i++ {
		bs.updateGasPrice([]*big.Int{
			big.NewInt(1 * TestGasPriceResolution),
			big.NewInt(2 * TestGasPriceResolution),
			big.NewInt(3 * TestGasPriceResolution),
			big.NewInt(4 * TestGasPriceResolution),
			big.NewInt(5 * TestGasPriceResolution),
		})
	}

	// empty blocks should not count
	bs.updateGasPrice([]*big.Int{})
	c.Assert(len(bs.gasCache), Equals, 99)
	c.Assert(bs.gasPrice.Cmp(big.NewInt(0)), Equals, 0)

	// now we should get the median of medians
	bs.updateGasPrice([]*big.Int{
		big.NewInt(1 * TestGasPriceResolution),
		big.NewInt(2 * TestGasPriceResolution),
		big.NewInt(3 * TestGasPriceResolution),
		big.NewInt(4 * TestGasPriceResolution),
		big.NewInt(5 * TestGasPriceResolution),
	})
	c.Assert(len(bs.gasCache), Equals, 100)
	c.Assert(bs.gasPrice.String(), Equals, big.NewInt(3*TestGasPriceResolution).String())

	// add 49 more blocks with 2x the median and we should get the same
	for i := 0; i < 49; i++ {
		bs.updateGasPrice([]*big.Int{
			big.NewInt(2 * TestGasPriceResolution),
			big.NewInt(4 * TestGasPriceResolution),
			big.NewInt(6 * TestGasPriceResolution),
			big.NewInt(8 * TestGasPriceResolution),
			big.NewInt(10 * TestGasPriceResolution),
		})
	}
	c.Assert(len(bs.gasCache), Equals, 100)
	c.Assert(bs.gasPrice.String(), Equals, big.NewInt(3*TestGasPriceResolution).String())

	// after one more block with 2x the median we should get 2x
	bs.updateGasPrice([]*big.Int{
		big.NewInt(2 * TestGasPriceResolution),
		big.NewInt(4 * TestGasPriceResolution),
		big.NewInt(6 * TestGasPriceResolution),
		big.NewInt(8 * TestGasPriceResolution),
		big.NewInt(10 * TestGasPriceResolution),
	})
	c.Assert(bs.gasPrice.String(), Equals, big.NewInt(6*TestGasPriceResolution).String())

	// add 50 more blocks with half the median and we should get the same
	for i := 0; i < 50; i++ {
		bs.updateGasPrice([]*big.Int{
			big.NewInt(TestGasPriceResolution),
		})
	}
	c.Assert(len(bs.gasCache), Equals, 100)
	c.Assert(bs.gasPrice.String(), Equals, big.NewInt(6*TestGasPriceResolution).String())

	// after one more block with half the median we should get half
	bs.updateGasPrice([]*big.Int{
		big.NewInt(TestGasPriceResolution),
	})
	c.Assert(bs.gasPrice.String(), Equals, big.NewInt(TestGasPriceResolution).String())
}
