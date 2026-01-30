package common

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_tronAddressToBytes2(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test1",
			args: args{
				address: "TQwhbNxVL4WXJCzp7kpDCw8T1VRi6vuWuz",
			},
			want: common.Hex2Bytes("A440ec08b651d04f11EE538a30e691578e3FD983"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "test2",
			args: args{
				address: "TQwhbNxVL4WXJCzp7kpdCw8T1VRi6vuWuz",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid tron address checksum")
			},
		},
		{
			name: "test3",
			args: args{
				address: "TQwhbNxVL4WXJCzp7kpDCw8T1VRi6vuWuzTQwhbNxVL4WXJCzp7kpDCw8T1VRi6vuWuz",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid tron address length")
			},
		},
		{
			name: "test4",
			args: args{
				address: "TQwhbNxVL4WXJCzp7kpDCw8T1VRi6vuWuz",
			},
			want: common.Hex2Bytes("A440ec08b651d04f11EE538a30e691578e3FD983"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "test5",
			args: args{
				address: "TXcb8NicbbiT1sfSuNRZH19XggX1ph3Aoz",
			},
			want: common.Hex2Bytes("ED6C808451AEEB4B0399788984D2E08D0ECB3B3C"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tronAddressToBytes(tt.args.address)
			t.Log("want: ", tt.want)
			t.Log("got: ", got)
			t.Log("addr: ", common.Bytes2Hex(got))
			if !tt.wantErr(t, err, fmt.Sprintf("tronAddressToBytes(%v)", tt.args.address)) {
				t.Log("error: ", err)
				return
			}
			assert.Equalf(t, tt.want, got, "tronAddressToBytes(%v)", tt.args.address)
		})
	}
}

// 0xA440ec08b651d04f11EE538a30e691578e3FD983
// 41A440EC08B651D04F11EE538A30E691578E3FD983

func Test_solanaAddressToBytes(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test1",
			args: args{
				address: "",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid solana address length")
			},
		},
		{
			name: "test2",
			args: args{
				address: "3u6qZdDUjmdAiZSx2Khq5y",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid solana address length")
			},
		},
		{
			name: "test3",
			args: args{
				address: "3u6qZdDUjmdAiZSx2Khq5yoTqW9XVQHUNpLcqCpbhdei3u6qZdDUjmdAiZSx2Khq5yoTqW9XVQHUNpLcqCpbhdei",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid solana address length")
			},
		},
		{
			name: "test4",
			args: args{
				address: "3u6qZdDUjmdAiZSx2Khq5yoTqW9XVQHUNpLcqCpbhdei",
			},
			want: base58.Decode("3u6qZdDUjmdAiZSx2Khq5yoTqW9XVQHUNpLcqCpbhdei"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := solanaAddressToBytes(tt.args.address)
			if !tt.wantErr(t, err, fmt.Sprintf("solanaAddressToBytes(%v)", tt.args.address)) {
				return
			}
			assert.Equalf(t, tt.want, got, "solanaAddressToBytes(%v)", tt.args.address)
		})
	}
}

func Test_evmAddressToBytes(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test1",
			args: args{
				address: "0x0000000000000000000000000000000000000000",
			},
			want: common.Hex2Bytes("0000000000000000000000000000000000000000"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "test2",
			args: args{
				address: "",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "empty evm address")
			},
		},
		{
			name: "test3",
			args: args{
				address: "0x",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "empty evm address")
			},
		},
		{
			name: "test4",
			args: args{
				address: "0x0hhhh",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid evm address: 0x0hhhh")
			},
		},
		{
			name: "test5",
			args: args{
				address: "0x0eb16a9cfdf8e3a4471ef190ee63de5a24f387",
			},
			want: []byte{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "invalid evm address: 0x0eb16a9cfdf8e3a4471ef190ee63de5a24f387")
			},
		},
		{
			name: "test5",
			args: args{
				address: "0x0eb16a9cfdf8e3a4471ef190ee63de5a24f38787",
			},
			want: common.Hex2Bytes("0eb16a9cfdf8e3a4471ef190ee63de5a24f38787"),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evmAddressToBytes(tt.args.address)
			t.Log("want: ", tt.want)
			t.Log("got: ", got)
			t.Log("addr: ", common.Bytes2Hex(got))
			if !tt.wantErr(t, err, fmt.Sprintf("evmAddressToBytes(%v)", tt.args.address)) {
				return
			}
			assert.Equalf(t, tt.want, got, "evmAddressToBytes(%v)", tt.args.address)
		})
	}
}
