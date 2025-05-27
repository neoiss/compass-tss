package types

import (
	. "gopkg.in/check.v1"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

type EventSuite struct{}

var _ = Suite(&EventSuite{})

func (s EventSuite) TestSwapEvent(c *C) {
	evt := NewEventSwap(
		common.ETHAsset,
		cosmos.NewUint(5),
		cosmos.NewUint(5),
		cosmos.NewUint(5),
		cosmos.ZeroUint(),
		GetRandomTx(),
		common.NewCoin(common.ETHAsset, cosmos.NewUint(100)),
		cosmos.NewUint(5),
	)
	c.Check(evt.Type(), Equals, "swap")
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestAddLiqudityEvent(c *C) {
	evt := NewEventAddLiquidity(
		common.ETHAsset,
		cosmos.NewUint(5),
		GetRandomRUNEAddress(),
		cosmos.NewUint(5),
		cosmos.NewUint(5),
		GetRandomTxHash(),
		GetRandomTxHash(),
		GetRandomETHAddress(),
	)
	c.Check(evt.Type(), Equals, "add_liquidity")
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestWithdrawEvent(c *C) {
	evt := NewEventWithdraw(
		common.ETHAsset,
		cosmos.NewUint(6),
		5000,
		cosmos.NewDec(0),
		GetRandomTx(),
		cosmos.NewUint(100),
		cosmos.NewUint(100),
	)
	c.Check(evt.Type(), Equals, "withdraw")
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestPool(c *C) {
	evt := NewEventPool(common.ETHAsset, PoolStatus_Available)
	c.Check(evt.Type(), Equals, "pool")
	c.Check(evt.Pool.String(), Equals, common.ETHAsset.String())
	c.Check(evt.Status.String(), Equals, PoolStatus_Available.String())
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestReward(c *C) {
	evt := NewEventRewards(
		cosmos.NewUint(300),
		[]PoolAmt{{common.ETHAsset, 30}, {common.BTCAsset, 40}},
		cosmos.NewUint(50),
		cosmos.NewUint(60),
		cosmos.NewUint(70),
	)
	c.Check(evt.Type(), Equals, "rewards")
	c.Check(evt.BondReward.String(), Equals, "300")
	c.Assert(evt.PoolRewards, HasLen, 2)
	c.Check(evt.PoolRewards[0].Asset.Equals(common.ETHAsset), Equals, true)
	c.Check(evt.PoolRewards[0].Amount, Equals, int64(30))
	c.Check(evt.PoolRewards[1].Asset.Equals(common.BTCAsset), Equals, true)
	c.Check(evt.PoolRewards[1].Amount, Equals, int64(40))
	c.Check(evt.DevFundReward.String(), Equals, "50")
	c.Check(evt.IncomeBurn.String(), Equals, "60")
	c.Check(evt.TcyStakeReward.String(), Equals, "70")
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestSlash(c *C) {
	evt := NewEventSlash(common.ETHAsset, []PoolAmt{
		{common.ETHAsset, -20},
		{common.RuneAsset(), 30},
	})
	c.Check(evt.Type(), Equals, "slash")
	c.Check(evt.Pool, Equals, common.ETHAsset)
	c.Assert(evt.SlashAmount, HasLen, 2)
	c.Check(evt.SlashAmount[0].Asset, Equals, common.ETHAsset)
	c.Check(evt.SlashAmount[0].Amount, Equals, int64(-20))
	c.Check(evt.SlashAmount[1].Asset, Equals, common.RuneAsset())
	c.Check(evt.SlashAmount[1].Amount, Equals, int64(30))
	events, err := evt.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestEventGas(c *C) {
	eg := NewEventGas()
	c.Assert(eg, NotNil)
	eg.UpsertGasPool(GasPool{
		Asset:    common.ETHAsset,
		AssetAmt: cosmos.NewUint(1000),
		RuneAmt:  cosmos.ZeroUint(),
	})
	c.Assert(eg.Pools, HasLen, 1)
	c.Assert(eg.Pools[0].Asset, Equals, common.ETHAsset)
	c.Assert(eg.Pools[0].RuneAmt.Equal(cosmos.ZeroUint()), Equals, true)
	c.Assert(eg.Pools[0].AssetAmt.Equal(cosmos.NewUint(1000)), Equals, true)

	eg.UpsertGasPool(GasPool{
		Asset:    common.ETHAsset,
		AssetAmt: cosmos.NewUint(1234),
		RuneAmt:  cosmos.NewUint(1024),
	})
	c.Assert(eg.Pools, HasLen, 1)
	c.Assert(eg.Pools[0].Asset, Equals, common.ETHAsset)
	c.Assert(eg.Pools[0].RuneAmt.Equal(cosmos.NewUint(1024)), Equals, true)
	c.Assert(eg.Pools[0].AssetAmt.Equal(cosmos.NewUint(2234)), Equals, true)

	eg.UpsertGasPool(GasPool{
		Asset:    common.BTCAsset,
		AssetAmt: cosmos.NewUint(1024),
		RuneAmt:  cosmos.ZeroUint(),
	})
	c.Assert(eg.Pools, HasLen, 2)
	c.Assert(eg.Pools[1].Asset, Equals, common.BTCAsset)
	c.Assert(eg.Pools[1].AssetAmt.Equal(cosmos.NewUint(1024)), Equals, true)
	c.Assert(eg.Pools[1].RuneAmt.Equal(cosmos.ZeroUint()), Equals, true)

	eg.UpsertGasPool(GasPool{
		Asset:    common.BTCAsset,
		AssetAmt: cosmos.ZeroUint(),
		RuneAmt:  cosmos.ZeroUint(),
	})

	c.Assert(eg.Pools, HasLen, 2)
	c.Assert(eg.Pools[1].Asset, Equals, common.BTCAsset)
	c.Assert(eg.Pools[1].AssetAmt.Equal(cosmos.NewUint(1024)), Equals, true)
	c.Assert(eg.Pools[1].RuneAmt.Equal(cosmos.ZeroUint()), Equals, true)

	eg.UpsertGasPool(GasPool{
		Asset:    common.BTCAsset,
		AssetAmt: cosmos.ZeroUint(),
		RuneAmt:  cosmos.NewUint(3333),
	})

	c.Assert(eg.Pools, HasLen, 2)
	c.Assert(eg.Pools[1].Asset, Equals, common.BTCAsset)
	c.Assert(eg.Pools[1].AssetAmt.Equal(cosmos.NewUint(1024)), Equals, true)
	c.Assert(eg.Pools[1].RuneAmt.Equal(cosmos.NewUint(3333)), Equals, true)
	events, err := eg.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (s EventSuite) TestEventFee(c *C) {
	event := NewEventFee(GetRandomTxHash(), common.Fee{
		Coins: common.Coins{
			common.NewCoin(common.ETHAsset, cosmos.NewUint(1024)),
		},
		PoolDeduct: cosmos.NewUint(1023),
	}, cosmos.NewUint(5))
	c.Assert(event.Type(), Equals, FeeEventType)
	evts, err := event.Events()
	c.Assert(err, IsNil)
	c.Assert(evts, HasLen, 1)
}

func (s EventSuite) TestEventDonate(c *C) {
	e := NewEventDonate(common.ETHAsset, GetRandomTx())
	c.Check(e.Type(), Equals, "donate")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventRefund(c *C) {
	e := NewEventRefund(1, "refund", GetRandomTx(), common.NewFee(common.Coins{
		common.NewCoin(common.ETHAsset, cosmos.NewUint(100)),
	}, cosmos.ZeroUint()))
	c.Check(e.Type(), Equals, "refund")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventBond(c *C) {
	e := NewEventBond(cosmos.NewUint(100), BondType_bond_paid, GetRandomTx(), &NodeAccount{}, nil)
	c.Check(e.Type(), Equals, "bond")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventReserve(c *C) {
	e := NewEventReserve(ReserveContributor{
		Address: GetRandomETHAddress(),
		Amount:  cosmos.NewUint(100),
	}, GetRandomTx())
	c.Check(e.Type(), Equals, "reserve")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventErrata(c *C) {
	e := NewEventErrata(GetRandomTxHash(), PoolMods{
		NewPoolMod(common.ETHAsset, cosmos.NewUint(100), true, cosmos.NewUint(200), true),
	})
	c.Check(e.Type(), Equals, "errata")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventOutbound(c *C) {
	e := NewEventOutbound(GetRandomTxHash(), GetRandomTx())
	c.Check(e.Type(), Equals, "outbound")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventSlashPoint(c *C) {
	e := NewEventSlashPoint(GetRandomBech32Addr(), 100, "what ever")
	c.Check(e.Type(), Equals, "slash_points")
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}

func (EventSuite) TestEventPoolStageCost(c *C) {
	e := NewEventPoolBalanceChanged(NewPoolMod(common.BTCAsset, cosmos.NewUint(100), false, cosmos.ZeroUint(), false), "test")
	c.Check(e.Type(), Equals, PoolBalanceChangeEventType)
	events, err := e.Events()
	c.Check(err, IsNil)
	c.Check(events, NotNil)
}
