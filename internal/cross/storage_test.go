package cross_test

import (
	"fmt"
	"testing"

	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/cross"
)

func TestCrossStorage_Range(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		path string
		opts config.LevelDBOptions
		// Named input parameters for target function.
		key     string
		limit   int64
		want    []*cross.CrossMapping
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
			key:   "cross:00000000000000000020:order",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := cross.NewStorage(tt.path, tt.opts)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			for i := 0; i < 100; i++ {
				err := s.AddOrUpdateTx(&cross.CrossData{
					TxHash:  "src_tx_hash" + fmt.Sprint(i),
					Topic:   "test",
					Height:  100,
					OrderId: "order" + fmt.Sprint("%d", i),
				}, "src")
				if err != nil {
					t.Fatalf("could not write to storage: %v", err)
				}
			}
			got, gotErr := s.Range(tt.key, tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Range() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Range() succeeded unexpectedly")
			}

			for _, item := range got {
				t.Logf("got item: %+v", item.Key)
			}
		})
	}
}
