package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kilianmandscharo/work_hours/models"
	"github.com/stretchr/testify/assert"
)

const (
	BID                = 1
	BStart             = "2023-05-09T07:00:00Z"
	BEnd               = "2023-05-09T15:30:00Z"
	BStartUpdated      = "2023-05-09T08:30:00Z"
	BEndUpdated        = "2023-05-09T17:00:00Z"
	BHomeoffice        = false
	BHomeofficeUpdated = true
	PID                = 1
	PStart             = "2023-05-09T12:00:00Z"
	PEnd               = "2023-05-09T12:30:00Z"
	PStartUpdated      = "2023-05-09T13:00:00Z"
	PEndUpdated        = "2023-05-09T13:30:00Z"
	PBlockID           = 1
)

func AssertTestBlock(t *testing.T, b models.Block) {
	assert.Equal(t, BID, b.Id)
	assert.Equal(t, BStart, b.Start)
	assert.Equal(t, BEnd, b.End)
	assert.Equal(t, BHomeoffice, b.Homeoffice)
}

func AssertTestBlockUpdated(t *testing.T, b models.Block) {
	assert.Equal(t, BID, b.Id)
	assert.Equal(t, BStartUpdated, b.Start)
	assert.Equal(t, BEndUpdated, b.End)
	assert.Equal(t, BHomeofficeUpdated, b.Homeoffice)
}

func AssertTestPause(t *testing.T, p models.Pause) {
	assert.Equal(t, PID, p.Id)
	assert.Equal(t, PStart, p.Start)
	assert.Equal(t, PEnd, p.End)
	assert.Equal(t, PBlockID, p.BlockID)
}

func AssertTestPauseUpdated(t *testing.T, p models.Pause) {
	assert.Equal(t, PID, p.Id)
	assert.Equal(t, PStartUpdated, p.Start)
	assert.Equal(t, PEndUpdated, p.End)
	assert.Equal(t, PBlockID, p.BlockID)
}

func TestBlockCreate() models.BlockCreate {
	pauses := make([]models.PauseWithoutBlockID, 1)
	pauses[0].Start = PStart
	pauses[0].End = PEnd
	block := TestBlockCreateWithoutPause()
	block.Pauses = pauses
	return block
}

func TestBlockCreateWithoutPause() models.BlockCreate {
	return models.BlockCreate{
		Start:      BStart,
		End:        BEnd,
		Homeoffice: BHomeoffice,
	}
}

func TestBlock() models.Block {
	pauses := make([]models.Pause, 1)
	pauses[0] = TestPause()
	return models.Block{
		Id:         BID,
		Start:      BStart,
		End:        BEnd,
		Homeoffice: BHomeoffice,
		Pauses:     pauses,
	}
}

func TestBlockUpdated() models.Block {
	return models.Block{
		Id:         BID,
		Start:      BStartUpdated,
		End:        BEndUpdated,
		Homeoffice: BHomeofficeUpdated,
	}
}

func TestPauseCreate() models.PauseCreate {
	return models.PauseCreate{
		Start:   PStart,
		End:     PEnd,
		BlockID: PBlockID,
	}
}

func TestPause() models.Pause {
	return models.Pause{
		Id:      PID,
		Start:   PStart,
		End:     PEnd,
		BlockID: PBlockID,
	}
}

func TestPauseUpdated() models.Pause {
	return models.Pause{
		Id:      PID,
		Start:   PStartUpdated,
		End:     PEndUpdated,
		BlockID: PBlockID,
	}
}

func AssertRequest(t *testing.T, r *gin.Engine, token string, method string, route string, statusWant int) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, route, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)
	assert.Equal(t, statusWant, w.Code)
}

func AssertRequestWithBody(t *testing.T, r *gin.Engine, token string, method string, route string, data any, statusWant int) {
	reader := getTestReader(t, data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, route, reader)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)
	assert.Equal(t, statusWant, w.Code)
}

func getTestReader(t *testing.T, data any) *bytes.Reader {
	dataBytes, err := json.Marshal(data)
	assert.NoError(t, err)
	dataReader := bytes.NewReader(dataBytes)
	return dataReader
}

func CreateRangeTestBlocks() []models.BlockCreate {
	return []models.BlockCreate{
		{
			Start: "2023-05-09T07:00:00Z",
			End:   "2023-05-09T15:30:00Z",
		},
		{
			Start: "2023-06-09T07:00:00Z",
			End:   "2023-06-09T15:30:00Z",
		},
		{
			Start: "2023-07-09T07:00:00Z",
			End:   "2023-07-09T15:30:00Z",
		},
	}

}
