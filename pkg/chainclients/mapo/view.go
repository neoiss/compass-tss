package mapo

import (
	"math/big"

	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

type VaultInfo struct {
	PubKey  []byte
	Routers []VaultRouter
}

type VaultRouter struct {
	Chain  *big.Int
	Router []byte
}

func (b *Bridge) getPublickeys() (*VaultInfo, error) {
	method := constants.GetPublickeys
	input, err := b.viewAbi.Pack(method)
	if err != nil {
		return nil, err
	}

	ret := VaultInfo{}
	err = b.callContract(&ret, b.cfg.ViewController, method, input, b.viewAbi)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call %s", method)
	}

	return &ret, nil
}
