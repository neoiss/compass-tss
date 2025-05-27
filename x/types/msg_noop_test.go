package types

import (
	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
	. "gopkg.in/check.v1"
)

type MsgNoopSuite struct{}

var _ = Suite(&MsgNoopSuite{})

func (MsgNoopSuite) TestMsgNoop(c *C) {
	addr := GetRandomBech32Addr()
	c.Check(addr.Empty(), Equals, false)
	tx := common.ObservedTx{
		Tx:             GetRandomTx(),
		Status:         common.Status_done,
		OutHashes:      nil,
		BlockHeight:    1,
		Signers:        []string{addr.String()},
		ObservedPubKey: GetRandomPubKey(),
		FinaliseHeight: 1,
	}
	m := NewMsgNoOp(tx, addr, "")
	c.Check(m.ValidateBasic(), IsNil)
	EnsureMsgBasicCorrect(m, c)
	mEmpty := NewMsgNoOp(tx, cosmos.AccAddress{}, "")
	c.Assert(mEmpty.ValidateBasic(), NotNil)
}
