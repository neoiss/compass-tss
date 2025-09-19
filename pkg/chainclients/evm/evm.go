package evm

import _ "embed"

////go:embed abi/router.json
//var routerContractABI string

//go:embed abi/erc20.json
var erc20ContractABI string

//go:embed abi/gateway.json
var gatewayContractABI string

const (
	defaultDecimals = 18 // evm chains consolidate all decimals to 18 (wei)
	tenGwei         = 10000000000
)
