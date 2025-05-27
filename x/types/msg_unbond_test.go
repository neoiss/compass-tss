package types

import (
	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
	. "gopkg.in/check.v1"
)

type MsgUnBondSuite struct{}

var _ = Suite(&MsgUnBondSuite{})

func (mas *MsgUnBondSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

func (MsgUnBondSuite) TestMsgUnBond(c *C) {
	nodeAddr := GetRandomBech32Addr()
	txId := GetRandomTxHash()
	c.Check(txId.IsEmpty(), Equals, false)
	signerAddr := GetRandomBech32Addr()
	bondAddr := GetRandomETHAddress()
	txin := GetRandomTx()
	txinNoID := txin
	txinNoID.ID = ""
	msgApply := NewMsgUnBond(txin, nodeAddr, cosmos.NewUint(common.One), bondAddr, nil, signerAddr)
	c.Assert(msgApply.ValidateBasic(), IsNil)
	c.Assert(len(msgApply.GetSigners()), Equals, 1)
	c.Assert(msgApply.GetSigners()[0].Equals(signerAddr), Equals, true)
	c.Assert(NewMsgUnBond(txin, cosmos.AccAddress{}, cosmos.NewUint(common.One), bondAddr, nil, signerAddr).ValidateBasic(), NotNil)
	c.Assert(NewMsgUnBond(txin, nodeAddr, cosmos.ZeroUint(), bondAddr, nil, signerAddr).ValidateBasic(), IsNil)
	c.Assert(NewMsgUnBond(txinNoID, nodeAddr, cosmos.NewUint(common.One), bondAddr, nil, signerAddr).ValidateBasic(), NotNil)
	c.Assert(NewMsgUnBond(txin, nodeAddr, cosmos.NewUint(common.One), "", nil, signerAddr).ValidateBasic(), NotNil)
	c.Assert(NewMsgUnBond(txin, nodeAddr, cosmos.NewUint(common.One), bondAddr, nil, cosmos.AccAddress{}).ValidateBasic(), NotNil)
}
