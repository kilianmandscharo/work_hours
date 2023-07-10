package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bID                = 1
	bStart             = "2023-05-09T07:00:00Z"
	bEnd               = "2023-05-09T15:30:00Z"
	bStartUpdated      = "2023-05-09T08:30:00Z"
	bEndUpdated        = "2023-05-09T17:00:00Z"
	bHomeoffice        = false
	bHomeofficeUpdated = true
	pID                = 1
	pStart             = "2023-05-09T12:00:00Z"
	pEnd               = "2023-05-09T12:30:00Z"
	pStartUpdated      = "2023-05-09T13:00:00Z"
	pEndUpdated        = "2023-05-09T13:30:00Z"
	pBlockID           = 1
)

func assertTestBlock(t *testing.T, b Block) {
	assert.Equal(t, bID, b.Id)
	assert.Equal(t, bStart, b.Start)
	assert.Equal(t, bEnd, b.End)
	assert.Equal(t, bHomeoffice, b.Homeoffice)
}

func assertTestBlockUpdated(t *testing.T, b Block) {
	assert.Equal(t, bID, b.Id)
	assert.Equal(t, bStartUpdated, b.Start)
	assert.Equal(t, bEndUpdated, b.End)
	assert.Equal(t, bHomeofficeUpdated, b.Homeoffice)
}

func assertTestPause(t *testing.T, p Pause) {
	assert.Equal(t, pID, p.Id)
	assert.Equal(t, pStart, p.Start)
	assert.Equal(t, pEnd, p.End)
	assert.Equal(t, pBlockID, p.BlockID)
}

func assertTestPauseUpdated(t *testing.T, p Pause) {
	assert.Equal(t, pID, p.Id)
	assert.Equal(t, pStartUpdated, p.Start)
	assert.Equal(t, pEndUpdated, p.End)
	assert.Equal(t, pBlockID, p.BlockID)
}

func testBlockCreate() BlockCreate {
	pauses := make([]PauseWithoutBlockID, 1)
	pauses[0].Start = pStart
	pauses[0].End = pEnd
	block := testBlockCreateWithoutPause()
	block.Pauses = pauses
	return block
}

func testBlockCreateWithoutPause() BlockCreate {
	return BlockCreate{
		Start:      bStart,
		End:        bEnd,
		Homeoffice: bHomeoffice,
	}
}

func testBlock() Block {
	pauses := make([]Pause, 1)
	pauses[0] = testPause()
	return Block{
		Id:         bID,
		Start:      bStart,
		End:        bEnd,
		Homeoffice: bHomeoffice,
		Pauses:     pauses,
	}
}

func testBlockUpdated() Block {
	return Block{
		Id:         bID,
		Start:      bStartUpdated,
		End:        bEndUpdated,
		Homeoffice: bHomeofficeUpdated,
	}
}

func testPauseCreate() PauseCreate {
	return PauseCreate{
		Start:   pStart,
		End:     pEnd,
		BlockID: pBlockID,
	}
}

func testPause() Pause {
	return Pause{
		Id:      pID,
		Start:   pStart,
		End:     pEnd,
		BlockID: pBlockID,
	}
}

func testPauseUpdated() Pause {
	return Pause{
		Id:      pID,
		Start:   pStartUpdated,
		End:     pEndUpdated,
		BlockID: pBlockID,
	}
}

func getNewTestDatabase() *DB {
	db, err := newTestDatabase()
	if err != nil {
		log.Fatalf("ERROR: could not open test database, %v", err)
	}
	err = db.init()
	if err != nil {
		log.Fatalf("ERROR: could not initialize test database, %v", err)
	}
	return db
}
