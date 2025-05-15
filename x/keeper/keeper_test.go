package keeper

import (
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
)

func FundModule(c *C, ctx cosmos.Context, k Keeper, name string, amt uint64) {
	coin := common.NewCoin(common.RuneNative, cosmos.NewUint(amt))
	err := k.MintToModule(ctx, ModuleName, coin)
	c.Assert(err, IsNil)
	err = k.SendFromModuleToModule(ctx, ModuleName, name, common.NewCoins(coin))
	c.Assert(err, IsNil)
}

func FundAccount(c *C, ctx cosmos.Context, k Keeper, addr cosmos.AccAddress, amt uint64) {
	coin := common.NewCoin(common.RuneNative, cosmos.NewUint(amt))
	err := k.MintToModule(ctx, ModuleName, coin)
	c.Assert(err, IsNil)
	err = k.SendFromModuleToAccount(ctx, ModuleName, addr, common.NewCoins(coin))
	c.Assert(err, IsNil)
}
