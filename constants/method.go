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
// Method of GasService
// -----------------------------------------------------------------
const (
	GetNetworkFeeInfo = "getNetworkFeeInfo"
	GetNetworkFee     = "getNetworkFee"
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
	GetPublickeys      = "getPublicKeys"
	GetLastTxOutHeight = "getLastTxOutHeight"
	GetLastTxInHeight  = "getLastTxInHeight"
	GetVault           = "getVault"
)
