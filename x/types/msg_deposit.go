package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/protobuf/proto"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
	"github.com/mapprotocol/compass-tss/constants"
	"gitlab.com/thorchain/thornode/v3/api/types"
)

var (
	_ sdk.Msg              = &MsgDeposit{}
	_ sdk.HasValidateBasic = &MsgDeposit{}
	_ sdk.LegacyMsg        = &MsgDeposit{}
)

// NewMsgDeposit is a constructor function for NewMsgDeposit
func NewMsgDeposit(coins common.Coins, memo string, signer cosmos.AccAddress) *MsgDeposit {
	return &MsgDeposit{
		Coins:  coins,
		Memo:   memo,
		Signer: signer,
	}
}

// ValidateBasic implements HasValidateBasic
// ValidateBasic is now ran in the message service router for messages that
// used to be routed using the external handler and only when HasValidateBasic is implemented.
// No versioning is used there.
func (m *MsgDeposit) ValidateBasic() error {
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress(m.Signer.String())
	}
	for _, coin := range m.Coins {
		if !coin.IsNative() {
			return cosmos.ErrUnknownRequest("all coins must be native to THORChain")
		}
	}
	if len([]byte(m.Memo)) > constants.MaxMemoSize {
		err := fmt.Errorf("memo must not exceed %d bytes: %d", constants.MaxMemoSize, len([]byte(m.Memo)))
		return cosmos.ErrUnknownRequest(err.Error())
	}

	// deposit only allowed with one coin
	if len(m.Coins) != 1 {
		return errors.New("only one coin is allowed")
	}

	if err := m.Coins[0].Asset.Valid(); err != nil {
		return fmt.Errorf("invalid coin: %w", err)
	}

	return nil
}

// GetSigners defines whose signature is required
func (m *MsgDeposit) GetSigners() []cosmos.AccAddress {
	return []cosmos.AccAddress{m.Signer}
}

func MsgDepositCustomGetSigners(m proto.Message) ([][]byte, error) {
	msgSend, ok := m.(*types.MsgDeposit)
	if !ok {
		return nil, fmt.Errorf("can't cast as MsgDeposit: %T", m)
	}
	return [][]byte{msgSend.Signer}, nil
}
