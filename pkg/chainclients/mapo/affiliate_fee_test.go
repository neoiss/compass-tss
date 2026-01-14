package mapo

import (
	"errors"
	"testing"
)

func TestBridge_GetAffiliateIDByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
		err     error
	}{
		{
			name: "t-1",
			args: args{
				name: "buttertest",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "t-1",
			args: args{
				name: "aabb",
			},
			want:    0,
			wantErr: true,
			err:     errors.New("affiliate not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bridge := getBridgeForTest(t)
			got, err := bridge.GetAffiliateIDByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAffiliateIDByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("GetAffiliateIDByName() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("GetAffiliateIDByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBridge_GetAffiliateIDByAlias(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
		err     error
	}{
		{
			name: "t-1",
			args: args{
				name: "b",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "t-1",
			args: args{
				name: "aabb",
			},
			want:    0,
			wantErr: true,
			err:     errors.New("affiliate not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bridge := getBridgeForTest(t)
			got, err := bridge.GetAffiliateIDByAlias(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAffiliateIDByAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("GetAffiliateIDByAlias() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("GetAffiliateIDByAlias() got = %v, want %v", got, tt.want)
			}
		})
	}
}
