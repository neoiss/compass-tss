package types

import (
	"errors"

	se "github.com/cosmos/cosmos-sdk/types/errors"
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
)

type MsgObservedTxInSuite struct{}

var _ = Suite(&MsgObservedTxInSuite{})

func (s *MsgObservedTxInSuite) TestMsgObservedTxIn(c *C) {
	var err error
	pk := GetRandomPubKey()
	tx := common.NewObservedTx(GetRandomTx(), 55, pk, 55)
	acc := GetRandomBech32Addr()
	tx.Tx.ToAddress, err = pk.GetAddress(tx.Tx.Coins[0].Asset.Chain)
	c.Assert(err, IsNil)

	m := NewMsgObservedTxIn(common.ObservedTxs{tx}, acc)
	EnsureMsgBasicCorrect(m, c)

	m1 := NewMsgObservedTxIn(nil, acc)
	c.Assert(m1.ValidateBasic(), NotNil)
	m2 := NewMsgObservedTxIn(common.ObservedTxs{tx}, cosmos.AccAddress{})
	c.Assert(m2.ValidateBasic(), NotNil)

	// will not accept observations with pre-determined signers. This is
	// important to ensure an observer can fake signers from other node accounts
	// *IMPORTANT* DON'T REMOVE THIS CHECK
	tx.Signers = append(tx.Signers, GetRandomBech32Addr().String())
	m3 := NewMsgObservedTxIn(common.ObservedTxs{tx}, acc)
	c.Assert(m3.ValidateBasic(), NotNil)

	tx4 := common.NewObservedTx(GetRandomTx(), 1, pk, 1)
	m4 := NewMsgObservedTxIn(common.ObservedTxs{tx4}, acc)
	err4 := m4.ValidateBasic()
	c.Assert(err4, NotNil)
	c.Assert(errors.Is(err4, se.ErrUnknownRequest), Equals, true)

	tx5 := common.NewObservedTx(GetRandomTx(), 1, pk, 1)
	tx5.Tx.ToAddress, err = pk.GetAddress(tx.Tx.Coins[0].Asset.Chain)
	c.Assert(err, IsNil)
	tx5.OutHashes = []string{
		GetRandomTxHash().String(),
	}
	m5 := NewMsgObservedTxIn(common.ObservedTxs{tx5}, acc)
	err5 := m5.ValidateBasic()
	c.Assert(err5, NotNil)
	c.Assert(errors.Is(err4, se.ErrUnknownRequest), Equals, true)

	tx6 := common.NewObservedTx(GetRandomTx(), 1, pk, 1)
	tx6.Tx.FromAddress = common.Address("")
	m6 := NewMsgObservedTxIn(common.ObservedTxs{tx6}, acc)
	err6 := m6.ValidateBasic()
	c.Assert(err6, NotNil)
	c.Assert(errors.Is(err4, se.ErrUnknownRequest), Equals, true)

	tx7 := common.NewObservedTx(GetRandomTx(), 1, "whatever", 1)
	m7 := NewMsgObservedTxIn(common.ObservedTxs{tx7}, acc)
	err7 := m7.ValidateBasic()
	c.Assert(err7, NotNil)
	c.Assert(errors.Is(err4, se.ErrUnknownRequest), Equals, true)
}
