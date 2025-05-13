package tcysmartcontract

import (
	"github.com/mapprotocol/compass-tss/common"
)

func IsTCYSmartContractAddress(address common.Address) bool {
	return address.String() == TCYSmartContractAddress
}

func GetTCYSmartContractAddress() (common.Address, error) {
	return common.NewAddress(TCYSmartContractAddress)
}
