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
	_ sdk.Msg              = &MsgErrataTxQuorum{}
	_ sdk.HasValidateBasic = &MsgErrataTxQuorum{}
	_ sdk.LegacyMsg        = &MsgErrataTxQuorum{}
)

// NewMsgErrataTxQuorum is a constructor function for MsgErrataTxQuorum
func NewMsgErrataTxQuorum(tx *common.QuorumErrataTx, signer cosmos.AccAddress) *MsgErrataTxQuorum {
	return &MsgErrataTxQuorum{
		QuoErrata: tx,
		Signer:    signer,
	}
}

// ValidateBasic implements HasValidateBasic
// ValidateBasic is now ran in the message service router handler for messages that
// used to be routed using the external handler and only when HasValidateBasic is implemented.
// No versioning is used there.
func (m *MsgErrataTxQuorum) ValidateBasic() error {
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress(m.Signer.String())
	}
	tx := m.QuoErrata.ErrataTx
	if tx.Id.IsEmpty() {
		return cosmos.ErrUnknownRequest("Tx ID cannot be empty")
	}
	if tx.Chain.IsEmpty() {
		return cosmos.ErrUnknownRequest("chain cannot be empty")
	}

	return nil
}

// GetSigners defines whose signature is required
func (m *MsgErrataTxQuorum) GetSigners() []cosmos.AccAddress {
	return quorumSignersCommon(m.QuoErrata.Attestations)
}

func MsgErrataTxQuorumCustomGetSigners(m proto.Message) ([][]byte, error) {
	msg, ok := m.(*types.MsgErrataTxQuorum)
	if !ok {
		return nil, fmt.Errorf("can't cast as MsgErrataTxQuorum: %T", m)
	}

	return quorumSignersApiCommon(msg.QuoErrata.Attestations), nil
}
