package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)

	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		block := testBlockCreate()
		blockBytes, err := json.Marshal(block)
		assert.NoError(t, err)
		blockReader := bytes.NewReader(blockBytes)
		req, _ := http.NewRequest(http.MethodPost, "/block", blockReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var newBlock Block
		err = json.Unmarshal([]byte(w.Body.String()), &newBlock)
		assert.NoError(t, err)
		assert.EqualValues(t, testBlock(), newBlock)
	})
}

func TestUpdateBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/block", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		block := testBlockUpdated()
		blockBytes, err := json.Marshal(block)
		assert.NoError(t, err)
		blockReader := bytes.NewReader(blockBytes)
		req, _ := http.NewRequest(http.MethodPut, "/block", blockReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeleteBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/block/a", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/block/%d", bID), nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetBlockByIDRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	_, err := db.addBlock(testBlockCreate())
	assert.NoError(t, err)
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block/a", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no block", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/block/%d", 2), nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("block found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/block/%d", bID), nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var block Block
		err = json.Unmarshal([]byte(w.Body.String()), &block)
		assert.NoError(t, err)
		assert.EqualValues(t, testBlock(), block)
	})
}

func TestGetAllBlocksRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	_, err := db.addBlock(testBlockCreate())
	assert.NoError(t, err)
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("blocks found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var blocks []Block
		err = json.Unmarshal([]byte(w.Body.String()), &blocks)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(blocks))
		assert.EqualValues(t, testBlock(), blocks[0])
	})
}

func TestAddPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	_, err := db.addBlock(testBlockCreateWithoutPause())
	assert.NoError(t, err)
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("wrong block ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		pause := testPauseCreate()
		pause.BlockID = 2
		pauseBytes, err := json.Marshal(pause)
		assert.NoError(t, err)
		pauseReader := bytes.NewReader(pauseBytes)
		req, _ := http.NewRequest(http.MethodPost, "/pause", pauseReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		pause := testPauseCreate()
		pauseBytes, err := json.Marshal(pause)
		assert.NoError(t, err)
		pauseReader := bytes.NewReader(pauseBytes)
		req, _ := http.NewRequest(http.MethodPost, "/pause", pauseReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var newPause Pause
		err = json.Unmarshal([]byte(w.Body.String()), &newPause)
		assert.NoError(t, err)
		assert.EqualValues(t, testPause(), newPause)
	})
}

func TestUpdatePauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	_, err := db.addBlock(testBlockCreateWithoutPause())
	assert.NoError(t, err)
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		pause := testPauseUpdated()
		pauseBytes, err := json.Marshal(pause)
		assert.NoError(t, err)
		pauseReader := bytes.NewReader(pauseBytes)
		req, _ := http.NewRequest(http.MethodPost, "/pause", pauseReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeletePauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	_, err := db.addBlock(testBlockCreate())
	assert.NoError(t, err)
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/pause/a", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/pause/%d", pID), nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
