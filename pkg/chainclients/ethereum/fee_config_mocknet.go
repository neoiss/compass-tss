//go:build mocknet
// +build mocknet

package ethereum

import "github.com/mapprotocol/compass-tss/common"

const (
	// initialGasPrice overrides the initial gas price in mocknet to force a reported fee.
	initialGasPrice = 2 * common.One * 100
)
