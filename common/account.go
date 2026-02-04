package common

import "github.com/mapprotocol/compass-tss/common/cosmos"

// Account define a struct to hold account information across all chain
type Account struct {
	Sequence      int64
	AccountNumber int64
	Balance       cosmos.Uint
	HasMemoFlag   bool
}

// NewAccount create a new instance of Account
func NewAccount(sequence, accountNumber int64, balance cosmos.Uint, hasMemoFlag bool) Account {
	return Account{
		Sequence:      sequence,
		AccountNumber: accountNumber,
		Balance:       balance,
		HasMemoFlag:   hasMemoFlag,
	}
}
