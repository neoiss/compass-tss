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
				chain:   "Btc",
				orderID: "0x9abf316cd5f9ef297569f29a746e78cebd45a0e3c950e0ed3e030be3c73d420c",
			},
			want: "M<|Btc|0x9abf316cd5f9ef297569f29a746e78cebd45a0e3c950e0ed3e030be3c73d420c",
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

func TestParseRefundMemo(t *testing.T) {
	mem, err := ParseMemo("M<|Btc|0x9abf316cd5f9ef297569f29a746e78cebd45a0e3c950e0ed3e030be3c73d420c")
	if err != nil {
		t.Errorf("ParseMemo() error = %v", err)
	}
	t.Log(mem.GetType())
}
