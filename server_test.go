package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var token string

func init() {
	env, err := envVariables()
	if err != nil {
		log.Fatal("could not load env file")
	}
	token, err = createToken(env.email, env.tokenKey)
	if err != nil {
		log.Fatal("could not create token")
	}
}

func TestAddBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)

	gin.SetMode(gin.TestMode)

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/block/%d", bID), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no block", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/block/%d", 2), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("block found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/block/%d", bID), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no blocks available", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("blocks found", func(t *testing.T) {
		_, err := db.addBlock(testBlockCreate())
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/pause/%d", pID), nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStartBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block_start", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block_start", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestEndBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("block not started", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("pause still active", func(t *testing.T) {
		_, err := db.startBlock()
		assert.NoError(t, err)
		_, err = db.startPause()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		_, err := db.endPause()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/block_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStartPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid request no active block", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_start", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		_, err := db.startBlock()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_start", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid request pause already active", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_start", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestEndPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid request no active block", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid request no active pause", func(t *testing.T) {
		_, err := db.startBlock()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		_, err := db.startPause()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/pause_end", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetCurrentBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block_current", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		_, err := db.startBlock()
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/block_current", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLoginRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)
	envTest, err := envTestVariables()
	assert.NoError(t, err)

	t.Run("invalid body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid email", func(t *testing.T) {
		w := httptest.NewRecorder()
		login := Login{Email: "invalid@gmail.com", Password: envTest.password}
		loginBytes, err := json.Marshal(login)
		assert.NoError(t, err)
		loginReader := bytes.NewReader(loginBytes)
		req, _ := http.NewRequest(http.MethodPost, "/login", loginReader)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid password", func(t *testing.T) {
		w := httptest.NewRecorder()
		login := Login{Email: envTest.email, Password: "987654321"}
		loginBytes, err := json.Marshal(login)
		assert.NoError(t, err)
		loginReader := bytes.NewReader(loginBytes)
		req, _ := http.NewRequest(http.MethodPost, "/login", loginReader)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		login := Login{Email: envTest.email, Password: envTest.password}
		loginBytes, err := json.Marshal(login)
		assert.NoError(t, err)
		loginReader := bytes.NewReader(loginBytes)
		req, _ := http.NewRequest(http.MethodPost, "/login", loginReader)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRefreshRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("token still valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/refresh", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
