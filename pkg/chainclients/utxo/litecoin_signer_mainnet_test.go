//go:build !testnet
// +build !testnet

package utxo

import (
	"github.com/ltcsuite/ltcd/chaincfg"
	. "gopkg.in/check.v1"
)

func (s *LitecoinSignerSuite) TestGetChainCfg(c *C) {
	param := s.client.getChainCfgLTC()
	c.Assert(param, Equals, &chaincfg.MainNetParams)
}
