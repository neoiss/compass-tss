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
// Method of relay
// -----------------------------------------------------------------
const (
	RelaySigned       = "relaySigned"
	BridgeRelay       = "BridgeRelay"
	BridgeCompleted   = "BridgeCompleted"
	BridgeRelaySigned = "BridgeRelaySigned"
)

// -----------------------------------------------------------------
// Method of tssManager
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
// Method of Gateway
// -----------------------------------------------------------------
const (
	TransferAllowance = "transferAllowance"
	TransferOut       = "transferOut"
	TransferOutCall   = "transferOutCall"
)
