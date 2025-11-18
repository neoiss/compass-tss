package utxo

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const broadcastCachePrefix = "broadcast-v1-"

const (
	TxStatusSent uint8 = iota + 1
	TxStatusConfirmed
)

type TransactionStorage struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

type TxStatus struct {
	Status    uint8
	TxHash    string
	Timestamp int64
}

func NewTransactionSStorage(db *leveldb.DB) (*TransactionStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("db parameter is nil")
	}
	return &TransactionStorage{
		db:     db,
		logger: log.With().Str("module", "broadcast-cache").Logger(),
	}, nil
}

func (bs *TransactionStorage) SetTxStatus(hash string, status uint8) error {
	key := []byte(hash)
	txStatus := TxStatus{
		Status:    status,
		TxHash:    hash,
		Timestamp: time.Now().Unix(),
	}
	value, err := json.Marshal(txStatus)
	if err != nil {
		bs.logger.Error().Err(err).Interface("txStatus", txStatus).Msg("fail to marshal tx status")
		return err
	}
	if err := bs.db.Put(key, value, nil); err != nil {
		bs.logger.Error().Err(err).Msg("fail to set tx status")
		return err
	}
	return nil
}

func (bs *TransactionStorage) GetTxStatus(hash string) (status uint8, err error) {
	key := []byte(hash)

	ok, err := bs.db.Has(key, nil)
	if !ok || err != nil {
		return
	}
	buf, err := bs.db.Get(key, nil)
	if err != nil {
		bs.logger.Error().Err(err).Msg("fail to get tx status")
		return
	}
	txStatus := TxStatus{}
	if err = json.Unmarshal(buf, &txStatus); err != nil {
		bs.logger.Error().Err(err).Msg("fail to unmarshal tx status")
		return
	}
	return txStatus.Status, nil
}

func (bs *TransactionStorage) List() []TxStatus {
	iterator := bs.db.NewIterator(util.BytesPrefix([]byte(broadcastCachePrefix)), nil)
	defer iterator.Release()
	var results []TxStatus
	for iterator.Next() {
		buf := iterator.Value()
		if len(buf) == 0 {
			continue
		}

		var item TxStatus
		if err := json.Unmarshal(buf, &item); err != nil {
			bs.logger.Error().Err(err).Msg("fail to unmarshal to txout store item")
			continue
		}

		// ignore already confirmed tx
		if item.Status == TxStatusConfirmed {
			continue
		}

		results = append(results, item)
	}

	sort.SliceStable(results, func(i, j int) bool { return results[i].Timestamp < results[j].Timestamp })
	return results
}

func (bs *TransactionStorage) ListConfirmed() []TxStatus {
	iterator := bs.db.NewIterator(util.BytesPrefix([]byte(broadcastCachePrefix)), nil)
	defer iterator.Release()
	var results []TxStatus
	for iterator.Next() {
		buf := iterator.Value()
		if len(buf) == 0 {
			continue
		}

		var item TxStatus
		if err := json.Unmarshal(buf, &item); err != nil {
			bs.logger.Error().Err(err).Msg("fail to unmarshal to txout store item")
			continue
		}

		// ignore already confirmed tx
		if item.Status != TxStatusConfirmed {
			continue
		}

		results = append(results, item)
	}

	sort.SliceStable(results, func(i, j int) bool { return results[i].Timestamp < results[j].Timestamp })
	return results
}

func (bs *TransactionStorage) Remove(hash string) error {
	return bs.db.Delete([]byte(hash), nil)
}

func (bs *TransactionStorage) getBroadcastKey(hash string) string {
	return fmt.Sprintf("%s%s", broadcastCachePrefix, hash)
}

// Close underlying db
func (bs *TransactionStorage) Close() error {
	return bs.db.Close()
}
