package mapo

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBridge_GetChainIDFromFusionReceiver(t *testing.T) {

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				"Uni",
			},
			want:    big.NewInt(130),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBridgeForTest(t)
			got, err := b.GetChainIDFromFusionReceiver(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainIDFromFusionReceiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokenDecimals() got = %v, want %v", got, tt.want)
			}
		})
	}
}
