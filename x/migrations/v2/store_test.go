package v2_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	modtestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/mapprotocol/compass-tss/x/types"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type MigrationsV2Suite struct{}

var _ = Suite(&MigrationsV2Suite{})

func (s *MigrationsV2Suite) TestV2Migrations(c *C) {
	encodingConfig := modtestutil.MakeTestEncodingConfig()
	storeKey := storetypes.NewKVStoreKey("thorchain")
	ctx := testutil.DefaultContext(storeKey, storetypes.NewTransientStoreKey("transient_test"))
	key := "_ver//"
	store := ctx.KVStore(storeKey)

	// Version from v2.138.3, but can be any int64
	ver := types.ProtoInt64{Value: 138}
	store.Set([]byte(key), encodingConfig.Codec.MustMarshal(&ver))
	c.Check(store.Has([]byte(key)), Equals, true)

	//err := v2.MigrateStore(ctx, storeKey)
	//c.Assert(err, IsNil)

	c.Check(store.Has([]byte(key)), Equals, false)
}
