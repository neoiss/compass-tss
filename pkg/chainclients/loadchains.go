package chainclients

import (
	"github.com/mapprotocol/compass-tss/internal/keys"
	"time"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/ethereum"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/evm"
	//"github.com/mapprotocol/compass-tss/pkg/chainclients/gaia"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/utxo"
	//"github.com/mapprotocol/compass-tss/pkg/chainclients/xrp"
	"github.com/mapprotocol/compass-tss/pubkeymanager"
	"github.com/mapprotocol/compass-tss/tss/go-tss/tss"
	"github.com/rs/zerolog/log"
)

// ChainClient exports the shared type.
type ChainClient = types.ChainClient

// LoadChains returns chain clients from chain configuration
func LoadChains(thorKeys *keys.Keys,
	cfg map[common.Chain]config.BifrostChainConfiguration,
	server *tss.TssServer,
	thorchainBridge shareTypes.Bridge,
	m *metrics.Metrics,
	pubKeyValidator pubkeymanager.PubKeyValidator,
	poolMgr mapo.PoolManager,
) (chains map[common.Chain]ChainClient, restart chan struct{}) {
	logger := log.Logger.With().Str("module", "bifrost").Logger()

	chains = make(map[common.Chain]ChainClient)
	restart = make(chan struct{})
	failedChains := []common.Chain{}

	loadChain := func(chain config.BifrostChainConfiguration) (ChainClient, error) {
		switch chain.ChainID {
		case common.ETHChain:
			return ethereum.NewClient(thorKeys, chain, server, thorchainBridge, m, pubKeyValidator, poolMgr)
		case common.BSCChain: // ,common.BASEChain:
			return evm.NewEVMClient(thorKeys, chain, server, thorchainBridge, m, pubKeyValidator, poolMgr)
		//case common.GAIAChain:
		//	return gaia.NewCosmosClient(thorKeys, chain, server, thorchainBridge, m)
		case common.BTCChain, common.DOGEChain:
			return utxo.NewClient(thorKeys, chain, server, thorchainBridge, m)
		//case common.XRPChain:
		//	return xrp.NewClient(thorKeys, chain, server, thorchainBridge, m)
		default:
			log.Fatal().Msgf("chain %s is not supported", chain.ChainID)
			return nil, nil
		}
	}

	for _, chain := range cfg {
		if chain.Disabled {
			logger.Info().Msgf("%s chain is disabled by configure", chain.ChainID)
			continue
		}

		client, err := loadChain(chain)
		if err != nil {
			logger.Error().Err(err).Stringer("chain", chain.ChainID).Msg("failed to load chain")
			failedChains = append(failedChains, chain.ChainID)
			continue
		}

		// trunk-ignore-all(golangci-lint/forcetypeassert)
		switch chain.ChainID {
		case common.BTCChain, common.BCHChain, common.LTCChain, common.DOGEChain:
			pubKeyValidator.RegisterCallback(client.(*utxo.Client).RegisterPublicKey)
		}
		chains[chain.ChainID] = client
	}

	// watch failed chains minutely and restart bifrost if any succeed init
	if len(failedChains) > 0 {
		go func() {
			tick := time.NewTicker(time.Minute)
			for range tick.C {
				for _, chain := range failedChains {
					ccfg := cfg[chain]
					ccfg.BlockScanner.DBPath = "" // in-memory db

					_, err := loadChain(ccfg)
					if err == nil {
						logger.Info().Stringer("chain", chain).Msg("chain loaded, restarting bifrost")
						close(restart)
						return
					} else {
						logger.Error().Err(err).Stringer("chain", chain).Msg("failed to load chain")
					}
				}
			}
		}()
	}

	return chains, restart
}
