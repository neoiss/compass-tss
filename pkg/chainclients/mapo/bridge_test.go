package mapo

import (
	"context"
	"github.com/ethereum/go-ethereum"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/constants"
	selfAbi "github.com/mapprotocol/compass-tss/pkg/abi"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func Test_Bridge_PostNetworkFee(t *testing.T) {
	ethClient, err := ethclient.Dial("https://testnet-rpc.maplabs.io")
	assert.Nil(t, err)
	pkStr := os.Getenv("pri_key")
	priKey, err := ecrypto.HexToECDSA(pkStr)
	addr := ecommon.HexToAddress("0xad76db9c043fB5386D8D5C4634F55bbAda559B29")
	assert.Nil(t, err)

	ai, err := selfAbi.New(maintainerAbi)
	assert.Nil(t, err)

	to := ecommon.HexToAddress("0x0EdA5e4015448A2283662174DD7def3C3d262D38")

	input, err := ai.PackInput(constants.VoteNetworkFeeOfMaintainer,
		big.NewInt(1),
		big.NewInt(1360095883558914),
		big.NewInt(882082),
		big.NewInt(100000000), // gasPrice
		big.NewInt(1000000),   // gasLimit
		big.NewInt(1500000))   // swapGasLimit
	assert.Nil(t, err)

	head, err := ethClient.HeaderByNumber(context.Background(), nil)
	assert.Nil(t, err)

	gasFeeCap := head.BaseFee

	createdTx := ethereum.CallMsg{
		From:     addr,
		To:       &to,
		GasPrice: gasFeeCap,
		Value:    nil,
		Data:     input,
	}

	t.Log("input ", ecommon.Bytes2Hex(input))
	t.Log("gasFeeCap ", gasFeeCap)

	gasLimit, err := ethClient.EstimateGas(context.Background(), createdTx)
	assert.Nil(t, err)

	nonce, err := ethClient.NonceAt(context.Background(), addr, nil)
	assert.Nil(t, err)

	// create tx
	tipCap := new(big.Int).Mul(gasFeeCap, big.NewInt(10))
	tipCap.Div(tipCap, big.NewInt(100))
	td := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Value:     nil,
		To:        &to,
		Gas:       gasLimit,
		GasTipCap: tipCap,
		GasFeeCap: gasFeeCap,
		Data:      input,
	})

	signedTx, err := types.SignTx(td, types.NewLondonSigner(big.NewInt(212)), priKey)
	assert.Nil(t, err)

	err = ethClient.SendTransaction(context.Background(), signedTx)
	assert.Nil(t, err)

	t.Log("postGasFee tx successfully, tx ================= ", signedTx.Hash().Hex())

}
