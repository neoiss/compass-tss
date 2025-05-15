package types

import (
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	"gitlab.com/thorchain/thornode/v3/common/cosmos"
)

type MsgRagnarokSuite struct{}

var _ = Suite(&MsgRagnarokSuite{})

func (MsgRagnarokSuite) TestMsgRagnarokSuite(c *C) {
	txID := GetRandomTxHash()
	eth := GetRandomETHAddress()
	acc1 := GetRandomBech32Addr()
	tx := common.NewObservedTx(common.NewTx(
		txID,
		eth,
		GetRandomETHAddress(),
		common.Coins{common.NewCoin(common.ETHAsset, cosmos.OneUint())},
		common.Gas{common.NewCoin(common.ETHAsset, cosmos.NewUint(common.One))},
		"ragnarok:10",
	), 12, GetRandomPubKey(), 12)
	m := NewMsgRagnarok(tx, 10, acc1)
	EnsureMsgBasicCorrect(m, c)

	inputs := []struct {
		txID        common.TxID
		blockHeight int64
		sender      common.Address
		signer      cosmos.AccAddress
	}{
		{
			txID:        common.TxID(""),
			blockHeight: 1,
			sender:      eth,
			signer:      acc1,
		},
		{
			txID:        txID,
			blockHeight: 0,
			sender:      eth,
			signer:      acc1,
		},
		{
			txID:        txID,
			blockHeight: 1,
			sender:      common.NoAddress,
			signer:      acc1,
		},
		{
			txID:        txID,
			blockHeight: 1,
			sender:      eth,
			signer:      cosmos.AccAddress{},
		},
	}

	for _, item := range inputs {
		tx = common.NewObservedTx(common.NewTx(
			item.txID,
			item.sender,
			GetRandomETHAddress(),
			common.Coins{common.NewCoin(common.ETHAsset, cosmos.OneUint())},
			common.Gas{common.NewCoin(common.ETHAsset, cosmos.NewUint(common.One))},
			"",
		), 12, GetRandomPubKey(), 12)
		m = NewMsgRagnarok(tx, item.blockHeight, item.signer)
		err := m.ValidateBasic()
		c.Assert(err, NotNil, Commentf("%s", err.Error()))
	}
}
