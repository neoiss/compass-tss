package keeperv1

import (
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
)

type KeeperLiquidityProviderSuite struct{}

var _ = Suite(&KeeperLiquidityProviderSuite{})

func (mas *KeeperLiquidityProviderSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

func (s *KeeperLiquidityProviderSuite) TestLiquidityProvider(c *C) {
	ctx, k := setupKeeperForTest(c)
	asset := common.ETHAsset

	lp, err := k.GetLiquidityProvider(ctx, asset, GetRandomRUNEAddress())
	c.Assert(err, IsNil)
	c.Check(lp.PendingRune, NotNil)
	c.Check(lp.Units, NotNil)

	lp = LiquidityProvider{
		Asset:        asset,
		Units:        cosmos.NewUint(12),
		RuneAddress:  GetRandomRUNEAddress(),
		AssetAddress: GetRandomBTCAddress(),
	}
	k.SetLiquidityProvider(ctx, lp)
	lp, err = k.GetLiquidityProvider(ctx, asset, lp.RuneAddress)
	c.Assert(err, IsNil)
	c.Check(lp.Asset.Equals(asset), Equals, true)
	c.Check(lp.Units.Equal(cosmos.NewUint(12)), Equals, true)
	iter := k.GetLiquidityProviderIterator(ctx, common.ETHAsset)
	c.Check(iter, NotNil)
	iter.Close()
	k.RemoveLiquidityProvider(ctx, lp)
}
