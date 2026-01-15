package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mapprotocol/compass-tss/cmd/compass/docs"

	"github.com/mapprotocol/compass-tss/internal/cross"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// -------------------------------------------------------------------------------------
// cross Server
// -------------------------------------------------------------------------------------

// CrossServer to back cross-data
type CrossServer struct {
	logger    zerolog.Logger
	s         *http.Server
	dbStorage *cross.CrossStorage
}

// NewCrossServer create a new instance of health server
func NewCrossServer(addr string, dbStorage *cross.CrossStorage) *CrossServer {
	hs := &CrossServer{
		logger:    log.With().Str("module", "cross").Logger(),
		dbStorage: dbStorage,
	}
	s := &http.Server{
		Addr:              addr,
		Handler:           hs.newHandler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	hs.s = s

	return hs
}

func (s *CrossServer) newHandler() http.Handler {
	router := mux.NewRouter()
	router.Handle("/ping", http.HandlerFunc(s.pingHandler)).Methods(http.MethodGet)
	router.Handle("/cross/chain/height", http.HandlerFunc(s.chainHeight)).Methods(http.MethodGet)
	router.Handle("/cross/chain/height/txs", http.HandlerFunc(s.chainTxs)).Methods(http.MethodGet)
	router.Handle("/cross/order", http.HandlerFunc(s.crossSignel)).Methods(http.MethodGet)
	router.Handle("/cross/pending/tx", http.HandlerFunc(s.pendingTx)).Methods(http.MethodGet)
	router.Handle("/cross/tx", http.HandlerFunc(s.crossFindByTx)).Methods(http.MethodGet)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

func (s *CrossServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CrossSignelResponse is the response for cross signel
type CrossSignelResponse struct {
	Data *cross.CrossSet `json:"data"`
}

// ChainHeightResponse
type ChainHeightResponse struct {
	Height string `json:"height"`
}

// ChainHeightResponse
type ChainOrderIdResponse struct {
	Height string   `json:"height"`
	Set    []string `json:"set"`
}

// PendingTxResponse
type PendingTxResponse struct {
	Txs []string `json:"txs"`
}

// get tx record by orderId
// @Summary      通过orderId获取交易记录
// @Description  通过orderId获取交易记录
// @Tags         交易记录
// @Accept       json
// @Produce      json
// @Param        orderId query string true "orderId"
// @Success      200  {object}  CrossSignelResponse
// @Failure      400  {object}  nil  "bad request"
// @Router       /cross/order [get]
func (s *CrossServer) crossSignel(w http.ResponseWriter, request *http.Request) {
	orderId := request.URL.Query().Get("orderId")
	crossData, err := s.dbStorage.GetCrossData(orderId)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	crossData.StatusStr = crossData.Status.String()
	res := &CrossSignelResponse{
		Data: crossData,
	}

	// write the response
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to write to response")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to write to response")
		}
	}
}

// get tx record by txHash
// @Summary      通过 txHash 获取交易记录
// @Description  通过 txHash 获取交易记录
// @Tags         交易记录
// @Accept       json
// @Produce      json
// @Param        tx query string true "txHash"
// @Success      200  {object}  CrossSignelResponse
// @Failure      400  {object}  nil  "bad request"
// @Router       /cross/tx [get]
func (s *CrossServer) crossFindByTx(w http.ResponseWriter, request *http.Request) {
	key := request.URL.Query().Get("tx")
	s.logger.Info().Any("tx", key).Msg("get cross signel by tx")
	crossData, err := s.dbStorage.GetCrossDataByTx(key)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	crossData.StatusStr = crossData.Status.String()
	res := &CrossSignelResponse{
		Data: crossData,
	}

	// write the response
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to write to response")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to write to response")
		}
	}
}

// get scanner chain txs
// @Summary      获取高度对应的交易集群
// @Description  根据 chainId 获取扫描交易集合
// @Tags         交易记录
// @Accept       json
// @Produce      json
// @Param        chainId query string true "1"
// @Param        height query string true "12245"
// @Success      200  {object}   ChainOrderIdResponse
// @Failure      400  {object}  nil  "bad request"
// @Router       /cross/chain/height/txs [get]
func (s *CrossServer) chainTxs(w http.ResponseWriter, request *http.Request) {
	chainId := request.URL.Query().Get("chainId")
	height := request.URL.Query().Get("height")
	s.logger.Info().Any("chainId", chainId).Any("height", height).Msg("get chain height")
	heightI, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderIdSet, err := s.dbStorage.GetOrderIdSet(chainId, heightI)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	scannerHeight, err := s.dbStorage.GetChainHeight(chainId)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &ChainOrderIdResponse{
		Height: scannerHeight,
		Set:    orderIdSet,
	}

	// write the response
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to write to response")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to write to response")
		}
	}
}

// get scanner chain height
// @Summary      获取扫描高度
// @Description  根据 chainId 获取扫描高度
// @Tags         交易记录
// @Accept       json
// @Produce      json
// @Param        chainId query string true "chainId"
// @Success      200  {object}   ChainHeightResponse
// @Failure      400  {object}  nil  "bad request"
// @Router       /cross/chain/height [get]
func (s *CrossServer) chainHeight(w http.ResponseWriter, request *http.Request) {
	chainId := request.URL.Query().Get("chainId")
	s.logger.Info().Any("chainId", chainId).Msg("get chain height")
	height, err := s.dbStorage.GetChainHeight(chainId)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &ChainHeightResponse{
		Height: height,
	}

	// write the response
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to write to response")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to write to response")
		}
	}
}

// get pending txs by chainId
// @Summary      获取pending的交易列表
// @Description  根据 chainId 获取pending的交易列表
// @Tags         交易记录
// @Accept       json
// @Produce      json
// @Param        chainId query string true "chainId"
// @Success      200  {object}  PendingTxResponse
// @Failure      400  {object}  nil  "bad request"
// @Router       /cross/pending/tx [get]
func (s *CrossServer) pendingTx(w http.ResponseWriter, request *http.Request) {
	chainId := request.URL.Query().Get("chainId")
	s.logger.Info().Any("chainId", chainId).Msg("get chain pending txs")
	txs, err := s.dbStorage.GetPendingSet(chainId)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &PendingTxResponse{
		Txs: txs,
	}

	// write the response
	jsonBytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to write to response")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = w.Write(jsonBytes)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to write to response")
		}
	}
}

// Start health server
func (s *CrossServer) Start() error {
	if s.s == nil {
		return errors.New("invalid http server instance")
	}
	if err := s.s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("fail to start http server: %w", err)
		}
	}
	return nil
}

func (s *CrossServer) Stop() error {
	s.logger.Info().Msg("shutting down cross server...")
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.s.Shutdown(c)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to shutdown the Cross server gracefully")
	}
	return err
}
