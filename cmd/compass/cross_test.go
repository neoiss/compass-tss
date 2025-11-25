package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/mapprotocol/compass-tss/config"
	"github.com/mapprotocol/compass-tss/internal/cross"
	"github.com/stretchr/testify/assert"
)

func TestCrossServer(t *testing.T) {
	crossStorage, err := cross.NewStorage("./test", config.LevelDBOptions{})
	assert.Nil(t, err)
	s := NewCrossServer("127.0.0.1:8080", crossStorage)
	assert.NotNil(t, s)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.Start()
		assert.Nil(t, err)
	}()
	time.Sleep(time.Second)
	assert.NotNil(t, err)
}

func TestPingHandler(t *testing.T) {
	crossStorage, err := cross.NewStorage("./test", config.LevelDBOptions{})
	assert.Nil(t, err)
	s := NewCrossServer("127.0.0.1:8080", crossStorage)
	assert.NotNil(t, s)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()
	s.pingHandler(res, req)
	assert.Equal(t, res.Code, http.StatusOK)
}

func TestGetP2pIDHandler(t *testing.T) {
	crossStorage, err := cross.NewStorage("./test", config.LevelDBOptions{})
	assert.Nil(t, err)
	s := NewCrossServer("127.0.0.1:8080", crossStorage)
	assert.NotNil(t, s)

	req := httptest.NewRequest(http.MethodGet, "/cross/list", nil)
	res := httptest.NewRecorder()
	s.crossList(res, req)
	assert.Equal(t, res.Code, http.StatusOK)
}
