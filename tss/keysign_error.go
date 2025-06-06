package tss

import (
	"fmt"

	tssMessages "github.com/mapprotocol/compass-tss/p2p/messages"

	"github.com/mapprotocol/compass-tss/x/types"
)

// KeysignError is a custom error create to include which party to blame
type KeysignError struct {
	Blame types.Blame
}

// NewKeysignError create a new instance of KeysignError
func NewKeysignError(blame types.Blame) KeysignError {
	return KeysignError{
		Blame: blame,
	}
}

// Error implement error interface
func (k KeysignError) Error() string {
	return fmt.Sprintf("fail to complete TSS keysign, reason:%s, round:%s, culprit:%+v", k.Blame.FailReason, k.Blame.Round, k.Blame.BlameNodes)
}

func (k KeysignError) IsRound7() bool {
	return k.Blame.Round == tssMessages.KEYSIGN7
}
