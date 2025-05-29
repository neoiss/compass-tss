package mapo

import "math/big"

func (b *Bridge) Stop() {
	b.ethClient.Close()
	close(b.stopChan)
}

// GetGasPrice gets gas price from eth scanner
func (b *Bridge) GetGasPrice() *big.Int {
	return b.gasPrice
}
