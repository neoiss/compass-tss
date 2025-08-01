package utxo

import (
	shareTypes "github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"

	stypes "github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/rs/zerolog"
)

// SignCheckpoint is used to checkpoint the built transaction before signing, for use in
// round 7 signing errors which must reuse the same inputs.
type SignCheckpoint struct {
	UnsignedTx        []byte           `json:"unsigned_tx"`
	IndividualAmounts map[string]int64 `json:"individual_amounts"`
}

func PostKeysignFailure(
	bridge shareTypes.Bridge,
	tx stypes.TxOutItem,
	logger zerolog.Logger,
	thorchainHeight int64,
	utxoErr error,
) error {
	// // PostKeysignFailure only once per SignTx, to not broadcast duplicate messages.
	// var keysignError tss.KeysignError
	// if errors.As(utxoErr, &keysignError) {
	// 	if len(keysignError.Blame.BlameNodes) == 0 {
	// 		// TSS doesn't know which node to blame
	// 		utxoErr = multierror.Append(utxoErr, fmt.Errorf("fail to sign UTXO"))
	// 		return fmt.Errorf("fail to sign the message: %w", utxoErr)
	// 	}

	// 	// key sign error forward the keysign blame to thorchain
	// 	txID, err := bridge.PostKeysignFailure(keysignError.Blame, thorchainHeight, tx.Memo, tx.Coins, tx.VaultPubKey)
	// 	if err != nil {
	// 		logger.Error().Err(err).Msg("fail to post keysign failure to thorchain")
	// 		utxoErr = multierror.Append(utxoErr, fmt.Errorf("fail to post keysign failure to THORChain: %w", err))
	// 		return fmt.Errorf("fail to sign the message: %w", utxoErr)
	// 	}
	// 	logger.Info().Str("tx_id", txID).Msgf("post keysign failure to thorchain")
	// }
	// return utxoErr

	return nil
}
