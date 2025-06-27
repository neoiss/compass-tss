//go:build !testnet
// +build !testnet

package avaxtokens

import (
	_ "embed"
)

//go:embed avax_mainnet_latest.json
var AVAXTokenListRaw []byte
