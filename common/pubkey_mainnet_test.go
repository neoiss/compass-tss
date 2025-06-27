//go:build !testnet
// +build !testnet

package common

import (
	"encoding/hex"

	. "gopkg.in/check.v1"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

func (s *PubKeyTestSuite) TestPubKeyGetAddress(c *C) {
	for _, d := range s.keyData {
		privB, _ := hex.DecodeString(d.priv)
		pubB, _ := hex.DecodeString(d.pub)
		priv := secp256k1.PrivKey(privB)
		pubKey := priv.PubKey()
		pubT, _ := pubKey.(secp256k1.PubKey)
		pub := pubT[:]

		c.Assert(hex.EncodeToString(pub), Equals, hex.EncodeToString(pubB))

		tmp, err := codec.FromCmtPubKeyInterface(pubKey)
		c.Assert(err, IsNil)
		pubBech32, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, tmp)
		c.Assert(err, IsNil)

		pk, err := NewPubKey(pubBech32)
		c.Assert(err, IsNil)

		addrETH, err := pk.GetAddress(ETHChain)
		c.Assert(err, IsNil)
		c.Assert(addrETH.String(), Equals, d.addrETH.mainnet)

		addrBTC, err := pk.GetAddress(BTCChain)
		c.Assert(err, IsNil)
		c.Assert(addrBTC.String(), Equals, d.addrBTC.mainnet)

		addrLTC, err := pk.GetAddress(LTCChain)
		c.Assert(err, IsNil)
		c.Assert(addrLTC.String(), Equals, d.addrLTC.mainnet)

		addrBCH, err := pk.GetAddress(BCHChain)
		c.Assert(err, IsNil)
		c.Assert(addrBCH.String(), Equals, d.addrBCH.mainnet)

		addrDOGE, err := pk.GetAddress(DOGEChain)
		c.Assert(err, IsNil)
		c.Assert(addrDOGE.String(), Equals, d.addrDOGE.mainnet)
	}
}
