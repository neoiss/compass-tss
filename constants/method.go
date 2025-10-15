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
	RelaySigned = "relaySigned"
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
	GetNetworkFee     = "getNetworkFee"
	VoteTxIn          = "voteTxIn"
	VoteTxOut         = "voteTxOut"
	IsOrderExecuted   = "isOrderExecuted"
	GetTSSStatus      = "getTSSStatus"
)

// -----------------------------------------------------------------
// Method of GasService
// -----------------------------------------------------------------
const (
	GetNetworkFeeInfo = "getNetworkFeeInfo"
)

// -----------------------------------------------------------------
// Method of Gateway
// -----------------------------------------------------------------

const (
	BridgeIn = "bridgeIn"
)

var ZeroHash = "0x0000000000000000000000000000000000000000"
