//go:build !testnet
// +build !testnet

package utxo

import (
	"github.com/btcsuite/btcd/chaincfg"
	. "gopkg.in/check.v1"
)

func (s *BitcoinSignerSuite) TestGetChainCfg(c *C) {
	param := s.client.getChainCfgBTC()
	c.Assert(param, Equals, &chaincfg.MainNetParams)
}
