package blockscanner

import (
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/mapclient"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/mapprotocol/compass-tss/metrics"
)

func TestPackage(t *testing.T) { TestingT(t) }

var m *metrics.Metrics

type BlockScannerTestSuite struct {
	m      *metrics.Metrics
	bridge shareTypes.ThorchainBridge
	cfg    config.BifrostClientConfiguration
	keys   *mapclient.Keys
}

var _ = Suite(&BlockScannerTestSuite{})

func (s *BlockScannerTestSuite) TearDownSuite(c *C) {
}

func (s *BlockScannerTestSuite) TestNewBlockScanner(c *C) {
	mss := NewMockScannerStorage()
	cbs, err := NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight: 1, // avoids querying thorchain for block height
	}, mss, nil, nil, DummyFetcher{})
	c.Check(cbs, IsNil)
	c.Check(err, NotNil)
	cbs, err = NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight: 1, // avoids querying thorchain for block height
	}, mss, nil, nil, DummyFetcher{})
	c.Check(cbs, IsNil)
	c.Check(err, NotNil)
	cbs, err = NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight: 1, // avoids querying thorchain for block height
	}, mss, m, s.bridge, DummyFetcher{})
	c.Check(cbs, NotNil)
	c.Check(err, IsNil)
}

const (
	blockBadResult  = `{ "jsonrpc": "2.0", "id": "", "result": { "block_meta": { "block_id": { "hash": "D063E5F1562F93D46FD4F01CA24813DD60B919D1C39CC34EF1DBB0EA07D0F7F8"1EB49C7042E5622189EDD4FA" } } } }`
	lastBlockResult = `[ { "chain": "ETH", "last_observed_in": 1, "last_signed_out": 1, "thorchain": 3 }]`
)

func (s *BlockScannerTestSuite) TestBlockScanner(c *C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.RequestURI, mapclient.MimirEndpoint):
			buf, err := os.ReadFile("../../test/fixtures/endpoints/mimir/mimir.json")
			c.Assert(err, IsNil)
			_, err = w.Write(buf)
			c.Assert(err, IsNil)
		case strings.HasPrefix(r.RequestURI, "/thorchain/lastblock"):
			// NOTE: weird pattern in GetBlockHeight uses first thorchain height.
			_, err := w.Write([]byte(`[
          {
            "chain": "NOOP",
            "lastobservedin": 0,
            "lastsignedout": 0,
            "thorchain": 0
          }
        ]`))
			c.Assert(err, IsNil)
		}
	})
	mss := NewMockScannerStorage()
	server := httptest.NewServer(h)
	defer server.Close()
	bridge, err := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       server.Listener.Addr().String(),
		ChainRPC:        server.Listener.Addr().String(),
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}, s.m, s.keys)
	c.Assert(err, IsNil)

	cbs, err := NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1, // avoids querying thorchain for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second * 30,
		HTTPRequestWriteTimeout:    time.Second * 30,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		ChainID:                    common.ETHChain,
	}, mss, m, bridge, DummyFetcher{})
	c.Check(cbs, NotNil)
	c.Check(err, IsNil)
	var counter int
	go func() {
		for item := range cbs.GetMessages() {
			_ = item
			counter++
		}
	}()
	globalChan := make(chan types.TxIn)
	nfChan := make(chan common.NetworkFee)
	cbs.Start(globalChan, nfChan)
	time.Sleep(time.Second * 1)
	cbs.Stop()
}

func (s *BlockScannerTestSuite) TestBadBlock(c *C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logf("================>:%s", r.RequestURI)
		switch {
		case strings.HasPrefix(r.RequestURI, mapclient.MimirEndpoint):
			buf, err := os.ReadFile("../../test/fixtures/endpoints/mimir/mimir.json")
			c.Assert(err, IsNil)
			_, err = w.Write(buf)
			c.Assert(err, IsNil)
		case strings.HasPrefix(r.RequestURI, "/block"): // trying to get block
			if _, err := w.Write([]byte(blockBadResult)); err != nil {
				c.Error(err)
			}
		}
	})
	mss := NewMockScannerStorage()
	server := httptest.NewTLSServer(h)
	defer server.Close()
	bridge, err := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       server.Listener.Addr().String(),
		ChainRPC:        server.Listener.Addr().String(),
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}, s.m, s.keys)
	c.Assert(err, IsNil)
	cbs, err := NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1, // avoids querying thorchain for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second * 30,
		HTTPRequestWriteTimeout:    time.Second * 30,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		ChainID:                    common.ETHChain,
	}, mss, m, bridge, DummyFetcher{})
	c.Check(cbs, NotNil)
	c.Check(err, IsNil)
	cbs.Start(make(chan types.TxIn), make(chan common.NetworkFee))
	time.Sleep(time.Second * 1)
	cbs.Stop()
}

func (s *BlockScannerTestSuite) TestBadConnection(c *C) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, mapclient.MimirEndpoint) {
			buf, err := os.ReadFile("../../test/fixtures/endpoints/mimir/mimir.json")
			c.Assert(err, IsNil)
			_, err = w.Write(buf)
			c.Assert(err, IsNil)
		}
	})
	mss := NewMockScannerStorage()
	server := httptest.NewServer(h)
	defer server.Close()
	bridge, err := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       server.Listener.Addr().String(),
		ChainRPC:        server.Listener.Addr().String(),
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}, s.m, s.keys)
	c.Assert(err, IsNil)

	cbs, err := NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1, // avoids querying thorchain for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second,
		HTTPRequestWriteTimeout:    time.Second,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		ChainID:                    common.ETHChain,
	}, mss, m, bridge, DummyFetcher{})
	c.Check(cbs, NotNil)
	c.Check(err, IsNil)
	cbs.Start(make(chan types.TxIn), make(chan common.NetworkFee))
	time.Sleep(time.Second * 1)
	cbs.Stop()
}

func (s *BlockScannerTestSuite) TestIsChainPaused(c *C) {
	mimirMap := map[string]int{
		"HaltETHChain":         0,
		"SolvencyHaltETHChain": 0,
		"HaltChainGlobal":      0,
		"NodePauseChainGlobal": 0,
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logf("================>:%s", r.RequestURI)
		switch {
		case strings.HasPrefix(r.RequestURI, mapclient.LastBlockEndpoint):
			if _, err := w.Write([]byte(lastBlockResult)); err != nil {
				c.Error(err)
			}
		case strings.HasPrefix(r.RequestURI, mapclient.MimirEndpoint):
			parts := strings.Split(r.RequestURI, "/key/")
			mimirKey := parts[1]

			mimirValue := 0
			if val, found := mimirMap[mimirKey]; found {
				mimirValue = val
			}

			if _, err := w.Write([]byte(strconv.Itoa(mimirValue))); err != nil {
				c.Error(err)
			}
		}
	})

	// setup scanner
	mss := NewMockScannerStorage()
	server := httptest.NewServer(h)
	defer server.Close()
	bridge, err := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       server.Listener.Addr().String(),
		ChainRPC:        server.Listener.Addr().String(),
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}, s.m, s.keys)
	c.Assert(err, IsNil)

	cbs, err := NewBlockScanner(config.BifrostBlockScannerConfiguration{
		StartBlockHeight:           1, // avoids querying thorchain for block height
		BlockScanProcessors:        1,
		HTTPRequestTimeout:         time.Second,
		HTTPRequestReadTimeout:     time.Second * 30,
		HTTPRequestWriteTimeout:    time.Second * 30,
		MaxHTTPRequestRetry:        3,
		BlockHeightDiscoverBackoff: time.Second,
		BlockRetryInterval:         time.Second,
		ChainID:                    common.ETHChain,
	}, mss, m, bridge, DummyFetcher{})
	c.Check(cbs, NotNil)
	c.Check(err, IsNil)

	// Should not be paused
	isHalted := cbs.isChainPaused()
	c.Assert(isHalted, Equals, false)

	// Setting Halt<chain>Chain should pause
	mimirMap["HaltETHChain"] = 2
	// Wait for one block's time so as to replace the cache with an updated query.
	time.Sleep(constants.MAPRelayChainBlockTime)
	isHalted = cbs.isChainPaused()
	c.Assert(isHalted, Equals, true)
	mimirMap["HaltETHChain"] = 0

	// Setting SolvencyHalt<chain>Chain should pause
	mimirMap["SolvencyHaltETHChain"] = 2
	// Wait for one block's time so as to replace the cache with an updated query.
	time.Sleep(constants.MAPRelayChainBlockTime)
	isHalted = cbs.isChainPaused()
	c.Assert(isHalted, Equals, true)
	mimirMap["SolvencyHaltETHChain"] = 0

	// Setting HaltChainGlobal should pause
	mimirMap["HaltChainGlobal"] = 2
	// Wait for one block's time so as to replace the cache with an updated query.
	time.Sleep(constants.MAPRelayChainBlockTime)
	isHalted = cbs.isChainPaused()
	c.Assert(isHalted, Equals, true)
	mimirMap["HaltChainGlobal"] = 0

	// Setting NodePauseChainGlobal should pause
	mimirMap["NodePauseChainGlobal"] = 4 // node pause only halts for an hour, so pause height needs to be larger than thor height
	// Wait for one block's time so as to replace the cache with an updated query.
	time.Sleep(constants.MAPRelayChainBlockTime)
	isHalted = cbs.isChainPaused()
	c.Assert(isHalted, Equals, true)
}
