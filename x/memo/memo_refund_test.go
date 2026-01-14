package thorchain

import (
	"testing"
)

func TestRefundMemo(t *testing.T) {
	type args struct {
		chain   string
		orderID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t-1",
			args: args{
				chain:   "btc_test",
				orderID: "0x8fa90e636cbc1dd79691e255fae3702205031798d1ca393f5042fddc64c1a122",
			},
			want: "M<|btc_test|0x8fa90e636cbc1dd79691e255fae3702205031798d1ca393f5042fddc64c1a122",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRefundMemo(tt.args.chain, tt.args.orderID); got.String() != tt.want {
				t.Errorf("NewRefundMemo() = %v, want %v", got, tt.want)
			}
		})
	}
}
