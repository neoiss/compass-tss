//go:build testnet
// +build testnet

package ethtokens

import (
	_ "embed"
)

//go:embed eth_testnet_latest.json
var ETHTokenListRaw []byte
