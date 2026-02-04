package xrp

import (
	"fmt"
	"math/big"
	"strconv"

	sdkmath "cosmossdk.io/math"

	"github.com/mapprotocol/compass-tss/common"

	txtypes "github.com/Peersyst/xrpl-go/xrpl/transaction/types"
)

func parseCurrencyAmount(coin txtypes.CurrencyAmount) (*big.Int, error) {
	if xrpAmount, ok := coin.(txtypes.XRPCurrencyAmount); ok {
		return big.NewInt(int64(xrpAmount.Uint64())), nil
	}
	if issuedAmount, ok := coin.(txtypes.IssuedCurrencyAmount); ok {
		amount, err := strconv.ParseInt(issuedAmount.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return big.NewInt(amount), nil
	}
	return nil, fmt.Errorf("invalid xrp currency type")
}

func convertFee(coin txtypes.CurrencyAmount) (sdkmath.Uint, error) {
	_, exists := GetAssetByXrpCurrency(coin)
	if !exists {
		return sdkmath.NewUint(0), fmt.Errorf("asset does not exist / not whitelisted by client")
	}

	amount, err := parseCurrencyAmount(coin)
	if err != nil {
		return sdkmath.NewUint(0), err
	}
	return sdkmath.NewUintFromBigInt(amount), nil
}

func fromThorchainToXrp(coin common.Coin) (txtypes.CurrencyAmount, error) {
	asset, exists := GetAssetByThorchainAsset(coin.Asset)
	if !exists {
		return nil, fmt.Errorf("asset (%s) does not exist / not whitelisted by client", coin.Asset)
	}

	decimals := asset.XrpDecimals
	amount := coin.Amount.BigInt()
	var exp big.Int
	if decimals > common.THORChainDecimals {
		// Decimals are more than native relay, so multiply...
		decimalDiff := decimals - common.THORChainDecimals
		amount.Mul(amount, exp.Exp(big.NewInt(10), big.NewInt(decimalDiff), nil))
	} else if decimals < common.THORChainDecimals {
		// Decimals are less than native relay, so divide...
		decimalDiff := common.THORChainDecimals - decimals
		amount.Quo(amount, exp.Exp(big.NewInt(10), big.NewInt(decimalDiff), nil))
	}

	if asset.XrpKind == txtypes.ISSUED {
		return txtypes.IssuedCurrencyAmount{
			Issuer:   txtypes.Address(asset.XrpIssuer),
			Currency: asset.XrpCurrency,
			Value:    amount.String(),
		}, nil
	}

	return txtypes.XRPCurrencyAmount(amount.Uint64()), nil
}
