package utxo

import (
	"github.com/mapprotocol/compass-tss/common"
	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/keys"
	"github.com/mapprotocol/compass-tss/metrics"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/mapo"
	"github.com/mapprotocol/compass-tss/pkg/chainclients/shared/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func GetMetricForTest() (*metrics.Metrics, error) {
	return metrics.NewMetrics(config.BifrostMetricsConfiguration{
		Enabled:      false,
		ListenPort:   9000,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Chains:       common.Chains{common.ETHChain},
	})
}

func getBridgeForTest(t *testing.T) types.Bridge {
	m, err := GetMetricForTest()
	assert.Nil(t, err)

	bridgeCfg := config.BifrostClientConfiguration{
		ChainID:         "map",
		ChainHost:       "https://testnet-rpc.maplabs.io",
		SignerPasswd:    "password",
		ChainHomeFolder: "./",
		Maintainer:      "0x0EdA5e4015448A2283662174DD7def3C3d262D38",
		ViewController:  "0x7Ea4dFBa2fA7de4C18395aCD391D9E67bECA47A6",
	}

	name := "test-eth"
	//  dont push
	keyStorePath := "/Users/t/Documents/compass-tss/keys/0x25fa71d4f689f4b65eb6d020a414090828281d51"
	kb, keyStore, err := keys.GetKeyringKeybase(keyStorePath, name)
	assert.Nil(t, err)

	k := keys.NewKeysWithKeybase(kb, name, "123456", keyStore)
	bridge, err := mapo.NewBridge(bridgeCfg, m, k)
	assert.Nil(t, err)
	return bridge
}

//func TestGetAsgardAddress2PubKeyMapped(t *testing.T) {
//	type args struct {
//		chain string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    map[common.Address][]byte
//		wantErr bool
//	}{
//		{
//			name: "t-1",
//			args: args{
//				chain: "212",
//			},
//			want:    nil,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bridge := getBridgeForTest(t)
//			chain, err := common.NewChain(tt.args.chain)
//			if err != nil {
//				t.Errorf("NewChain() error = %v", err)
//				return
//			}
//			got, err := GetAsgardAddress2PubKeyMapped(chain, bridge)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetAsgardAddress2PubKeyMapped() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetAsgardAddress2PubKeyMapped() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
