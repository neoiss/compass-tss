package mapo

import (
	"fmt"
	"math/big"

	"github.com/mapprotocol/compass-tss/constants"
	"github.com/pkg/errors"
)

func (b *Bridge) Stop() {
	b.ethClient.Close()
	close(b.stopChan)
}

// GetGasPrice gets gas price from eth scanner
func (b *Bridge) GetGasPrice() *big.Int {
	return b.gasPrice
}

// GetConstants from relay node
func (b *Bridge) GetConstants() (map[string]int64, error) {
	keys := []string{
		"ChurnInterval",
		"SigningTransactionPeriod",
	}
	input, err := b.cfgAbi.Pack(constants.BatchGetIntValue, keys)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pack input of batchGetIntValue")
	}
	var retSlice []*big.Int
	err = b.callContract(&retSlice, b.cfg.Configuration, constants.BatchGetIntValue, input, b.cfgAbi)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract batchGetIntValue")
	}

	ret := make(map[string]int64)
	for i, v := range retSlice {
		ret[keys[i]] = v.Int64()
	}

	return ret, nil
}

// GetMimir - get mimir settings
func (b *Bridge) GetMimir(key string) (int64, error) {
	input, err := b.cfgAbi.Pack(constants.GetIntValue, key)
	if err != nil {
		return 0, errors.Wrap(err, "unable to pack input of getIntValue")
	}
	var ret *big.Int
	err = b.callContract(&ret, b.cfg.Configuration, constants.GetIntValue, input, b.cfgAbi)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call contract getIntValue")
	}
	return ret.Int64(), nil
}

// GetMimirWithRef is a helper function to more readably insert references (such as Asset MimirString or Chain) into Mimir key templates.
func (b *Bridge) GetMimirWithRef(template, ref string) (int64, error) {
	key := fmt.Sprintf(template, ref)
	return b.GetMimir(key)
}

// GetMimirWithBytes is a helper function to more readably insert references (such as Asset MimirString or Chain) into Mimir key templates.
func (b *Bridge) GetMimirWithBytes(template, ref string) ([]byte, error) {
	key := fmt.Sprintf(template, ref)
	method := constants.GetBytesValue
	input, err := b.cfgAbi.Pack(method, key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pack input of getBytesValue")
	}
	var ret []byte
	err = b.callContract(&ret, b.cfg.Configuration, method, input, b.cfgAbi)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call contract getBytesValue")
	}
	return ret, nil
}
