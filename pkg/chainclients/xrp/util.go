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

func decimalToXrp(amount *big.Int) (txtypes.CurrencyAmount, error) {
	asset, exists := GetAssetByThorchainAsset(common.XRPAsset)
	if !exists {
		return nil, fmt.Errorf("asset (xrp) does not exist / not whitelisted by client")
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
