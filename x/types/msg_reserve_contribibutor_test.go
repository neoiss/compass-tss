package types

import (
	"errors"

	se "github.com/cosmos/cosmos-sdk/types/errors"
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

type MsgReserveContributorSuite struct{}

var _ = Suite(&MsgReserveContributorSuite{})

func (s *MsgReserveContributorSuite) TestMsgReserveContributor(c *C) {
	addr := GetRandomETHAddress()
	amt := cosmos.NewUint(378 * common.One)
	res := NewReserveContributor(addr, amt)
	signer := GetRandomBech32Addr()

	msg := NewMsgReserveContributor(GetRandomTx(), res, signer)
	c.Check(msg.Contributor.IsEmpty(), Equals, false)
	c.Check(msg.Signer.Equals(signer), Equals, true)
	EnsureMsgBasicCorrect(msg, c)

	tx1 := GetRandomTx()
	tx1.FromAddress = ""
	msg1 := NewMsgReserveContributor(tx1, res, signer)
	err1 := msg1.ValidateBasic()
	c.Assert(err1, NotNil)
	c.Assert(errors.Is(err1, se.ErrUnknownRequest), Equals, true)

	msg2 := NewMsgReserveContributor(GetRandomTx(), res, cosmos.AccAddress{})
	err2 := msg2.ValidateBasic()
	c.Assert(err2, NotNil)
	c.Assert(errors.Is(err2, se.ErrInvalidAddress), Equals, true)

	msg3 := NewMsgReserveContributor(GetRandomTx(), NewReserveContributor("", cosmos.ZeroUint()), signer)
	err3 := msg3.ValidateBasic()
	c.Assert(err3, NotNil)
	c.Assert(errors.Is(err3, se.ErrUnknownRequest), Equals, true)
}
