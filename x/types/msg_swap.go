package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/common/cosmos"
)

// MaxAffiliateFeeBasisPoints basis points for withdrawals
const MaxAffiliateFeeBasisPoints = 1_000

var (
	_ sdk.Msg              = &MsgSwap{}
	_ sdk.HasValidateBasic = &MsgSwap{}
	_ sdk.LegacyMsg        = &MsgSwap{}
)

// NewMsgSwap is a constructor function for MsgSwap
func NewMsgSwap(tx common.Tx, target common.Asset, destination common.Address, tradeTarget cosmos.Uint, affAddr common.Address, affPts cosmos.Uint, agg, aggregatorTargetAddr string, aggregatorTargetLimit *cosmos.Uint, otype OrderType, quan, interval uint64, signer cosmos.AccAddress) *MsgSwap {
	return &MsgSwap{
		Tx:                      tx,
		TargetAsset:             target,
		Destination:             destination,
		TradeTarget:             tradeTarget,
		AffiliateAddress:        affAddr,
		AffiliateBasisPoints:    affPts,
		Signer:                  signer,
		Aggregator:              agg,
		AggregatorTargetAddress: aggregatorTargetAddr,
		AggregatorTargetLimit:   aggregatorTargetLimit,
		OrderType:               otype,
		StreamQuantity:          quan,
		StreamInterval:          interval,
	}
}

func (m *MsgSwap) IsStreaming() bool {
	return m.StreamInterval > 0
}

func (m *MsgSwap) GetStreamingSwap() StreamingSwap {
	return NewStreamingSwap(
		m.Tx.ID,
		m.StreamQuantity,
		m.StreamInterval,
		m.TradeTarget,
		m.Tx.Coins[0].Amount,
	)
}

// ValidateBasic runs stateless checks on the message
func (m *MsgSwap) ValidateBasic() error {
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress(m.Signer.String())
	}
	if err := m.Tx.Valid(); err != nil {
		return cosmos.ErrUnknownRequest(err.Error())
	}
	if m.TargetAsset.IsEmpty() {
		return cosmos.ErrUnknownRequest("swap Target cannot be empty")
	}
	if len(m.Tx.Coins) > 1 {
		return cosmos.ErrUnknownRequest("not expecting multiple coins in a swap")
	}
	if m.Tx.Coins.IsEmpty() {
		return cosmos.ErrUnknownRequest("swap coin cannot be empty")
	}
	for _, coin := range m.Tx.Coins {
		if coin.Asset.Equals(m.TargetAsset) {
			return cosmos.ErrUnknownRequest("swap Source and Target cannot be the same.")
		}
	}
	if m.Destination.IsEmpty() {
		return cosmos.ErrUnknownRequest("swap Destination cannot be empty")
	}
	// TODO: remove this check on hardfork
	if m.AffiliateAddress.IsEmpty() && !m.AffiliateBasisPoints.IsZero() {
		return cosmos.ErrUnknownRequest("swap affiliate address is empty while affiliate basis points is non-zero")
	}
	if !m.AffiliateBasisPoints.IsZero() && m.AffiliateBasisPoints.GT(cosmos.NewUint(MaxAffiliateFeeBasisPoints)) {
		return cosmos.ErrUnknownRequest(fmt.Sprintf("affiliate fee basis points can't be more than %d", MaxAffiliateFeeBasisPoints))
	}
	if !m.Destination.IsNoop() && !m.Destination.IsChain(m.TargetAsset.GetChain()) {
		return cosmos.ErrUnknownRequest("swap destination address is not the same chain as the target asset")
	}
	if !m.AffiliateAddress.IsEmpty() && !m.AffiliateAddress.IsChain(common.THORChain) {
		return cosmos.ErrUnknownRequest("swap affiliate address must be a THOR address")
	}
	if len(m.Aggregator) != 0 && len(m.AggregatorTargetAddress) == 0 {
		return cosmos.ErrUnknownRequest("aggregator target asset address is empty")
	}
	if len(m.AggregatorTargetAddress) > 0 && len(m.Aggregator) == 0 {
		return cosmos.ErrUnknownRequest("aggregator is empty")
	}
	return nil
}

// GetSigners defines whose signature is required
func (m *MsgSwap) GetSigners() []cosmos.AccAddress {
	return []cosmos.AccAddress{m.Signer}
}

func (m *MsgSwap) GetTotalAffiliateFee() cosmos.Uint {
	return common.GetSafeShare(
		m.AffiliateBasisPoints,
		cosmos.NewUint(10000),
		m.Tx.Coins[0].Amount,
	)
}
