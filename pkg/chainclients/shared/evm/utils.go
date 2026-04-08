package evm

import (
	"fmt"

	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// IsSmartContractCall - determine if the transaction is a smart contract call and thus should
// be parsed using the SmartContractLogParser - these txs may have a DepositEvent from
// the THORChain Router. This is determined by checking if the tx data is at least 4
// bytes and there is at least one log in the receipt. It is possible for a smart
// contract call to have no logs or no data, but these cannot be THORChain deposits, so
// they can be parsed as a normal tx. On the other hand, simple ETH/ERC20 transfer CAN
// have logs & data, but again these will not be THORChain deposits or outbounds.
func IsSmartContractCall(tx *etypes.Transaction, receipt *etypes.Receipt) bool {
	data := tx.Data()
	if len(data) < 4 {
		return false
	}
	if len(receipt.Logs) == 0 {
		return false
	}
	return true
}

func GenerateOrderID(chainID, txHash string) ecommon.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("%s%s", chainID, txHash)))
}
