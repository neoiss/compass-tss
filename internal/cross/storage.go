package cross

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/db"
	"github.com/mapprotocol/compass-tss/mapclient/types"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
)

// CrossStorage save the ondeck tx in item to key value store, in case bifrost restart
type CrossStorage struct {
	db   *leveldb.DB
	mu   sync.Mutex
	ch   chan *ChanStruct
	stop chan struct{}
}

const (
	CrossChainPrefix = "meta:order:%s"   // orderId
	KeyOfTxHash      = "meta:tx:%s"      // txHash
	KeyOfChainHeight = "meta:height:%s"  // chainId
	KeyOfOrderIdSet  = "meta:set:%s:%s"  // chainId:startHeight
	KeyOfPendingTx   = "meta:pending:%s" // chainId
)

type StatusOfCross int64

const (
	StatusOfInit StatusOfCross = iota
	StatusOfPending
	StatusOfSend
	StatusOfCompleted
	StatusOfFailed
)

func (s StatusOfCross) String() string {
	switch s {
	case StatusOfInit:
		return "init"
	case StatusOfPending:
		return "pending"
	case StatusOfSend:
		return "send"
	case StatusOfCompleted:
		return "completed"
	case StatusOfFailed:
		return "failed"
	}
	return ""
}

const (
	TypeOfSrcChain         = "src"
	TypeOfRelayChain       = "relay"
	TypeOfRelaySignedChain = "relay_signed"
	TypeOfSendDst          = "send_dst"
	TypeOfDstChain         = "dst"
	TypeOfMapDstChain      = "map_dst"
)

// CrossData
type CrossData struct {
	TxHash           string `json:"tx_hash" example:""  `
	Topic            string `json:"topic" example:""  `
	Height           int64  `json:"height" example:81507414 `
	OrderId          string `json:"order_id" example:"" `
	LogIndex         uint   `json:"log_index" example:1 `
	Chain            string `json:"chain" example:"" `
	ChainAndGasLimit string `json:"chain_and_gas_limit" example:"" `
	Timestamp        int64  `json:"timestamp" example: 1767097427 `
	IsMemoized       bool   `json:"is_memoized" example: false `
}

type ChanStruct struct {
	CrossData *CrossData
	Type      string
}

// CrossSet
type CrossSet struct {
	Src         *CrossData    `json:"src" `      // source chain transaction
	Relay       *CrossData    `json:"relay"  `   // relay chain transaction
	RelaySigned *CrossData    `json:"relay"  `   // relay signed transaction , The front end ignores this field.
	Dest        *CrossData    `json:"dest" `     // target Chain Transactions
	MapDst      *CrossData    `json:"map_dest" ` // map dest Transactions
	Now         int64         `json:"now" `
	Status      StatusOfCross `json:"status"`
	StatusStr   string        `json:"status_str"`
	OrderId     string        `json:"order_id"`
}

// NewStorage create a new instance of LevelDBScannerStorage
func NewStorage(path string, opts config.LevelDBOptions) (*CrossStorage, error) {
	ldb, err := db.NewLevelDB(path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create observer storage: %w", err)
	}

	return &CrossStorage{
		db:   ldb,
		mu:   sync.Mutex{},
		ch:   make(chan *ChanStruct, 100),
		stop: make(chan struct{}),
	}, nil
}

func (s *CrossStorage) Start() {
	log.Info().Msg("starting cross storage")
	go func() {
		for {
			select {
			case <-s.stop:
				return
			case ele, ok := <-s.ch:
				if !ok {
					log.Error().Msg("cross storage channel closed")
					return
				}
				err := s.HandlerCrossData(ele)
				if err != nil {
					log.Error().Any("ele", ele).Err(err).Msg("fail to handle cross data")
				}
			}
		}
	}()
}

func (s *CrossStorage) Stop() {
	log.Info().Msg("stop cross storage")
	close(s.stop)
}

// createOrderIDKey creates a unique key for a TxIn based on prefix, chain, mempool, blockheight
func (s *CrossStorage) createOrderIDKey(orderId string) string {
	return fmt.Sprintf(CrossChainPrefix, orderId)
}

func (s *CrossStorage) createTxKey(txHash string) string {
	return fmt.Sprintf(KeyOfTxHash, txHash)
}

func (s *CrossStorage) createChainHeightKey(chainId string) string {
	return fmt.Sprintf(KeyOfChainHeight, chainId)
}

func (s *CrossStorage) createOrderIdSetKey(chainId string, startHeight int64) string {
	cacheHeight := startHeight / 100 * 100
	return fmt.Sprintf(KeyOfOrderIdSet, chainId, strconv.FormatInt(cacheHeight, 10))
}

func (s *CrossStorage) createPendingKey(chainId string) string {
	return fmt.Sprintf(KeyOfPendingTx, chainId)
}

func TxInConvertCross(txIn *types.TxInItem) *CrossData {
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
	return &CrossData{
		TxHash:           txIn.Tx,
		Topic:            txIn.Topic,
		Height:           height,
		OrderId:          txIn.OrderId.String(),
		LogIndex:         txIn.LogIndex,
		Chain:            fromChain,
		ChainAndGasLimit: cgl,
		Timestamp:        time.Now().Unix(),
	}
}

func TxOutConvertCross(txOut *types.TxOutItem) *CrossData {
	return &CrossData{
		TxHash:           txOut.TxHash,
		Topic:            txOut.Topics,
		Height:           txOut.Height,
		OrderId:          txOut.OrderId.String(),
		LogIndex:         txOut.LogIndex,
		Chain:            txOut.Chain.String(),
		ChainAndGasLimit: txOut.ChainAndGasLimit.String(),
		Timestamp:        time.Now().Unix(),
	}
}

// AddOrUpdateTx adds or updates a single TxIn in storage
func (s *CrossStorage) AddOrUpdateTx(insertData *CrossData, _type string) {
	s.ch <- &ChanStruct{
		CrossData: insertData,
		Type:      _type,
	}
}

func (s *CrossStorage) HandlerCrossData(ele *ChanStruct) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pendingKey := s.createPendingKey(ele.CrossData.Chain)
	pendingTxs, err := s.GetPendingSet(ele.CrossData.Chain)
	if err != nil {
		return fmt.Errorf("GetPending failed err :%w", err)
	}
	if ele.CrossData.IsMemoized {
		pendingTxs = append(pendingTxs, ele.CrossData.TxHash)
		pendingTxsData, _ := json.Marshal(pendingTxs)
		return s.db.Put([]byte(pendingKey), pendingTxsData, nil)
	}

	key := s.createOrderIDKey(ele.CrossData.OrderId)
	ret, err := s.GetCrossData(ele.CrossData.OrderId)
	if err != nil {
		return fmt.Errorf("fail to get crossData: %w", err)
	}
	changeStatus := ret.Status
	switch ele.Type {
	case TypeOfSrcChain:
		ret.Src = ele.CrossData
		changeStatus = StatusOfInit
	case TypeOfRelayChain:
		if ret.Src == nil { // map sending tx
			ret.Src = ele.CrossData
		}
		ret.Relay = ele.CrossData
		changeStatus = StatusOfPending
	case TypeOfRelaySignedChain:
		// dont change status
		ret.RelaySigned = ele.CrossData
	case TypeOfSendDst:
		ret.Dest = ele.CrossData
		changeStatus = StatusOfSend
	case TypeOfDstChain:
		ret.Dest = ele.CrossData
		changeStatus = StatusOfCompleted
	case TypeOfMapDstChain:
		if ret.Relay == nil {
			ret.Relay = ele.CrossData
		}
		if ret.Dest == nil { // send to map
			ret.Dest = ele.CrossData
		}
		ret.MapDst = ele.CrossData
		changeStatus = StatusOfCompleted
	default:
		return fmt.Errorf("invalid type:%s", ele.Type)
	}
	if ret.Status < changeStatus {
		ret.Status = changeStatus
	}
	data, err := json.Marshal(ret)
	if err != nil {
		return fmt.Errorf("fail to marshal tx to json: %w", err)
	}

	txSetKey := s.createOrderIdSetKey(ele.CrossData.Chain, ele.CrossData.Height)
	orderIdSet, err := s.GetOrderIdSet(ele.CrossData.Chain, ele.CrossData.Height)
	orderIdSet = append(orderIdSet, ele.CrossData.OrderId)

	batch := new(leveldb.Batch)
	batch.Put([]byte(key), data)
	batch.Put([]byte(s.createTxKey(ele.CrossData.TxHash)), []byte(ele.CrossData.OrderId))
	batch.Put([]byte(s.createChainHeightKey(ele.CrossData.Chain)), []byte(strconv.Itoa(int(ele.CrossData.Height))))
	orderIdSetData, _ := json.Marshal(orderIdSet)
	batch.Put([]byte(txSetKey), orderIdSetData)
	if len(pendingTxs) > 0 {
		// rm this tx from pending
		newPendingTxs := make([]string, 0, len(pendingTxs)-1)
		for _, tx := range pendingTxs {
			if tx == ele.CrossData.TxHash {
				continue
			}
			newPendingTxs = append(newPendingTxs, tx)
		}

		pendingTxsData, _ := json.Marshal(newPendingTxs)
		batch.Put([]byte(pendingKey), pendingTxsData)
	}

	return s.db.Write(batch, nil)
}

func (s *CrossStorage) GetCrossData(orderId string) (*CrossSet, error) {
	key := s.createOrderIDKey(orderId)
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
	ret.StatusStr = ret.Status.String()
	return ret, nil
}

func (s *CrossStorage) GetCrossDataByTx(txHash string) (*CrossSet, error) {
	orderIdBytes, err := s.db.Get([]byte(s.createTxKey(txHash)), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}
	ret := &CrossSet{
		OrderId: string(orderIdBytes),
	}
	if len(orderIdBytes) == 0 {
		return ret, nil
	}
	retBytes, err := s.db.Get([]byte(s.createOrderIDKey(string(orderIdBytes))), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}
	if len(retBytes) == 0 {
		return ret, nil
	}
	err = json.Unmarshal(retBytes, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *CrossStorage) GetOrderIdSet(chainId string, height int64) ([]string, error) {
	key := s.createOrderIdSetKey(chainId, height)
	retBytes, err := s.db.Get([]byte(key), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}

	ret := make([]string, 0)
	if len(retBytes) == 0 {
		return ret, nil
	}

	err = json.Unmarshal(retBytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *CrossStorage) GetPendingSet(chainId string) ([]string, error) {
	key := s.createPendingKey(chainId)
	retBytes, err := s.db.Get([]byte(key), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, err
	}

	ret := make([]string, 0)
	if len(retBytes) == 0 {
		return ret, nil
	}

	err = json.Unmarshal(retBytes, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *CrossStorage) GetChainHeight(chainId string) (string, error) {
	key := s.createChainHeightKey(chainId)
	retBytes, err := s.db.Get([]byte(key), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return "", err
	}

	return string(retBytes), nil
}

func (s *CrossStorage) DeleteTx(key string) error {
	return s.db.Delete([]byte(key), nil)
}

func (s *CrossStorage) Close() error {
	return s.db.Close()
}
