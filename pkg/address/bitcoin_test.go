package address

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"testing"
)

func TestEncodeBitcoinAddress(t *testing.T) {
	type args struct {
		addr    string
		network *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-P2WPKH",
			args: args{
				addr:    "bc1qwydpn4jtq2cu6d4s3e92qt4t2wahun2292rsqc",
				network: &chaincfg.MainNetParams,
			},
			want:    "0x00711a19d64b02b1cd36b08e4aa02eab53bb7e4d4a",
			wantErr: false,
		},
		{
			name: "t-P2TR",
			args: args{
				addr:    "tb1pa4ae7wfehun4744cckujgesw5a4qjfl76a2vxfxg8cmqssj0d5mq4dxmy7",
				network: &chaincfg.TestNet3Params,
			},
			want:    "0x01ed7b9f3939bf275f56b8c5b924660ea76a0927fed754c324c83e3608424f6d36",
			wantErr: false,
		},
		{
			name: "t-P2PKH",
			args: args{
				addr:    "1Bzn3WN5rgFwCBRXbauKFtTB798yBxk36e",
				network: &chaincfg.MainNetParams,
			},
			want:    "0x01789e915de8e7a6805438e00852960287bbc127dc",
			wantErr: false,
		},
		{
			name: "t-P2SH",
			args: args{
				addr:    "3F9JDXk2JLxVFi95tZMUPcowoPjbQc7Cti",
				network: &chaincfg.MainNetParams,
			},
			want:    "0x0593920d01bbbeb3c649118c498334849df4bc6846",
			wantErr: false,
		},
		{
			name: "t-P2WSH",
			args: args{
				addr:    "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
				network: &chaincfg.MainNetParams,
			},
			want:    "0x00701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := btcutil.DecodeAddress(tt.args.addr, tt.args.network)
			if err != nil {
				t.Fatal(err)
			}

			got, err := EncodeBitcoinAddress(address)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBitcoinAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeBitcoinAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBitcoinAddressToBytes(t *testing.T) {
	type args struct {
		addr    string
		network *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-P2WPKH",
			args: args{
				addr:    "bc1qwydpn4jtq2cu6d4s3e92qt4t2wahun2292rsqc",
				network: &chaincfg.MainNetParams,
			},
			want:    "00711a19d64b02b1cd36b08e4aa02eab53bb7e4d4a",
			wantErr: false,
		},
		{
			name: "t-P2TR",
			args: args{
				addr:    "tb1pa4ae7wfehun4744cckujgesw5a4qjfl76a2vxfxg8cmqssj0d5mq4dxmy7",
				network: &chaincfg.TestNet3Params,
			},
			want:    "01ed7b9f3939bf275f56b8c5b924660ea76a0927fed754c324c83e3608424f6d36",
			wantErr: false,
		},
		{
			name: "t-P2PKH",
			args: args{
				addr:    "1Bzn3WN5rgFwCBRXbauKFtTB798yBxk36e",
				network: &chaincfg.MainNetParams,
			},
			want:    "01789e915de8e7a6805438e00852960287bbc127dc",
			wantErr: false,
		},
		{
			name: "t-P2SH",
			args: args{
				addr:    "3F9JDXk2JLxVFi95tZMUPcowoPjbQc7Cti",
				network: &chaincfg.MainNetParams,
			},
			want:    "0593920d01bbbeb3c649118c498334849df4bc6846",
			wantErr: false,
		},
		{
			name: "t-P2WSH",
			args: args{
				addr:    "3AkndCisp1kXWuJJXS63BZRhPhtkwbYuys",
				network: &chaincfg.MainNetParams,
			},
			want:    "05636fb1c8e782208ba5b9517dd002d35ee3f2ebc7",
			wantErr: false,
		},
		{
			name: "t-1",
			args: args{
				addr:    "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
				network: &chaincfg.MainNetParams,
			},
			want:    "00701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := btcutil.DecodeAddress(tt.args.addr, tt.args.network)
			if err != nil {
				t.Fatal(err)
			}

			got, err := EncodeBitcoinAddressToBytes(address)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBitcoinAddressToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decoded, err := hex.DecodeString(tt.want)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, decoded) {
				t.Errorf("EncodeBitcoinAddressToBytes() got = %v, want %v", got, decoded)
			}
			t.Log(decoded)
		})
	}
}

func TestDecodeBitcoinAddress(t *testing.T) {
	type args struct {
		addr    string
		network *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-P2WPKH",
			args: args{
				addr:    "0x00711a19d64b02b1cd36b08e4aa02eab53bb7e4d4a",
				network: &chaincfg.MainNetParams,
			},
			want:    "bc1qwydpn4jtq2cu6d4s3e92qt4t2wahun2292rsqc",
			wantErr: false,
		},
		{
			name: "t-P2TR",
			args: args{
				addr:    "0x01ed7b9f3939bf275f56b8c5b924660ea76a0927fed754c324c83e3608424f6d36",
				network: &chaincfg.TestNet3Params,
			},
			want:    "tb1pa4ae7wfehun4744cckujgesw5a4qjfl76a2vxfxg8cmqssj0d5mq4dxmy7",
			wantErr: false,
		},
		{
			name: "t-P2TR-without-0x",
			args: args{
				addr:    "01ed7b9f3939bf275f56b8c5b924660ea76a0927fed754c324c83e3608424f6d36",
				network: &chaincfg.TestNet3Params,
			},
			want:    "tb1pa4ae7wfehun4744cckujgesw5a4qjfl76a2vxfxg8cmqssj0d5mq4dxmy7",
			wantErr: false,
		},
		{
			name: "t-P2WSH",
			args: args{
				addr:    "0x05636fb1c8e782208ba5b9517dd002d35ee3f2ebc7",
				network: &chaincfg.MainNetParams,
			},
			want:    "3AkndCisp1kXWuJJXS63BZRhPhtkwbYuys",
			wantErr: false,
		},
		{
			name: "t-P2PKH",
			args: args{
				addr:    "0x01789e915de8e7a6805438e00852960287bbc127dc",
				network: &chaincfg.MainNetParams,
			},
			want:    "1Bzn3WN5rgFwCBRXbauKFtTB798yBxk36e",
			wantErr: false,
		},
		{
			name: "t-P2SH",
			args: args{
				addr:    "0x0593920d01bbbeb3c649118c498334849df4bc6846",
				network: &chaincfg.MainNetParams,
			},
			want:    "3F9JDXk2JLxVFi95tZMUPcowoPjbQc7Cti",
			wantErr: false,
		},
		{
			name: "t-P2SH-without-0x",
			args: args{
				addr:    "0593920d01bbbeb3c649118c498334849df4bc6846",
				network: &chaincfg.MainNetParams,
			},
			want:    "3F9JDXk2JLxVFi95tZMUPcowoPjbQc7Cti",
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				addr:    "ff",
				network: &chaincfg.MainNetParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "unsupporte-public-key-prefix",
			args: args{
				addr:    "09593920d01bbbeb3c649118c498334849df4bc6846",
				network: &chaincfg.MainNetParams,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBitcoinAddress(tt.args.addr, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBitcoinAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.String() != tt.want {
				t.Errorf("DecodeBitcoinAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
