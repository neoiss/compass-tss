//go:build testnet
// +build testnet

package basetokens

import _ "embed"

//go:embed base_testnet_latest.json
var BASETokenListRaw []byte
