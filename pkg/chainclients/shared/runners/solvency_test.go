package runners

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/mapprotocol/compass-tss/cmd"
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/metrics"
	mapclient "github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/mapprotocol/compass-tss/x/types"
)

func TestPackage(t *testing.T) { TestingT(t) }

type SolvencyTestSuite struct {
	sp   *DummySolvencyCheckProvider
	m    *metrics.Metrics
	cfg  config.BifrostClientConfiguration
	keys *mapclient.Keys
}

var _ = Suite(&SolvencyTestSuite{})

func (s *SolvencyTestSuite) SetUpSuite(c *C) {
	sp := &DummySolvencyCheckProvider{}
	s.sp = sp

	m, _ := metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9090,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.ETHChain},
	})
	s.m = m

	cfg := config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       "localhost",
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb := ckeys.NewInMemory(cdc)
	_, _, err := kb.NewMnemonic(cfg.SignerName, ckeys.English, cmd.THORChainHDPath, cfg.SignerPasswd, hd.Secp256k1)
	c.Assert(err, IsNil)
	s.cfg = cfg
	s.keys = mapclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd)

	c.Assert(err, IsNil)
}

func (s *SolvencyTestSuite) TestSolvencyCheck(c *C) {
	mimirMap := map[string]int{
		"HaltETHChain":         0,
		"SolvencyHaltETHChain": 0,
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logf("================>:%s", r.RequestURI)
		if strings.HasPrefix(r.RequestURI, mapclient.MimirEndpoint) {
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

	server := httptest.NewServer(h)
	defer server.Close()
	bridge, _ := mapclient.NewThorchainBridge(config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       server.Listener.Addr().String(),
		ChainRPC:        server.Listener.Addr().String(),
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: ".",
	}, s.m, s.keys)

	stopchan := make(chan struct{})
	wg := &sync.WaitGroup{}

	// Happy path, shouldn't check solvency if nothing halted (chain clients will report solvency)
	s.sp.ResetChecks()
	wg.Add(1)
	go SolvencyCheckRunner(common.ETHChain, s.sp, bridge, stopchan, wg, constants.MAPRelayChainBlockTime)
	time.Sleep(time.Second * 6)

	c.Assert(s.sp.ShouldReportSolvencyRan, Equals, false)
	c.Assert(s.sp.ReportSolvencyRun, Equals, false)

	// Admin halted, still don't check solvency
	mimirMap["HaltETHChain"] = 1
	s.sp.ResetChecks()
	wg.Add(1)
	go SolvencyCheckRunner(common.ETHChain, s.sp, bridge, stopchan, wg, constants.MAPRelayChainBlockTime)
	time.Sleep(time.Second * 6)

	c.Assert(s.sp.ShouldReportSolvencyRan, Equals, false)
	c.Assert(s.sp.ReportSolvencyRun, Equals, false)

	// Double-spend check halted chain client, check solvency here
	mimirMap["HaltETHChain"] = 10
	s.sp.ResetChecks()
	wg.Add(1)
	go SolvencyCheckRunner(common.ETHChain, s.sp, bridge, stopchan, wg, constants.MAPRelayChainBlockTime)
	time.Sleep(time.Second * 6)

	c.Assert(s.sp.ShouldReportSolvencyRan, Equals, true)
	c.Assert(s.sp.ReportSolvencyRun, Equals, true)
	mimirMap["HaltETHChain"] = 0

	// Solvency halted chain, need to report solvency here as chain client is paused
	mimirMap["SolvencyHaltETHChain"] = 1
	s.sp.ResetChecks()
	wg.Add(1)
	go SolvencyCheckRunner(common.ETHChain, s.sp, bridge, stopchan, wg, constants.MAPRelayChainBlockTime)
	time.Sleep(time.Second * 6)

	c.Assert(s.sp.ShouldReportSolvencyRan, Equals, true)
	c.Assert(s.sp.ReportSolvencyRun, Equals, true)
}

// Mock SolvencyCheckProvider
type DummySolvencyCheckProvider struct {
	ShouldReportSolvencyRan bool
	ReportSolvencyRun       bool
}

func (d *DummySolvencyCheckProvider) ResetChecks() {
	d.ShouldReportSolvencyRan = false
	d.ReportSolvencyRun = false
}

func (d *DummySolvencyCheckProvider) GetHeight() (int64, error) {
	return 0, nil
}

func (d *DummySolvencyCheckProvider) ShouldReportSolvency(height int64) bool {
	d.ShouldReportSolvencyRan = true
	return true
}

func (d *DummySolvencyCheckProvider) ReportSolvency(height int64) error {
	d.ReportSolvencyRun = true
	return nil
}

func (s *SolvencyTestSuite) TestIsVaultSolvent(c *C) {
	vault := types.Vault{
		BlockHeight: 1,
		PubKey:      types.GetRandomPubKey(),
		Coins: common.NewCoins(
			common.NewCoin(common.ETHAsset, cosmos.NewUint(102400000000)),
		),
		Type:   types.VaultType_AsgardVault,
		Status: types.VaultStatus_ActiveVault,
	}
	acct := common.Account{
		Sequence:      0,
		AccountNumber: 0,
		Coins:         common.NewCoins(common.NewCoin(common.ETHAsset, cosmos.NewUint(102400000000))),
	}
	c.Assert(IsVaultSolvent(acct, vault, cosmos.NewUint(0)), Equals, true)
	acct = common.Account{
		Sequence:      0,
		AccountNumber: 0,
		Coins:         common.NewCoins(common.NewCoin(common.ETHAsset, cosmos.NewUint(102305000000))),
	}
	c.Assert(IsVaultSolvent(acct, vault, cosmos.NewUint(80000*120)), Equals, true)
	acct = common.Account{
		Sequence:      0,
		AccountNumber: 0,
		Coins:         common.NewCoins(common.NewCoin(common.ETHAsset, cosmos.NewUint(102205000000))),
	}
	c.Assert(IsVaultSolvent(acct, vault, cosmos.NewUint(80000*120)), Equals, false)
}
