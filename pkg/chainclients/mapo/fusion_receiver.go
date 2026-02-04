package mapo

import (
	ecommon "github.com/ethereum/go-ethereum/common"
)

func (b *Bridge) GetFusionReceiver() ecommon.Address {
	return ecommon.HexToAddress(b.cfg.FusionReceiver)
}
