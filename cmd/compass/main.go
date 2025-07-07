package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	golog "github.com/ipfs/go-log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"io"

	tcommon "github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/constants"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/observer"
	"github.com/mapprotocol/compass-tss/p2p"
	"github.com/mapprotocol/compass-tss/pkg/chainclients"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/mapprotocol/compass-tss/signer"
	ctss "github.com/mapprotocol/compass-tss/tss"
	"github.com/mapprotocol/compass-tss/tss/go-tss/common"
	"github.com/mapprotocol/compass-tss/tss/go-tss/tss"
)

// THORNode define version / revision here , so THORNode could inject the version from CI pipeline if THORNode want to
var (
	version  string
	revision string
)

const (
	serverIdentity = "compass-tss"
)

func printVersion() {
	fmt.Printf("%s v%s, rev %s\n", serverIdentity, version, revision)
}

func main() {
	showVersion := flag.Bool("version", false, "Shows version")
	logLevel := flag.StringP("log-level", "l", "info", "Log Level")
	pretty := flag.BoolP("pretty-log", "p", false, "Enables unstructured prettified logging. This is useful for local debugging")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	initPrefix()
	initLog(*logLevel, *pretty)
	config.Init()
	config.InitBifrost()
	cfg := config.GetBifrost()

	// metrics
	m, err := metrics.NewMetrics(cfg.Metrics)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create metric instance")
	}
	if err = m.Start(); err != nil {
		log.Fatal().Err(err).Msg("fail to start metric collector")
	}
	if len(cfg.Thorchain.SignerName) == 0 {
		log.Fatal().Msg("signer name is empty")
	}
	if len(cfg.Thorchain.SignerPasswd) == 0 {
		log.Fatal().Msg("signer password is empty")
	}
	kb, _, err := keys.GetKeyringKeybase(cfg.Thorchain.PrivateKey, cfg.Thorchain.SignerName)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
	}

	k := keys.NewKeysWithKeybase(kb, cfg.Thorchain.SignerName, cfg.Thorchain.SignerPasswd, cfg.Thorchain.PrivateKey)
	// map bridge
	mapBridge, err := mapo.NewBridge(cfg.Thorchain, m, k)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create new map bridge")
	}
	if err = mapBridge.EnsureNodeWhitelistedWithTimeout(); err != nil {
		log.Fatal().Err(err).Msg("node account is not whitelisted, can't start")
	}
	// PubKey Manager
	pubkeyMgr, err := pubkeymanager.NewPubKeyManager(mapBridge, m)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create pubkey manager")
	}
	if err = pubkeyMgr.Start(); err != nil {
		log.Fatal().Err(err).Msg("fail to start pubkey manager")
	}

	// automatically attempt to recover TSS key shares if they are missing
	if err = ctss.RecoverKeyShares(mapBridge); err != nil {
		log.Error().Err(err).Msg("fail to recover key shares")
	}

	// setup TSS signing
	priKey, err := k.GetPrivateKey()
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get private key")
	}

	tmPrivateKey := tcommon.CosmosPrivateKeyToTMPrivateKey(priKey)

	consts := constants.NewConstantValue()
	jailTimeKeygen := time.Duration(consts.GetInt64Value(constants.JailTimeKeygen)) * constants.MAPRelayChainBlockTime
	jailTimeKeysign := time.Duration(consts.GetInt64Value(constants.JailTimeKeysign)) * constants.MAPRelayChainBlockTime
	if cfg.Signer.KeygenTimeout >= jailTimeKeygen {
		log.Fatal().
			Stringer("keygenTimeout", cfg.Signer.KeygenTimeout).
			Stringer("keygenJail", jailTimeKeygen).
			Msg("keygen timeout must be shorter than jail time")
	}
	if cfg.Signer.KeysignTimeout >= jailTimeKeysign {
		log.Fatal().
			Stringer("keysignTimeout", cfg.Signer.KeysignTimeout).
			Stringer("keysignJail", jailTimeKeysign).
			Msg("keysign timeout must be shorter than jail time")
	}

	comm, stateManager, err := p2p.StartP2P(
		cfg.TSS,
		tmPrivateKey,
		constants.DefaultHome,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to start p2p")
	}

	tssIns, err := tss.NewTss(
		comm,
		stateManager,
		tmPrivateKey,
		common.TssConfig{
			EnableMonitor:   true,
			KeyGenTimeout:   cfg.Signer.KeygenTimeout,
			KeySignTimeout:  cfg.Signer.KeysignTimeout,
			PartyTimeout:    cfg.Signer.PartyTimeout,
			PreParamTimeout: cfg.Signer.PreParamTimeout,
		},
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create tss instance")
	}

	if err = tssIns.Start(); err != nil {
		log.Err(err).Msg("fail to start tss instance")
	}

	cfgChains := cfg.GetChains()

	// ensure we have a protocol for chain RPC Hosts
	for _, chainCfg := range cfgChains {
		if chainCfg.Disabled {
			continue
		}
		if len(chainCfg.RPCHost) == 0 {
			log.Fatal().Err(err).Stringer("chain", chainCfg.ChainID).Msg("missing chain RPC host")
			return
		}
		if !strings.HasPrefix(chainCfg.RPCHost, "http") {
			chainCfg.RPCHost = fmt.Sprintf("http://%s", chainCfg.RPCHost)
		}
	}
	poolMgr := mapo.NewPoolMgr(mapBridge)
	chains, restart := chainclients.LoadChains(k, cfgChains, tssIns, mapBridge, m, pubkeyMgr, poolMgr)
	if len(chains) == 0 {
		log.Fatal().Msg("fail to load any chains")
	}
	tssKeysignMetricMgr := metrics.NewTssKeysignMetricMgr()
	healthServer := NewHealthServer(cfg.TSS.InfoAddress, tssIns, chains)
	go func() {
		defer log.Info().Msg("health server exit")
		if err = healthServer.Start(); err != nil {
			log.Error().Err(err).Msg("fail to start health server")
		}
	}()

	ctx := context.Background()

	// start observer notifier
	ag, err := observer.NewAttestationGossip(comm.GetHost(), k, cfg.Thorchain.ChainEBifrost, mapBridge, m, cfg.AttestationGossip)

	// start observer
	obs, err := observer.NewObserver(pubkeyMgr, chains, mapBridge, m, cfgChains[tcommon.BTCChain].BlockScanner.DBPath, tssKeysignMetricMgr, ag)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create observer")
	}
	if err = obs.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("fail to start observer")
	}

	// enable observer to react to notifications from thornode
	// that come through the grpc connection within AttestationGossip.
	ag.SetObserverHandleObservedTxCommitted(obs)

	// start signer
	sign, err := signer.NewSigner(cfg, mapBridge, k, pubkeyMgr, tssIns, chains, m, tssKeysignMetricMgr, obs)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create instance of signer")
	}
	if err = sign.Start(); err != nil {
		log.Fatal().Err(err).Msg("fail to start signer")
	}

	// wait....
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ch:
	case <-restart:
	}
	log.Info().Msg("stop signal received")

	// stop observer
	if err = obs.Stop(); err != nil {
		log.Fatal().Err(err).Msg("fail to stop observer")
	}
	// stop signer
	if err = sign.Stop(); err != nil {
		log.Fatal().Err(err).Msg("fail to stop signer")
	}
	// stop go tss
	tssIns.Stop()
	if err = healthServer.Stop(); err != nil {
		log.Fatal().Err(err).Msg("fail to stop health server")
	}
}

func initPrefix() {
	//// todo
	//cosmosSDKConfg := cosmos.GetConfig()
	//cosmosSDKConfg.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	//cosmosSDKConfg.Seal()
}

func initLog(level string, pretty bool) {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Msgf("%s is not a valid log-level, falling back to 'info'", level)
	}
	var out io.Writer = os.Stdout
	if pretty {
		out = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	zerolog.SetGlobalLevel(l)
	log.Logger = log.Output(out).With().Caller().Str("service", serverIdentity).Logger()

	logLevel := golog.LevelInfo
	switch l {
	case zerolog.DebugLevel:
		logLevel = golog.LevelDebug
	case zerolog.InfoLevel:
		logLevel = golog.LevelInfo
	case zerolog.ErrorLevel:
		logLevel = golog.LevelError
	case zerolog.FatalLevel:
		logLevel = golog.LevelFatal
	case zerolog.PanicLevel:
		logLevel = golog.LevelPanic
	}
	golog.SetAllLoggers(logLevel)
	if err = golog.SetLogLevel("tss-lib", level); err != nil {
		log.Fatal().Err(err).Msg("fail to set tss-lib loglevel")
	}
}
