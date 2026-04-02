package memo_test

import (
	"testing"

	"github.com/mapprotocol/compass-tss/x/memo"
)

func TestParseMemo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		memo    string
		want    memo.Memo
		wantErr bool
	}{
		// {
		// 	name: "add liquidity memo",
		// 	memo: "M+|0x4a4f0d7d412f1d47fa45c434cecf05f2f8a434f7",
		// 	want: memo.AddLiquidityMemo{
		// 		MemoBase: memo.MemoBase{
		// 			TxType: memo.TxAdd,
		// 		},
		// 		Receiver: "0x4a4f0d7d412f1d47fa45c434cecf05f2f8a434f7",
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "txInbound memo",
			memo: "M>|Mapo|0xe72673a3fc50d39e83cffd4815c82b8192f5fe99027261d7c0d84ec10be06513",
			want: memo.InboundMemo{
				MemoBase: memo.MemoBase{
					TxType: memo.TxInbound,
				},
				Chain:   "Mapo",
				OrderID: "0xe72673a3fc50d39e83cffd4815c82b8192f5fe99027261d7c0d84ec10be06513",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := memo.ParseMemo(tt.memo)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseMemo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseMemo() succeeded unexpectedly")
			}
			// // TODO: update the condition below to compare got with tt.want.
			// if true {
			// 	t.Errorf("ParseMemo() = %v", got)
			// }
			t.Log("chain ", got.GetChain())
			t.Log("chain ", got.GetChain())
		})
	}
}
