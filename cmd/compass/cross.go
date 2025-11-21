package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mapprotocol/compass-tss/internal/cross"

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
	router.Handle("/cross/list", http.HandlerFunc(s.crossList)).Methods(http.MethodGet)
	router.Handle("/cross/signle", http.HandlerFunc(s.crossSignel)).Methods(http.MethodGet)
	return router
}

func (s *CrossServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type CrossSignelRequest struct {
	Key string `json:"key"`
}

type CrossSignelResponse struct {
	Data *cross.CrossSet `json:"data"`
}

type CrossListRequest struct {
	Key   string `json:"key"`
	Limit int64  `json:"limit"`
}

type CrossListResponse struct {
	Data []*cross.CrossMapping `json:"data"`
}

func (s *CrossServer) crossList(w http.ResponseWriter, request *http.Request) {
	reqData, err := io.ReadAll(request.Body)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := &CrossListRequest{}
	err = json.Unmarshal(reqData, req)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	crossData, err := s.dbStorage.Range(req.Key, req.Limit)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &CrossListResponse{
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

func (s *CrossServer) crossSignel(w http.ResponseWriter, request *http.Request) {
	reqData, err := io.ReadAll(request.Body)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := &CrossSignelRequest{}
	err = json.Unmarshal(reqData, req)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	crossData, err := s.dbStorage.GetCrossData(req.Key)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get cross data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
