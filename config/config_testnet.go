//go:build testnet
// +build testnet

package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/mapprotocol/compass-tss/constants"
)

const (
	rpcPort      = 26657
	p2pPort      = 26656
	ebifrostPort = 50051
)

func getSeedAddrs() []string {
	return []string{}
}

func assertBifrostHasSeeds() {
	// fail if seed file is missing or empty since compass-tss will hang
	seedPath := os.ExpandEnv(fmt.Sprintf("$HOME/%s/address_book.seed", constants.DefaultHome))
	fi, err := os.Stat(seedPath)
	if os.IsNotExist(err) {
		log.Warn().Msg("no seed file found")
	}
	if err != nil {
		log.Warn().Err(err).Msg("failed to stat seed file")
		return
	}
	if fi.Size() == 0 {
		log.Warn().Msg("seed file is empty")
	}
}
