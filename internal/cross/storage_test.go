package cross_test

import (
	"encoding/json"
	"testing"

	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/cross"
)

func TestCrossStorage_HandlerCrossData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		path string
		opts config.LevelDBOptions
		// Named input parameters for target function.
		key     cross.ChanStruct
		limit   int64
		wantErr bool
	}{
		{
			name: "success",
			path: "./testdata/cross_storage_range_success",
			opts: config.LevelDBOptions{
				BlockCacheCapacity:            1 << 20,
				CompactOnInit:                 true,
				CompactionTableSizeMultiplier: 1.0,
				FilterBitsPerKey:              10,
				WriteBuffer:                   1 << 20,
			},
			limit: 10,
			key: cross.ChanStruct{
				CrossData: &cross.CrossData{
					OrderId:          "0x015be4c33f51fbee02e13b93f5bc2089e2cde770810a5510225de4ce7a8375a1",
					Chain:            "1360095883558914",
					Height:           287013,
					LogIndex:         1,
					TxHash:           "3183f09f1a960c1d401dfa63f07ecab06bf586f3db84e88a3e1237c799ecd55e",
					Topic:            "0x8104943fdd0997a3240b59b381251572ac6ac81941e1af29845de70edca938a4",
					Timestamp:        1768536636,
					ChainAndGasLimit: "462816646496852369478131125222857703149698333785718784",
				},
				Type: "src",

				// CrossData: &cross.CrossData{
				// 	OrderId:          "0x015be4c33f51fbee02e13b93f5bc2089e2cde770810a5510225de4ce7a8375a1",
				// 	Chain:            "212",
				// 	Height:           287013,
				// 	LogIndex:         1,
				// 	TxHash:           "3183f09f1a960c1d401dfa63f07ecab06bf586f3db84e88a3e1237c799ecd55e",
				// 	Topic:            "0x8104943fdd0997a3240b59b381251572ac6ac81941e1af29845de70edca938a4",
				// 	Timestamp:        1768536636,
				// 	ChainAndGasLimit: "462816646496852369478131125222857703149698333785718784",
				// },
				// Type: "relay",
				// CrossData: &cross.CrossData{
				// 	OrderId:          "0x780a6ea54f538434331c3c1490a43c984e7e89d4434f53822437e97a98b0d146",
				// 	Chain:            "212",
				// 	Height:           20160120,
				// 	LogIndex:         3,
				// 	TxHash:           "0x6dfbefc5bd6a6d0399bd6d21c48015a35bf49db361a5b466da5286d4a32a5c01",
				// 	Topic:            "0x298a40641bd31f72c733761e0e85a6bd8a36909666ac2ed63a42c8015d025638",
				// 	Timestamp:        1768534943,
				// 	ChainAndGasLimit: "462816646496852369478131125222857703149698333785718784",
				// },
				// Type: "map_dst",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := cross.NewStorage(tt.path, tt.opts)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := s.HandlerCrossData(&tt.key)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Range() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Range() succeeded unexpectedly")
			}

		})
	}
}

func TestCrossStorage_GetCrossData(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		path string
		opts config.LevelDBOptions
		// Named input parameters for target function.
		orderId string
		want    *cross.CrossSet
		wantErr bool
	}{
		{
			name: "success",
			path: "./testdata/cross_storage_range_success",
			opts: config.LevelDBOptions{
				BlockCacheCapacity:            1 << 20,
				CompactOnInit:                 true,
				CompactionTableSizeMultiplier: 1.0,
				FilterBitsPerKey:              10,
				WriteBuffer:                   1 << 20,
			},
			orderId: "0x015be4c33f51fbee02e13b93f5bc2089e2cde770810a5510225de4ce7a8375a1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := cross.NewStorage(tt.path, tt.opts)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := s.GetCrossData(tt.orderId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCrossData() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCrossData() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				data, _ := json.Marshal(got)
				t.Errorf("GetCrossData() = %v, want %v", string(data), tt.want)
			}
		})
	}
}

func TestCrossStorage_GetChainHeight(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		path string
		opts config.LevelDBOptions
		// Named input parameters for target function.
		chainId string
		want    string
		wantErr bool
	}{
		{
			name: "success",
			path: "./testdata/cross_storage_range_success",
			opts: config.LevelDBOptions{
				BlockCacheCapacity:            1 << 20,
				CompactOnInit:                 true,
				CompactionTableSizeMultiplier: 1.0,
				FilterBitsPerKey:              10,
				WriteBuffer:                   1 << 20,
			},
			chainId: "212",
			want:    "20160120",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := cross.NewStorage(tt.path, tt.opts)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := s.GetChainHeight(tt.chainId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetChainHeight() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetChainHeight() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetChainHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}
