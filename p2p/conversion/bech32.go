package conversion

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SetupBech32Prefix() {
	config := sdk.GetConfig()
	// relay will import go-tss as a library , thus this is not needed, we copy the prefix here to avoid go-tss to import relay
	config.SetBech32PrefixForAccount("relay", "relaypub")
	config.SetBech32PrefixForValidator("relayv", "relayvpub")
	config.SetBech32PrefixForConsensusNode("relayc", "relaycpub")
}
