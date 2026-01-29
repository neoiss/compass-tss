package evm_test

import (
	"testing"
	"time"

	ethclient "github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/evm"
)

func TestEthRPC_GetBlockSafe(t *testing.T) {
	ethClient, err := ethclient.Dial("https://cf-rpc.omniservice.dev/filter/arbitrum")
	if err != nil {
		t.Fatalf("Dial faild, err is %v", err)
	}
	tests := []struct {
		name    string
		client  *ethclient.Client
		timeout time.Duration
		chain   string
		number  int64
		want    *evm.Block
		wantErr bool
	}{
		{
			name:    "success",
			client:  ethClient,
			timeout: time.Second * 10,
			chain:   "arb",
			number:  425674957,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := evm.NewEthRPC(tt.client, tt.timeout, tt.chain)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := e.GetBlockSafe(tt.number)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBlockSafe() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBlockSafe() succeeded unexpectedly")
			}

			// got1, gotErr := e.GetBlock(tt.number)
			// if gotErr != nil {
			// 	if !tt.wantErr {
			// 		t.Errorf("GetBlockSafe() failed: %v", gotErr)
			// 	}
			// 	return
			// }

			for _, ele := range got.Transactions {
				gas, ok := ele.GetGasPrice()
				// if gas.Cmp(got1.Transactions()[idx].GasPrice()) == 0 {
				// 	continue
				// }
				t.Log("txType-", ele.Type, "gas ----- ", gas, "ok", ok) //, "other ===", got1.Transactions()[idx].GasPrice())
			}
		})
	}
}
