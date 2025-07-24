package tss

import (
	"fmt"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"time"
)

const keySharesHex = ""

func Test_recoverKeyShares(t *testing.T) {
	type args struct {
		path       string
		keyShares  []byte
		passphrase string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path:       fmt.Sprintf("./localstate-%d.json", time.Now().UnixMilli()),
				keyShares:  common.Hex2Bytes(keySharesHex),
				passphrase: Mnemonic,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := recoverKeyShares(tt.args.path, tt.args.keyShares, tt.args.passphrase); (err != nil) != tt.wantErr {
				t.Errorf("recoverKeyShares() error  = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMnemonic(t *testing.T) {

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		t.Fatal(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mnemonic)
}
