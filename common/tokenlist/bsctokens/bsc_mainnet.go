//go:build !testnet
// +build !testnet

package bsctokens

import (
	_ "embed"
)

//go:embed bsc_mainnet_latest.json
var BSCTokenListRaw []byte
