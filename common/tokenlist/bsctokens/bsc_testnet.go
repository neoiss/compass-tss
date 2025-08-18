//go:build testnet
// +build testnet

package bsctokens

import _ "embed"

//go:embed bsc_testnet_latest.json
var BSCTokenListRaw []byte
