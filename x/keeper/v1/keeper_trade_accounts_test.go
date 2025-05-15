package keeperv1

import (
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
)

type KeeperTradeAccountSuite struct{}

var _ = Suite(&KeeperTradeAccountSuite{})

func (mas *KeeperTradeAccountSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

func (s *KeeperTradeAccountSuite) TestTradeAccount(c *C) {
	ctx, k := setupKeeperForTest(c)
	asset := common.ETHAsset
	addr := GetRandomBech32Addr()

	tr, err := k.GetTradeAccount(ctx, addr, asset)
	c.Assert(err, IsNil)
	c.Check(tr.Units.IsZero(), Equals, true)

	tr.Units = cosmos.NewUint(12)
	k.SetTradeAccount(ctx, tr)
	tr, err = k.GetTradeAccount(ctx, tr.Owner, asset)
	c.Assert(err, IsNil)
	c.Check(tr.Asset.Equals(asset), Equals, true)
	c.Check(tr.Units.Equal(cosmos.NewUint(12)), Equals, true)
	iter := k.GetTradeAccountIteratorWithAddress(ctx, addr)
	c.Check(iter, NotNil)
	iter.Close()
	k.RemoveTradeAccount(ctx, tr)
}
