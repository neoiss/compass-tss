package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/protobuf/proto"

	"github.com/mapprotocol/compass-tss/common"
	cosmos "github.com/mapprotocol/compass-tss/common/cosmos"
	"gitlab.com/thorchain/thornode/v3/api/types"
)

var (
	_ sdk.Msg              = &MsgNetworkFee{}
	_ sdk.HasValidateBasic = &MsgNetworkFee{}
	_ sdk.LegacyMsg        = &MsgNetworkFee{}
)

// NewMsgNetworkFee create a new instance of MsgNetworkFee
func NewMsgNetworkFee(blockHeight int64, chain common.Chain, transactionSize, transactionFeeRate uint64, signer cosmos.AccAddress) *MsgNetworkFee {
	return &MsgNetworkFee{
		BlockHeight:        blockHeight,
		Chain:              chain,
		TransactionSize:    transactionSize,
		TransactionFeeRate: transactionFeeRate,
		Signer:             signer,
	}
}

// ValidateBasic implements HasValidateBasic
// ValidateBasic is now ran in the message service router handler for messages that
// used to be routed using the external handler and only when HasValidateBasic is implemented.
// No versioning is used there.
func (m *MsgNetworkFee) ValidateBasic() error {
	if m.BlockHeight <= 0 {
		return cosmos.ErrUnknownRequest("block height can't be negative, or zero")
	}
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress(m.Signer.String())
	}
	if m.Chain.IsEmpty() {
		return cosmos.ErrUnknownRequest("chain can't be empty")
	}
	if m.TransactionSize <= 0 {
		return cosmos.ErrUnknownRequest("invalid transaction size")
	}
	if m.TransactionFeeRate <= 0 {
		return cosmos.ErrUnknownRequest("invalid transaction fee rate")
	}
	return nil
}

// GetSigners defines whose signature is required
func (m *MsgNetworkFee) GetSigners() []cosmos.AccAddress {
	return []cosmos.AccAddress{m.Signer}
}

func MsgNetworkFeeCustomGetSigners(m proto.Message) ([][]byte, error) {
	msg, ok := m.(*types.MsgNetworkFee)
	if !ok {
		return nil, fmt.Errorf("can't cast as MsgNetworkFee: %T", m)
	}
	return [][]byte{msg.Signer}, nil
}
