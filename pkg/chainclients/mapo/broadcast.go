package mapo

import (
	"context"
	"time"

	etypes "github.com/ethereum/go-ethereum/core/types"
	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
)

// Broadcast Broadcasts tx to mapBridge
func (b *Bridge) Broadcast(txOutItem *stypes.TxOutItem, hexTx []byte) (string, error) {
	// done
	b.broadcastLock.Lock()
	defer b.broadcastLock.Unlock()

	// decode the transaction
	tx := &etypes.Transaction{}
	if err := tx.UnmarshalJSON(hexTx); err != nil {
		return "", err
	}
	txID := tx.Hash().String()

	// get context with default timeout
	ctx, cancel := b.getTimeoutContext()
	defer cancel()

	// send the transaction
	if err := b.ethClient.SendTransaction(ctx, tx); !isAcceptableError(err) {
		b.logger.Error().Str("txid", txID).Err(err).Msg("failed to send transaction")
		return "", err
	}
	b.logger.Info().Str("memo", txOutItem.Memo).Str("txid", txID).Msg("broadcast tx")

	//// update the signer cache, send to map don’t need cache
	//if err := b.signerCacheManager.SetSigned(txOutItem.CacheHash(), txOutItem.CacheVault(b.GetChain()), txID); err != nil {
	//	b.logger.Err(err).Interface("txOutItem", txOutItem).Msg("fail to mark tx out item as signed")
	//}

	blockHeight, err := b.GetBlockHeight()
	if err != nil {
		b.logger.Err(err).Msg("fail to get current THORChain block height")
		// at this point , the tx already broadcast successfully , don't return an error
		// otherwise will cause the same tx to retry
	} else if err = b.AddSignedTxItem(txID, blockHeight, txOutItem.VaultPubKey.String(), txOutItem); err != nil {
		b.logger.Err(err).Str("hash", txID).Msg("fail to add signed tx item")
	}

	return txID, nil
}

func (b *Bridge) getTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}
