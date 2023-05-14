package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPragma(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	row := db.db.QueryRow("PRAGMA foreign_keys")
	assert.NotNil(t, row)
	var foreignKeys int
	if err := row.Scan(&foreignKeys); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, 1, foreignKeys)
}

func TestAddBlockWithPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	b, err := db.addBlock(testBlockCreate())
	assert.NoError(t, err)
	assertTestBlock(t, b)
	assert.Equal(t, 1, len(b.Pauses))
	assertTestPause(t, b.Pauses[0])
}

func TestAddBlockWithoutPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	b, err := db.addBlock(testBlockCreateWithoutPause())
	assert.NoError(t, err)
	assertTestBlock(t, b)
	assert.Equal(t, 0, len(b.Pauses))
}

func TestAddPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	_, err := db.addPause(testPauseCreate())
	assert.Error(t, err)

	db.addBlock(testBlockCreateWithoutPause())

	p, err := db.addPause(testPauseCreate())
	assert.NoError(t, err)
	assertTestPause(t, p)
}

func TestGetBlockByID(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	_, err := db.getBlockByID(bID)
	assert.Error(t, err)

	db.addBlock(testBlockCreate())

	b, err := db.getBlockByID(bID)
	assert.NoError(t, err)
	assertTestBlock(t, b)
	assert.Equal(t, 1, len(b.Pauses))
	assertTestPause(t, b.Pauses[0])
}

func TestGetPauseByID(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	_, err := db.getPauseByID(pID)
	assert.Error(t, err)

	db.addBlock(testBlockCreate())

	p, err := db.getPauseByID(pID)
	assert.NoError(t, err)
	assertTestPause(t, p)
}

func TestGetAllBlocks(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	blocks, err := db.getAllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(blocks))

	db.addBlock(testBlockCreate())

	blocks, err = db.getAllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(blocks))
	assertTestBlock(t, blocks[0])
	assert.Equal(t, 1, len(blocks[0].Pauses))
	assertTestPause(t, blocks[0].Pauses[0])
}

func TestDeleteBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	db.addBlock(testBlockCreate())

	err := db.deleteBlock(bID)
	assert.NoError(t, err)
	_, err = db.getPauseByID(pID)
	assert.Error(t, err)
	_, err = db.getBlockByID(bID)
	assert.Error(t, err)
}

func TestDeletePause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	db.addBlock(testBlockCreate())

	err := db.deletePause(pID)
	assert.NoError(t, err)
	_, err = db.getPauseByID(pID)
	assert.Error(t, err)
}

func TestUpdateBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	db.addBlock(testBlockCreate())

	err := db.updateBlock(testBlockUpdated())
	assert.NoError(t, err)

	b, err := db.getBlockByID(bID)
	assert.NoError(t, err)
	assertTestBlockUpdated(t, b)
}

func TestUpdatePause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	db.addBlock(testBlockCreate())

	err := db.updatePause(testPauseUpdated())
	assert.NoError(t, err)

	p, err := db.getPauseByID(pID)
	assertTestPauseUpdated(t, p)
}
