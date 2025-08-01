package utxo

import (
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/cmd"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
	mapclient "github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/utxo"
	ttypes "github.com/mapprotocol/compass-tss/x/types"
)

type BitcoinSuite struct {
	client *Client
	server *httptest.Server
	bridge shareTypes.Bridge
	cfg    config.BifrostChainConfiguration
	m      *metrics.Metrics
	keys   *keys.Keys
}

var _ = Suite(&BitcoinSuite{})

func (s *BitcoinSuite) SetUpSuite(c *C) {
	ttypes.SetupConfigForTest()

	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := cKeys.NewInMemory(cdc)
	_, _, err := kb.NewMnemonic(bob, cKeys.English, cmd.THORChainHDPath, password, hd.Secp256k1)
	c.Assert(err, IsNil)
	//s.keys = keys.NewKeysWithKeybase(kb, bob, password, os.Getenv("TEST_PRIVATE_KEY"))
	s.keys = keys.NewKeysWithKeybase(kb, bob, password, "c3e5c914c1e15b9271de78e739c7815b5f9af4d7bc448a4ad31968c6416dba00")
}

var btcChainRPCs = map[string]map[string]interface{}{}

func init() {
	//// map the method and params to the loaded fixture
	//loadFixture := func(path string) map[string]interface{} {
	//	f, err := os.Open(path)
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer f.Close()
	//	var data map[string]interface{}
	//	err = json.NewDecoder(f).Decode(&data)
	//	if err != nil {
	//		panic(err)
	//	}
	//	return data
	//}
	//
	//btcChainRPCs["getnetworkinfo"] = loadFixture("../../../../test/fixtures/btc/getnetworkinfo.json")
	//btcChainRPCs["getblockhash"] = loadFixture("../../../../test/fixtures/btc/blockhash.json")
	//btcChainRPCs["getblock"] = loadFixture("../../../../test/fixtures/btc/block_verbose.json")
	//btcChainRPCs["getblockcount"] = loadFixture("../../../../test/fixtures/btc/blockcount.json")
	//btcChainRPCs["importaddress"] = loadFixture("../../../../test/fixtures/btc/importaddress.json")
	//btcChainRPCs["listunspent"] = loadFixture("../../../../test/fixtures/btc/listunspent.json")
	//btcChainRPCs["getrawmempool"] = loadFixture("../../../../test/fixtures/btc/getrawmempool.json")
	//btcChainRPCs["getblockstats"] = loadFixture("../../../../test/fixtures/btc/blockstats.json")
	//btcChainRPCs["getrawtransaction-5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513"] = loadFixture("../../../../test/fixtures/btc/tx-5b08.json")
	//btcChainRPCs["getrawtransaction-54ef2f4679fb90af42e8d963a5d85645d0fd86e5fe8ea4e69dbf2d444cb26528"] = loadFixture("../../../../test/fixtures/btc/tx-54ef.json")
	//btcChainRPCs["getrawtransaction-64ef2f4679fb90af42e8d963a5d85645d0fd86e5fe8ea4e69dbf2d444cb26528"] = loadFixture("../../../../test/fixtures/btc/tx-64ef.json")
	//btcChainRPCs["getrawtransaction-74ef2f4679fb90af42e8d963a5d85645d0fd86e5fe8ea4e69dbf2d444cb26528"] = loadFixture("../../../../test/fixtures/btc/tx-74ef.json")
	//btcChainRPCs["getrawtransaction-27de3e1865c098cd4fded71bae1e8236fd27ce5dce6e524a9ac5cd1a17b5c241"] = loadFixture("../../../../test/fixtures/btc/tx-c241.json")
	//btcChainRPCs["getrawtransaction"] = loadFixture("../../../../test/fixtures/btc/tx.json")
	//btcChainRPCs["createwallet"] = loadFixture("../../../../test/fixtures/btc/createwallet.json")
}

var (
	rpcHost     = "129.226.149.150:38332"
	rpcUser     = "map-signet"
	rpcPassword = "TXC6c~Vl7Ln^@PQG"
)

func (s *BitcoinSuite) SetUpTest(c *C) {
	s.m = GetMetricForTest(c, common.BTCChain)
	s.cfg = config.BifrostChainConfiguration{
		ChainID:     "BTC",
		UserName:    rpcUser,
		Password:    rpcPassword,
		DisableTLS:  true,
		HTTPostMode: true,
		BlockScanner: config.BifrostBlockScannerConfiguration{
			StartBlockHeight: 1, // avoids querying thorchain for block height
		},
		RPCHost: rpcHost,
	}
	s.cfg.UTXO.TransactionBatchSize = 500
	s.cfg.UTXO.MaxMempoolBatches = 10
	s.cfg.UTXO.EstimatedAverageTxSize = 1000
	s.cfg.BlockScanner.MaxReorgRescanBlocks = 1
	ns := strconv.Itoa(time.Now().Nanosecond())

	thordir := filepath.Join(os.TempDir(), ns, ".thorcli")
	cfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "https://testnet-rpc.maplabs.io",
		ChainRPC:        "https://testnet-rpc.maplabs.io",
		SignerName:      bob,
		SignerPasswd:    password,
		ChainHomeFolder: thordir,
	}

	handleRPC := func(body []byte, rw http.ResponseWriter) {
		r := struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}{}

		err := json.Unmarshal(body, &r)
		c.Assert(err, IsNil)

		rw.Header().Set("Content-Type", "application/json")
		key := r.Method
		if r.Method == "getrawtransaction" {
			key = fmt.Sprintf("%s-%s", r.Method, r.Params[0])
		}
		if btcChainRPCs[key] == nil {
			key = r.Method
		}

		err = json.NewEncoder(rw).Encode(btcChainRPCs[key])
		c.Assert(err, IsNil)
	}
	handleBatchRPC := func(body []byte, rw http.ResponseWriter) {
		r := []struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
			ID     int           `json:"id"`
		}{}

		err := json.Unmarshal(body, &r)
		c.Assert(err, IsNil)

		rw.Header().Set("Content-Type", "application/json")
		result := make([]map[string]interface{}, len(r))
		for i, v := range r {
			key := v.Method
			if v.Method == "getrawtransaction" {
				key = fmt.Sprintf("%s-%s", v.Method, v.Params[0])
			}
			if btcChainRPCs[key] == nil {
				key = v.Method
			}
			result[i] = btcChainRPCs[key]
			result[i]["id"] = v.ID
		}

		err = json.NewEncoder(rw).Encode(result)
		c.Assert(err, IsNil)
	}

	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/" { // nolint
			body, _ := io.ReadAll(req.Body)
			if body[0] == '[' {
				handleBatchRPC(body, rw)
			} else {
				handleRPC(body, rw)
			}
		} else if strings.HasPrefix(req.RequestURI, "/thorchain/node/") {
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/nodeaccount/template.json")
		} else if req.RequestURI == "/thorchain/lastblock" {
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/lastblock/btc.json")
		} else if strings.HasPrefix(req.RequestURI, "/auth/accounts/") {
			_, err := rw.Write([]byte(`{ "jsonrpc": "2.0", "id": "", "result": { "height": "0", "result": { "value": { "account_number": "0", "sequence": "0" } } } }`))
			c.Assert(err, IsNil)
		} else if req.RequestURI == "/txs" {
			_, err := rw.Write([]byte(`{"height": "1", "txhash": "AAAA000000000000000000000000000000000000000000000000000000000000", "logs": [{"success": "true", "log": ""}]}`))
			c.Assert(err, IsNil)
		} else if strings.HasPrefix(req.RequestURI, mapclient.AsgardVault) {
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/vaults/asgard.json")
		} else if req.RequestURI == "/thorchain/mimir/key/MaxUTXOsToSpend" {
			_, err := rw.Write([]byte(`-1`))
			c.Assert(err, IsNil)
		} else if req.RequestURI == "/thorchain/vaults/pubkeys" {
			if common.CurrentChainNetwork == common.MainNet {
				httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/vaults/pubKeys-Mainnet.json")
			} else {
				httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/vaults/pubKeys.json")
			}
		}
	}))
	var err error
	//cfg.ChainHost = s.server.Listener.Addr().String()
	s.bridge, err = mapclient.NewBridge(cfg, s.m, s.keys)
	c.Assert(err, IsNil)
	//s.cfg.RPCHost = s.server.Listener.Addr().String()
	s.client, err = NewClient(s.keys, s.cfg, nil, s.bridge, s.m)
	s.client.disableVinZeroBatch = true
	s.client.globalNetworkFeeQueue = make(chan types.NetworkFee, 1)
	c.Assert(err, IsNil)
	c.Assert(s.client, NotNil)
}

func (s *BitcoinSuite) TearDownTest(_ *C) {
	s.server.Close()
}

func (s *BitcoinSuite) TestGetBlock(c *C) {
	block, err := s.client.getBlock(1696761)
	c.Assert(err, IsNil)
	c.Assert(block.Hash, Equals, "000000008de7a25f64f9780b6c894016d2c63716a89f7c9e704ebb7e8377a0c8")
	c.Assert(block.Tx[0].Txid, Equals, "31f8699ce9028e9cd37f8a6d58a79e614a96e3fdd0f58be5fc36d2d95484716f")
	c.Assert(len(block.Tx), Equals, 112)
}

func (s *BitcoinSuite) TestFetchTxs(c *C) {
	//var vaultPubKey common.PubKey
	var err error
	//if common.CurrentChainNetwork == common.MainNet {
	//	vaultPubKey, err = common.NewPubKey("thorpub1addwnpepqwprh5vd0rrk78kd98qjruuazwvapnxft7f86w7hlf768whxytpn5quf2gs") // from PubKeys-Mainnet.json
	//} else {
	//	vaultPubKey, err = common.NewPubKey("tthorpub1addwnpepqflvfv08t6qt95lmttd6wpf3ss8wx63e9vf6fvyuj2yy6nnyna576rfzjks") // from PubKeys.json
	//}
	//c.Assert(err, IsNil, Commentf(vaultPubKey.String()))
	//vaultAddress, err := vaultPubKey.GetAddress(s.client.GetChain())
	c.Assert(err, IsNil)
	//vaultAddressString := vaultAddress.String()

	txs, err := s.client.FetchTxs(261927, 261928)
	c.Assert(err, IsNil)
	c.Assert(txs.Chain, Equals, common.BTCChain)
	c.Assert(txs.TxArray[0].Height, Equals, int64(1696761))
	c.Assert(txs.TxArray[0].Tx, Equals, "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2")
	c.Assert(txs.TxArray[0].Sender, Equals, "tb1qdxxlx4r4jk63cve3rjpj428m26xcukjn5yegff")
	//c.Assert(txs.TxArray[0].To, Equals, vaultAddressString)
	c.Assert(len(txs.TxArray), Equals, 1)
}

func (s *BitcoinSuite) TestGetSender(c *C) {
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "31f8699ce9028e9cd37f8a6d58a79e614a96e3fdd0f58be5fc36d2d95484716f",
				Vout: 0,
			},
		},
	}
	sender, err := s.client.getSender(&tx, nil)
	c.Assert(err, IsNil)
	c.Assert(sender, Equals, "n3jYBjCzgGNydQwf83Hz6GBzGBhMkKfgL1")

	tx.Vin[0].Vout = 1
	sender, err = s.client.getSender(&tx, nil)
	c.Assert(err, IsNil)
	c.Assert(sender, Equals, "tb1qdxxlx4r4jk63cve3rjpj428m26xcukjn5yegff")
}

func (s *BitcoinSuite) TestGetMemo(c *C) {
	tx := btcjson.TxRawResult{
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Hex:       "6a1574686f72636861696e3a636f6e736f6c6964617465",
					ReqSigs:   0,
					Type:      "nulldata",
					Addresses: nil,
				},
			},
		},
	}
	memo, err := s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "thorchain:consolidate")

	tx = btcjson.TxRawResult{
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 737761703a6574682e3078633534633135313236393646334541373935366264396144343130383138654563414443466666663a30786335346331353132363936463345413739353662643961443431",
					Type: "nulldata",
					Hex:  "6a4c50737761703a6574682e3078633534633135313236393646334541373935366264396144343130383138654563414443466666663a30786335346331353132363936463345413739353662643961443431",
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 30383138654563414443466666663a3130303030303030303030",
					Type: "nulldata",
					Hex:  "6a1a30383138654563414443466666663a3130303030303030303030",
				},
			},
		},
	}
	memo, err = s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "swap:eth.0xc54c1512696F3EA7956bd9aD410818eEcADCFfff:0xc54c1512696F3EA7956bd9aD410818eEcADCFfff:10000000000")

	tx = btcjson.TxRawResult{
		Vout: []btcjson.Vout{},
	}
	memo, err = s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "")
}

func (s *BitcoinSuite) TestIgnoreTx(c *C) {
	var currentHeight int64 = 100

	// valid tx that will NOT be ignored
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.12345678,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
					Type:      "nulldata",
				},
			},
		},
	}
	ignored := s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, false)

	// tx with LockTime later than current height, so should be ignored
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.12345678,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
					Type:      "nulldata",
				},
			},
		},
		LockTime: uint32(currentHeight) + 1,
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// tx with LockTime equal to current height, so should not be ignored
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.12345678,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
					Type:      "nulldata",
				},
			},
		},
		LockTime: uint32(currentHeight),
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, false)

	// invalid tx missing Vout
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vout[0].Value == no coins
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vin[0].Txid means coinbase
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vin
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// invalid tx > 10 vout with coins we only expect 10 max
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, true)

	// valid tx == 2 vout with coins, 1 to vault, 1 with change back to user
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, false)

	// memo at first output should not ignore
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, false)

	// memo in the middle , should not ignore
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
		},
	}
	ignored = s.client.ignoreTx(&tx, currentHeight)
	c.Assert(ignored, Equals, false)
}

func (s *BitcoinSuite) TestGetGas(c *C) {
	// vin[0] returns value 0.19590108
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.12345678,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm: "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
				},
			},
		},
	}
	gas, err := s.client.getGas(&tx)
	c.Assert(err, IsNil)
	c.Assert(gas.Equals(common.Gas{common.NewCoin(common.BTCAsset, cosmos.NewUint(7244430))}), Equals, true)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm: "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
				},
			},
		},
	}
	gas, err = s.client.getGas(&tx)
	c.Assert(err, IsNil)
	c.Assert(gas.Equals(common.Gas{common.NewCoin(common.BTCAsset, cosmos.NewUint(149013))}), Equals, true)
}

func (s *BitcoinSuite) TestGetChain(c *C) {
	chain := s.client.GetChain()
	c.Assert(chain, Equals, common.BTCChain)
}

func (s *BitcoinSuite) TestGetHeight(c *C) {
	height, err := s.client.GetHeight()
	c.Assert(err, IsNil)
	c.Assert(height, Equals, int64(10))
}

func (s *BitcoinSuite) TestOnObservedTxIn(c *C) {
	pkey := ttypes.GetRandomPubKey()
	txIn := types.TxIn{
		Chain: common.BTCChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(1),
				Tx:     "31f8699ce9028e9cd37f8a6d58a79e614a96e3fdd0f58be5fc36d2d95484716f",
				Sender: "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				//To:          "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta := utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 1, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	s.client.OnObservedTxIn(*txIn.TxArray[0], 1)
	blockMeta, err := s.client.temporalStorage.GetBlockMeta(1)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)

	txIn = types.TxIn{
		Chain: common.BTCChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(2),
				Tx:     "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta = utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 2, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	s.client.OnObservedTxIn(*txIn.TxArray[0], 2)
	blockMeta, err = s.client.temporalStorage.GetBlockMeta(2)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)

	txIn = types.TxIn{
		Chain: common.BTCChain,
		TxArray: []*types.TxInItem{
			{
				Height: big.NewInt(3),
				Tx:     "44ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
			{
				Height: big.NewInt(3),
				Tx:     "54ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender: "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				//To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta = utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 3, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	for _, item := range txIn.TxArray {
		s.client.OnObservedTxIn(*item, 3)
	}

	blockMeta, err = s.client.temporalStorage.GetBlockMeta(3)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)
}

func (s *BitcoinSuite) TestProcessReOrg(c *C) {
	// can't get previous block meta should not error
	type response struct {
		Result btcjson.GetBlockVerboseResult `json:"result"`
	}
	res := response{}
	blockContent, err := os.ReadFile("../../../../test/fixtures/btc/block.json")
	c.Assert(err, IsNil)
	c.Assert(json.Unmarshal(blockContent, &res), IsNil)
	result := btcjson.GetBlockVerboseTxResult{
		Hash:         res.Result.Hash,
		PreviousHash: res.Result.PreviousHash,
		Height:       res.Result.Height,
	}
	// should not trigger re-org process
	reorgedItems, err := s.client.processReorg(&result)
	c.Assert(err, IsNil)
	c.Assert(reorgedItems, IsNil)

	// add one UTXO which will trigger the re-org process next
	previousHeight := result.Height - 1
	blockMeta := utxo.NewBlockMeta(ttypes.GetRandomTxHash().String(), previousHeight, ttypes.GetRandomTxHash().String())
	hash := "27de3e1865c098cd4fded71bae1e8236fd27ce5dce6e524a9ac5cd1a17b5c241"
	blockMeta.AddCustomerTransaction(hash)
	c.Assert(s.client.temporalStorage.SaveBlockMeta(previousHeight, blockMeta), IsNil)
	s.client.globalErrataQueue = make(chan types.ErrataBlock, 1)
	reorgedItems, err = s.client.processReorg(&result)
	c.Assert(err, IsNil)
	c.Assert(reorgedItems, NotNil)
	// make sure there is errata block in the queue
	c.Assert(s.client.globalErrataQueue, HasLen, 1)
	blockMeta, err = s.client.temporalStorage.GetBlockMeta(previousHeight)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)
}

func (s *BitcoinSuite) TestGetMemPool(c *C) {
	txIns, err := s.client.FetchMemPool(1024)
	c.Assert(err, IsNil)
	c.Assert(txIns.TxArray, HasLen, 1)

	// process it again , the tx will be ignored
	txIns, err = s.client.FetchMemPool(1024)
	c.Assert(err, IsNil)
	c.Assert(txIns.TxArray, HasLen, 0)
}

func (s *BitcoinSuite) TestGetOutput(c *C) {
	var vaultPubKey common.PubKey
	var err error
	if common.CurrentChainNetwork == common.MainNet {
		vaultPubKey, err = common.NewPubKey("thorpub1addwnpepqwprh5vd0rrk78kd98qjruuazwvapnxft7f86w7hlf768whxytpn5quf2gs") // from PubKeys-Mainnet.json
	} else {
		vaultPubKey, err = common.NewPubKey("tthorpub1addwnpepqflvfv08t6qt95lmttd6wpf3ss8wx63e9vf6fvyuj2yy6nnyna576rfzjks") // from PubKeys.json
	}
	c.Assert(err, IsNil, Commentf(vaultPubKey.String()))
	vaultAddress, err := vaultPubKey.GetAddress(s.client.GetChain())
	c.Assert(err, IsNil)
	vaultAddressString := vaultAddress.String()

	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err := s.client.getOutput(vaultAddressString, &tx, false)
	c.Assert(err, IsNil, Commentf(vaultAddressString))
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
		},
	}
	out, err = s.client.getOutput(vaultAddressString, &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
		},
	}
	out, err = s.client.getOutput(vaultAddressString, &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err = s.client.getOutput(vaultAddressString, &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{vaultAddressString},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err = s.client.getOutput(vaultAddressString, &tx, true)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, vaultAddressString)
	c.Assert(out.Value, Equals, 1.49655603)

	// invalid tx only multiple (positive-value) vout Addresses
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err = s.client.getOutput(vaultAddressString, &tx, true)
	c.Assert(err, NotNil)
}

func (s *BitcoinSuite) TestIsValidUTXO(c *C) {
	// normal pay to pubkey hash segwit
	c.Assert(s.client.isValidUTXO("00140653096f54ae1ae2d73291d15854aef08ebcfa8c"), Equals, true)
	// pubkey hash , bitcoin client doesn't use it
	c.Assert(s.client.isValidUTXO("76a91415fb126815935f6ae83a206d7d82f1065bc63e2588ac"), Equals, true)

	c.Assert(s.client.isValidUTXO("a914e51a3dd98ded55718ad2cf2ce7c8ff056394445787"), Equals, true)
	c.Assert(s.client.isValidUTXO("00483045022100995187373cabc9ef02e5dd2770519704054bff6e3b42f8eeeb1f08a40db527b50220380d4d1f471087c35ebdde4251a0c8fa38db688600020a189d2b19343d079c100147304402205c8a886fece4c40c96c47ee51cf8e32ff75251375a47c4ec0cec9193a8a747620220703045742cdec7a19e16aa071a4fe4333a6c1b587783b864d0a64ef87783b13a014c695221039b3fa7e3dd5f9caab777f0dd15a03f1011063a2bf205f96ad2b01540506109432103e7e00ea57b70cfd9493f1d7e482a2bfe4c785d8e9bef25eb1fd3a528bc452e072103f01a388aecf967af2d21a8578635e745a3990afdfde8099ac44bab3ecd9c042153ae"), Equals, false)
	c.Assert(s.client.isValidUTXO("51210281feb90c058c3436f8bc361930ae99fcfb530a699cdad141d7244bfcad521a1f51ae"), Equals, false)
	c.Assert(s.client.isValidUTXO("5121037953dbf08030f67352134992643d033417eaa6fcfb770c038f364ff40d7615882100bd2fda4cf456d64386a0756f580101a607c25bd8d6814693bdf16e2a7ba3e45c52ae"), Equals, false)
	c.Assert(s.client.isValidUTXO("524104d81fd577272bbe73308c93009eec5dc9fc319fc1ee2e7066e17220a5d47a18314578be2faea34b9f1f8ca078f8621acd4bc22897b03daa422b9bf56646b342a24104ec3afff0b2b66e8152e9018fe3be3fc92b30bf886b3487a525997d00fd9da2d012dce5d5275854adc3106572a5d1e12d4211b228429f5a7b2f7ba92eb0475bb14104b49b496684b02855bc32f5daefa2e2e406db4418f3b86bca5195600951c7d918cdbe5e6d3736ec2abf2dd7610995c3086976b2c0c7b4e459d10b34a316d5a5e753ae"), Equals, false)

	// V1_P2TR (pay-to-taproot) output
	c.Assert(s.client.isValidUTXO("5120f01002397e3cb9179d41f1e25412bd29fc8d22f8fe786758aeeacf137a4cbc5f"), Equals, true)
}

func TestFetchTxs(t *testing.T) {
	signerName := "first_cosmos"
	signerPasswd := os.Getenv("SIGNER_PASSWD")
	privateKey := "c3e5c914c1e15b9271de78e739c7815b5f9af4d7bc448a4ad31968c6416dba00"

	scfg := config.BifrostChainConfiguration{
		ChainID:     "BTC",
		UserName:    rpcUser,
		Password:    rpcPassword,
		DisableTLS:  true,
		HTTPostMode: true,
		BlockScanner: config.BifrostBlockScannerConfiguration{
			StartBlockHeight: 1, // avoids querying thorchain for block height
		},
		RPCHost: rpcHost,
	}
	scfg.UTXO.TransactionBatchSize = 500
	scfg.UTXO.MaxMempoolBatches = 10
	scfg.UTXO.EstimatedAverageTxSize = 1000
	scfg.BlockScanner.MaxReorgRescanBlocks = 1
	ns := strconv.Itoa(time.Now().Nanosecond())

	thordir := filepath.Join(os.TempDir(), ns, ".thorcli")
	cfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "https://testnet-rpc.maplabs.io",
		ChainRPC:        "https://testnet-rpc.maplabs.io",
		SignerName:      bob,
		SignerPasswd:    password,
		ChainHomeFolder: thordir,
	}

	kb, _, err := keys.GetKeyringKeybase(privateKey, signerName)
	if err != nil {
		t.Fatal(err)
	}

	ks := keys.NewKeysWithKeybase(kb, signerName, signerPasswd, privateKey)
	bridge, err := mapclient.NewBridge(cfg, m, ks)
	if err != nil {
		t.Fatal(err)
	}

	m, err = metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.DOGEChain, common.BCHChain, common.LTCChain, common.BTCChain},
	})
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(ks, scfg, nil, bridge, m)
	if err != nil {
		t.Fatal(err)
	}
	client.globalNetworkFeeQueue = make(chan types.NetworkFee, 1)

	txs, err := client.FetchTxs(261927, 261928)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txs)
}

func TestHex(t *testing.T) {
	t.Log(ethcommon.Hex2Bytes("0eb16a9cfdf8e3a4471ef190ee63de5a24f38787"))
	t.Log(ethcommon.Hex2Bytes("0x0eb16a9cfdf8e3a4471ef190ee63de5a24f38787"))
	t.Log(ethcommon.ParseHexOrString("0x0eb16a9cfdf8e3a4471ef190ee63de5a24f38787"))
}
