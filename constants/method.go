package constants

// -----------------------------------------------------------------
// Method of Maintainer
// -----------------------------------------------------------------

const (
	ElectionEpoch      = "electionEpoch"
	GetEpochInfo       = "getEpochInfo"
	GetMaintainerInfos = "getMaintainerInfos"
	Register           = "register"
	Heartbeat          = "heartbeat"
	VersionMethod      = "version"
)

// -----------------------------------------------------------------
// Method of Relay
// -----------------------------------------------------------------
const (
	RelaySigned     = "relaySigned"
	IsOrderExecuted = "isOrderExecuted"
	OrderInfos      = "orderInfos"
	Completed       = "completed" // custom method
)

// -----------------------------------------------------------------
// EventName of Relay
// -----------------------------------------------------------------
const (
	BridgeRelay       = "BridgeRelay"
	BridgeCompleted   = "BridgeCompleted"
	BridgeRelaySigned = "BridgeRelaySigned"
)

// -----------------------------------------------------------------
// Method of TssManager
// -----------------------------------------------------------------
const (
	VoteUpdateTssPool = "voteUpdateTssPool"
	GetKeyShare       = "getKeyShare"
	VoteNetworkFee    = "voteNetworkFee"
	VoteTxIn          = "voteTxIn"
	VoteTxOut         = "voteTxOut"
	GetTSSStatus      = "getTSSStatus"
)

// -----------------------------------------------------------------
// Method of Gateway
// -----------------------------------------------------------------

const (
	BridgeIn = "bridgeIn"
)

var ZeroHash = "0x0000000000000000000000000000000000000000"

// -----------------------------------------------------------------
// Method of View
// -----------------------------------------------------------------
const (
	GetPublickeys             = "getPublicKeys" // vaults
	GetLastTxOutHeight        = "getLastTxOutHeight"
	GetLastTxInHeight         = "getLastTxInHeight"
	GetVault                  = "getVault"
	GetNetworkFeeInfo         = "getNetworkFeeInfo" // gas
	GetInfoByNickname         = "getInfoByNickname" // affiliate related
	GetInfoByShortName        = "getInfoByShortName"
	GetChainByName            = "getChainByName" // token
	GetChainName              = "getChainName"
	GetTokenAddressByNickname = "getTokenAddressByNickname"
	GetTokenDecimals          = "getTokenDecimals"
)

// -----------------------------------------------------------------
// Method of configration
// -----------------------------------------------------------------

const (
	GetIntValue      = "getIntValue"
	GetAddressValue  = "getAddressValue"
	GetBoolValue     = "getBoolValue"
	GetStringValue   = "getStringValue"
	GetBytesValue    = "getBytesValue"
	BatchGetIntValue = "batchGetIntValue"
)
