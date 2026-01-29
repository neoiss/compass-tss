//go:build !testnet
// +build !testnet

package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/constants"
)

const (
	rpcPort      = 27147
	p2pPort      = 27146
	ebifrostPort = 50051
)

func getSeedAddrs() (addrs []string) {
	// todo get ips by contract
	if config.MAPO.SeedNodesEndpoint == "" {
		log.Warn().Msg("no seed nodes endpoint provided")
		return
	}
	return []string{}
}

func assertBifrostHasSeeds() {
	if config.MAPO.SeedNodesEndpoint == "" {
		log.Warn().Msg("no seed nodes endpoint provided, skipping seed file check")
		return
	}

	// fail if seed file is missing or empty since compass-tss will hang
	seedPath := os.ExpandEnv(fmt.Sprintf("$HOME/%s/address_book.seed", constants.DefaultHome))
	fi, err := os.Stat(seedPath)
	if os.IsNotExist(err) {
		log.Fatal().Msg("no seed file found")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to stat seed file")
	}
	if fi.Size() == 0 {
		log.Fatal().Msg("seed file is empty")
	}
}
