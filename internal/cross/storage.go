package cross

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/db"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/syndtr/goleveldb/leveldb"
)

// CrossStorage save the ondeck tx in item to key value store, in case bifrost restart
type CrossStorage struct {
	db *leveldb.DB
}

const (
	CrossChainPrefix = "cross"
)

type StatusOfCross int64

const (
	StatusOfInit StatusOfCross = iota
	StatusOfPending
	StatusOfSend
	StatusOfCompleted
	StatusOfFailed
)

const (
	TypeOfSrcChain   = "src"
	TypeOfRelayChain = "relay"
	TypeOfSendDst    = "send_dst"
	TypeOfDstChain   = "dst"
)

type CrossData struct {
	TxHash           string
	Topic            string
	Height           int64
	OrderId          string
	LogIndex         uint
	Chain            string
	ChainAndGasLimit string
	Timestamp        int64 // tx time
}

type CrossSet struct {
	Src    *CrossData    `json:"src"`
	Relay  *CrossData    `json:"relay"`
	Dest   *CrossData    `json:"dest"`
	Now    int64         `json:"now"`
	Status StatusOfCross `json:"status"`
}

type CrossMapping struct {
	Key      string    `json:"key"`
	CrossSet *CrossSet `json:"cross_set"`
}

// NewStorage create a new instance of LevelDBScannerStorage
func NewStorage(path string, opts config.LevelDBOptions) (*CrossStorage, error) {
	ldb, err := db.NewLevelDB(path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create observer storage: %w", err)
	}

	return &CrossStorage{db: ldb}, nil
}

// createTxKey creates a unique key for a TxIn based on prefix, chain, mempool, blockheight
func (s *CrossStorage) createTxKey(txIn *types.TxInItem) string {
	return fmt.Sprintf("%s:%s", CrossChainPrefix, txIn.OrderId.String())
}

// AddOrUpdateTx adds or updates a single TxIn in storage
func (s *CrossStorage) AddOrUpdateTx(txIn *types.TxInItem, _type string) error {
	key := s.createTxKey(txIn)
	height := int64(0)
	if txIn.Height != nil {
		height = txIn.Height.Int64()
	}
	fromChain := ""
	if txIn.FromChain != nil {
		fromChain = txIn.FromChain.String()
	}
	cgl := ""
	if txIn.ChainAndGasLimit != nil {
		cgl = txIn.ChainAndGasLimit.String()
	}
	insertData := &CrossData{
		TxHash:           txIn.Tx,
		Topic:            txIn.Topic,
		Height:           height,
		OrderId:          txIn.OrderId.String(),
		LogIndex:         txIn.LogIndex,
		Chain:            fromChain,
		ChainAndGasLimit: cgl,
		Timestamp:        txIn.Timestamp,
	}

	ret, err := s.GetCrossData(key)
	if err != nil {
		return fmt.Errorf("fail to get crossData: %w", err)
	}
	switch _type {
	case TypeOfSrcChain:
		ret.Src = insertData
		ret.Status = StatusOfInit
	case TypeOfRelayChain:
		ret.Relay = insertData
		ret.Status = StatusOfPending
	case TypeOfSendDst:
		ret.Dest = insertData
		ret.Status = StatusOfSend
	case TypeOfDstChain:
		ret.Dest = insertData
		ret.Status = StatusOfCompleted
	default:
		return fmt.Errorf("invalid type:%s", _type)
	}
	data, err := json.Marshal(ret)
	if err != nil {
		return fmt.Errorf("fail to marshal tx to json: %w", err)
	}
	return s.db.Put([]byte(key), data, nil)
}

func (s *CrossStorage) GetCrossData(key string) (*CrossSet, error) {
	retBytes, err := s.db.Get([]byte(key), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}
	ret := &CrossSet{}
	if len(retBytes) == 0 {
		return ret, nil
	}
	err = json.Unmarshal(retBytes, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *CrossStorage) Range(key string, limit int64) ([]*CrossMapping, error) {
	ret := make([]*CrossMapping, 0, limit)
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()
	if key != "" {
		ok := iter.Seek([]byte(key))
		if !ok {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		ele := &CrossSet{}
		err := json.Unmarshal(value, ele)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &CrossMapping{
			Key:      string(key),
			CrossSet: ele,
		})
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	return ret, nil
}

func (s *CrossStorage) DeleteTx(txIn *types.TxInItem) error {
	key := s.createTxKey(txIn)
	return s.db.Delete([]byte(key), nil)
}

func (s *CrossStorage) Close() error {
	return s.db.Close()
}
