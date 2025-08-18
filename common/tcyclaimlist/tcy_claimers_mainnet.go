//go:build !testnet
// +build !testnet

package tcyclaimlist

import (
	_ "embed"
)

//go:embed tcy_claimers_mainnet.json
var TCYClaimsListRaw []byte
