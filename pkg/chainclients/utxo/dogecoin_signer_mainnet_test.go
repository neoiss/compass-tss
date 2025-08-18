//go:build !testnet
// +build !testnet

package utxo

import (
	"github.com/eager7/dogd/chaincfg"
	. "gopkg.in/check.v1"
)

func (s *DogecoinSignerSuite) TestGetChainCfg(c *C) {
	param := s.client.getChainCfgDOGE()
	c.Assert(param, Equals, &chaincfg.MainNetParams)
}
