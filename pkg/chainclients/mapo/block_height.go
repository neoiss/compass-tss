package mapo

func (b *Bridge) GetMapBlockHeight() (int64, error) {
	return b.ethRpc.GetBlockHeight()
}
