package mapo

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func (b *Bridge) Stop() {
	b.ethClient.Close()
	close(b.stopChan)
}

// GetGasPrice gets gas price from eth scanner
func (b *Bridge) GetGasPrice() *big.Int {
	return b.gasPrice
}

// GetConstants from thornode
func (b *Bridge) GetConstants() (map[string]int64, error) {
	var constantStr = `{
		"int_64_values": {
			"AllowWideBlame": 0,
			"AsgardSize": 40,
			"BadValidatorRedline": 3,
			"BankSendEnabled": 0,
			"BlocksPerYear": 5256000,
			"BondSlashBan": 500000000000,
			"ChurnInterval": 43200,
			"ChurnMigrateRounds": 5,
			"ChurnOutForLowVersionBlocks": 21600,
			"ChurnRetryInterval": 720,
			"CloutLimit": 0,
			"CloutReset": 720,
			"DerivedDepthBasisPts": 0,
			"DerivedMinDepth": 100,
			"DerivedSlipMinBps": 0,
			"DesiredValidatorSet": 100,
			"DevFundSystemIncomeBps": 500,
			"DoubleBlockSignSlashPoints": 1000,
			"DoubleSignMaxAge": 24,
			"DynamicMaxAnchorCalcInterval": 14400,
			"DynamicMaxAnchorSlipBlocks": 201600,
			"DynamicMaxAnchorTarget": 0,
			"EVMDisableContractWhitelist": 0,
			"EmissionCurve": 6,
			"EnableAdvSwapQueue": 0,
			"EnableDerivedAssets": 0,
			"EnableOrderBooks": 0,
			"EnableUSDFees": 0,
			"FailKeygenSlashPoints": 720,
			"FailKeysignSlashPoints": 2,
			"FeeUSDRoundSignificantDigits": 2,
			"FundMigrationInterval": 360,
			"JailTimeKeygen": 4320,
			"JailTimeKeysign": 60,
			"KeygenRetryInterval": 0,
			"L1SlipMinBps": 0,
			"LackOfObservationPenalty": 2,
			"LendingLever": 3333,
			"LiquidityLockUpBlocks": 0,
			"LoanRepaymentMaturity": 0,
			"LoanStreamingSwapsInterval": 0,
			"MaxAffiliateFeeBasisPoints": 10000,
			"MaxAnchorBlocks": 300,
			"MaxAnchorSlip": 1500,
			"MaxAvailablePools": 100,
			"MaxBondProviders": 6,
			"MaxCR": 60000,
			"MaxMissingBlockChurnOut": 0,
			"MaxNodeToChurnOutForLowVersion": 1,
			"MaxOutboundAttempts": 0,
			"MaxOutboundFeeMultiplierBasisPoints": 30000,
			"MaxRuneSupply": -1,
			"MaxSwapsPerBlock": 100,
			"MaxSynthPerPoolDepth": 1700,
			"MaxSynthsForSaversYield": 0,
			"MaxTrackMissingBlock": 700,
			"MaxTxOutOffset": 720,
			"MigrationVaultSecurityBps": 0,
			"MinCR": 10000,
			"MinOutboundFeeMultiplierBasisPoints": 15000,
			"MinRuneForTCYStakeDistribution": 210000000000,
			"MinRunePoolDepth": 1000000000000,
			"MinSlashPointsForBadValidator": 100,
			"MinSwapsPerBlock": 10,
			"MinTCYForTCYStakeDistribution": 100000,
			"MinTxOutVolumeThreshold": 100000000000,
			"MinimumBondInRune": 100000000000000,
			"MinimumL1OutboundFeeUSD": 1000000,
			"MinimumNodesForBFT": 4,
			"MinimumPoolLiquidityFee": 0,
			"MissBlockSignSlashPoints": 1,
			"MissingBlockChurnOut": 0,
			"MultipleAffiliatesMaxCount": 5,
			"NativeOutboundFeeUSD": 2000000,
			"NativeTransactionFee": 2000000,
			"NativeTransactionFeeUSD": 2000000,
			"NodeOperatorFee": 500,
			"NodePauseChainBlocks": 720,
			"ObservationDelayFlexibility": 10,
			"ObserveSlashPoints": 1,
			"OperationalVotesMin": 3,
			"OutboundTransactionFee": 2000000,
			"POLBuffer": 0,
			"POLMaxNetworkDeposit": 0,
			"POLMaxPoolMovement": 100,
			"POLTargetSynthPerPoolDepth": 0,
			"PauseBond": 0,
			"PauseLoans": 1,
			"PauseOnSlashThreshold": 10000000000,
			"PauseUnbond": 0,
			"PendingLiquidityAgeLimit": 100800,
			"PendulumAssetsBasisPoints": 10000,
			"PendulumUseEffectiveSecurity": 0,
			"PendulumUseVaultAssets": 0,
			"PermittedSolvencyGap": 100,
			"PoolCycle": 43200,
			"PreferredAssetOutboundFeeMultiplier": 100,
			"RUNEPoolDepositMaturityBlocks": 1296000,
			"RUNEPoolEnabled": 0,
			"RUNEPoolHaltDeposit": 0,
			"RUNEPoolHaltWithdraw": 0,
			"RUNEPoolMaxReserveBackstop": 500000000000000,
			"RagnarokProcessNumOfLPPerIteration": 200,
			"RescheduleCoalesceBlocks": 0,
			"SaversEjectInterval": 0,
			"SaversStreamingSwapsInterval": 0,
			"SecuredAssetSlipMinBps": 5,
			"SigningTransactionPeriod": 300,
			"SlashPenalty": 15000,
			"StagedPoolCost": 1000000000,
			"StreamingSwapMaxLength": 14400,
			"StreamingSwapMaxLengthNative": 5256000,
			"StreamingSwapMinBPFee": 0,
			"StreamingSwapPause": 0,
			"SynthSlipMinBps": 0,
			"SynthYieldBasisPoints": 5000,
			"SynthYieldCycle": 0,
			"SystemIncomeBurnRateBps": 1,
			"TCYClaimingHalt": 1,
			"TCYClaimingSwapHalt": 1,
			"TCYStakeDistributionHalt": 1,
			"TCYStakeSystemIncomeBps": 1000,
			"TCYStakingHalt": 1,
			"TCYUnstakingHalt": 1,
			"TNSFeeOnSale": 1000,
			"TNSFeePerBlock": 20,
			"TNSFeePerBlockUSD": 20,
			"TNSRegisterFee": 1000000000,
			"TNSRegisterFeeUSD": 1000000000,
			"TVLCapBasisPoints": 0,
			"TargetOutboundFeeSurplusRune": 10000000000000,
			"TradeAccountsDepositEnabled": 1,
			"TradeAccountsEnabled": 0,
			"TradeAccountsSlipMinBps": 0,
			"TxOutDelayMax": 17280,
			"TxOutDelayRate": 2500000000,
			"ValidatorMaxRewardRatio": 1,
			"VirtualMultSynths": 2,
			"VirtualMultSynthsBasisPoints": 10000
		},
		"bool_values": {
			"StrictBondLiquidityRatio": true
		},
		"string_values": {
			"DefaultPoolStatus": "Staged",
			"DevFundAddress": ""
		}
	}`
	var result struct {
		Int64Values map[string]int64 `json:"int_64_values"`
	}
	err := json.Unmarshal([]byte(constantStr), &result)
	if err != nil {
		return nil, err
	}

	return result.Int64Values, nil
}

// GetMapVersion retrieve mapBridge version
// func (b *Bridge) GetMapVersion() (semver.Version, error) {
func (b *Bridge) GetMapVersion() (string, error) {
	// todo handler
	return "1.0.0", nil
}

// GetMimir - get mimir settings
func (b *Bridge) GetMimir(key string) (int64, error) {
	// todo handler
	switch key {
	case "MaxConfirmations-ETH":
		return 14, nil
	case "MaxConfirmations-BSC":
		return -1, nil
	case "SignerConcurrency":
		return 20, nil
	case "HALTSIGNINGBSC", "HALTSIGNINGETH", "HALTSIGNINGBTC", "HALTSIGNING":
		return 0, nil
	}
	return 0, nil
}

// GetMimirWithRef is a helper function to more readably insert references (such as Asset MimirString or Chain) into Mimir key templates.
func (b *Bridge) GetMimirWithRef(template, ref string) (int64, error) {
	// 'template' should be something like "Halt%sChain" (to halt an arbitrary specified chain)
	// or "Ragnarok-%s" (to halt the pool of an arbitrary specified Asset (MimirString used for Assets to join Chain and Symbol with a hyphen).
	key := fmt.Sprintf(template, ref)
	return b.GetMimir(key)
}

func (b *Bridge) SetObserverCache(key, value string) error {

	return nil
}
