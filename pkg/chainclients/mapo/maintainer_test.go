package mapo

import (
	"math/big"
	"testing"

	ecommon "github.com/ethereum/go-ethereum/common"
)

func TestBridge_genHash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		epoch   *big.Int
		members []ecommon.Address
		want    ecommon.Hash
		wantErr bool
	}{
		// {
		// 	name:  "test1",
		// 	epoch: big.NewInt(123),
		// 	members: []ecommon.Address{
		// 		ecommon.HexToAddress("0x69a99844D11bea5c6b73c84166e3c6B62Cd870F5"),
		// 		ecommon.HexToAddress("0x2b7588165556aB2fA1d30c520491C385BAa424d8"),
		// 		ecommon.HexToAddress("0xad76db9c043fB5386D8D5C4634F55bbAda559B29"),
		// 		ecommon.HexToAddress("0xE796bc0Ef665D5F730408a55AA0FF4e6f8B90920"),
		// 	},
		// 	want:    ecommon.HexToHash("0x4f57ce20df6123d20d52727e5f3d8c58255c2cc7d519b609eafb61aee42ed6ed"),
		// 	wantErr: false,
		// },
		{
			name:  "test1",
			epoch: big.NewInt(123),
			members: []ecommon.Address{
				ecommon.HexToAddress("0x2b7588165556aB2fA1d30c520491C385BAa424d8"),
				ecommon.HexToAddress("0xad76db9c043fB5386D8D5C4634F55bbAda559B29"),
				ecommon.HexToAddress("0x25fa71d4f689f4b65eb6d020a414090828281d51"),
			},
			want:    ecommon.HexToHash("0x5fcd379e0491f2af891592809c85478334d682556d599ec5dde8c690c7727024"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bridge
			got, gotErr := b.genHash(tt.epoch, tt.members, 100)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("genHash() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("genHash() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("genHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
