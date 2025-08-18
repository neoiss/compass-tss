package types

import (
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

type NetworkFee struct {
	Height              int64    `json:"height,omitempty"`
	ChainId             *big.Int `json:"chain,omitempty"`
	TransactionSize     uint64   `json:"transaction_size,omitempty"`
	TransactionSwapSize uint64   `json:"transaction_swap_size,omitempty"`
	TransactionRate     uint64   `json:"transaction_rate,omitempty"`
}

func (m *NetworkFee) Valid() error {
	if m.ChainId.Uint64() == 0 {
		return errors.New("chain can't be empty")
	}
	if m.Height <= 0 {
		return fmt.Errorf("height can't be zero or negative: %v", m.Height)
	}
	if m.TransactionSize <= 0 {
		return fmt.Errorf("transaction size can't be zero or negative: %v", m.TransactionSize)
	}
	if m.TransactionSwapSize <= 0 {
		return fmt.Errorf("transaction swap size can't be zero or negative: %v", m.TransactionSwapSize)
	}
	if m.TransactionRate <= 0 {
		return errors.New("transaction fee rate can't be zero")
	}
	return nil
}

func (m *NetworkFee) String() string {
	if m == nil {
		return "nil"
	}
	s := strings.Join([]string{`&NetworkFee{`,
		`Height:` + fmt.Sprintf("%v", m.Height) + `,`,
		`ChainId:` + fmt.Sprintf("%v", m.ChainId) + `,`,
		`TransactionSize:` + fmt.Sprintf("%v", m.TransactionSize) + `,`,
		`TransactionSwapSize:` + fmt.Sprintf("%v", m.TransactionSwapSize) + `,`,
		`TransactionRate:` + fmt.Sprintf("%v", m.TransactionRate) + `,`,
		`}`,
	}, "")
	return s
}
