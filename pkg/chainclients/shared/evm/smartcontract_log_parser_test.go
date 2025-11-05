package evm

import (
	"testing"

	ecommon "github.com/ethereum/go-ethereum/common"
)

func Test_parseChainAndGasLimit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cgl     ecommon.Hash
		want    *ChainAndGasLimit
		wantErr bool
	}{
		{
			name: "valid",
			cgl:  ecommon.HexToHash("0004d500000000020004d500000000020000000008f0d18000000000000186a0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parseChainAndGasLimit(tt.cgl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseChainAndGasLimit() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseChainAndGasLimit() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("parseChainAndGasLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
