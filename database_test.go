package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialCurrentIDs(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, -1, currentBlockID)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, -1, currentPauseID)
}

func TestSetCurrentBlockID(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	testID := 10

	err := db.setCurrentBlockID(testID)
	assert.NoError(t, err)
	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, testID, currentBlockID)
}

func TestSetCurrentPauseID(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	testID := 10

	err := db.setCurrentPauseID(testID)
	assert.NoError(t, err)
	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, testID, currentPauseID)
}

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

	rowsAffected, err := db.deleteBlock(bID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.deleteBlock(bID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)
	_, err = db.getPauseByID(pID)
	assert.Error(t, err)
	_, err = db.getBlockByID(bID)
	assert.Error(t, err)
}

func TestDeleteCurrentBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	_, err := db.startBlock(false)
	assert.NoError(t, err)
	_, err = db.startPause()
	assert.NoError(t, err)

	rowsAffected, err := db.deleteBlock(bID)
	assert.Equal(t, rowsAffected, 1)

	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, currentBlockID, -1)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, currentPauseID, -1)
}

func TestDeletePause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.deletePause(pID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.deletePause(pID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)
	_, err = db.getPauseByID(pID)
	assert.Error(t, err)
}

func TestDeleteCurrentPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	_, err := db.startBlock(false)
	assert.NoError(t, err)
	_, err = db.startPause()
	assert.NoError(t, err)

	rowsAffected, err := db.deletePause(pID)
	assert.Equal(t, rowsAffected, 1)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, currentPauseID, -1)
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

func TestUpdateBlockStart(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.updateBlockStart(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.updateBlockStart(bID, bStartUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.getBlockByID(bID)
	assert.NoError(t, err)
	assert.Equal(t, block.Start, bStartUpdated)
}

func TestUpdateBlockEnd(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.updateBlockEnd(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.updateBlockEnd(bID, bEndUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.getBlockByID(bID)
	assert.NoError(t, err)
	assert.Equal(t, block.End, bEndUpdated)
}

func TestUpdateBlockHomeoffice(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.updateBlockHomeoffice(10, true)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.updateBlockHomeoffice(bID, bHomeofficeUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.getBlockByID(bID)
	assert.NoError(t, err)
	assert.Equal(t, block.Homeoffice, bHomeofficeUpdated)
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

func TestUpdatePauseStart(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.updatePauseStart(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.updatePauseStart(bID, pStartUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	pause, err := db.getPauseByID(bID)
	assert.NoError(t, err)
	assert.Equal(t, pause.Start, pStartUpdated)
}

func TestUpdatePauseEnd(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	rowsAffected, err := db.updatePauseEnd(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.addBlock(testBlockCreate())

	rowsAffected, err = db.updatePauseEnd(bID, pEndUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	pause, err := db.getPauseByID(bID)
	assert.NoError(t, err)
	assert.Equal(t, pause.End, pEndUpdated)
}

func TestStartBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	t.Run("start successful", func(t *testing.T) {
		_, err := db.startBlock(false)
		assert.NoError(t, err)
		currentBlockID, err := db.getCurrentBlockID()
		assert.NoError(t, err)
		assert.NotEqual(t, -1, currentBlockID)
	})

	t.Run("block already active", func(t *testing.T) {
		_, err := db.startBlock(false)
		assert.Error(t, err)
	})
}

func TestEndBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.endBlock()
		assert.Error(t, err)
	})

	t.Run("end successful", func(t *testing.T) {
		newBlock, err := db.startBlock(false)
		assert.NoError(t, err)
		block, err := db.endBlock()
		assert.NoError(t, err)
		assert.Equal(t, newBlock.Id, block.Id)
		assert.Equal(t, newBlock.Start, block.Start)
		currentBlockID, err := db.getCurrentBlockID()
		assert.NoError(t, err)
		assert.Equal(t, -1, currentBlockID)
	})
}

func TestStartPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.startPause()
		assert.Error(t, err)
	})

	t.Run("start successful", func(t *testing.T) {
		_, err := db.startBlock(false)
		assert.NoError(t, err)
		_, err = db.startPause()
		assert.NoError(t, err)
		currentPauseID, err := db.getCurrentPauseID()
		assert.NoError(t, err)
		assert.NotEqual(t, -1, currentPauseID)
	})

	t.Run("pause already active", func(t *testing.T) {
		_, err := db.startPause()
		assert.Error(t, err)
	})
}

func TestEndPause(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	t.Run("no pause active", func(t *testing.T) {
		_, err := db.startBlock(false)
		assert.NoError(t, err)
		_, err = db.endPause()
		assert.Error(t, err)
	})

	t.Run("end successful", func(t *testing.T) {
		_, err := db.startPause()
		assert.NoError(t, err)
		_, err = db.endPause()
		assert.NoError(t, err)
		currentPauseID, err := db.getCurrentPauseID()
		assert.NoError(t, err)
		assert.Equal(t, -1, currentPauseID)
	})

}

func TestGetCurrentBlock(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.getCurrentBlock()
		assert.Error(t, err)
	})

	t.Run("get successful", func(t *testing.T) {
		newBlock, err := db.startBlock(false)
		assert.NoError(t, err)
		block, err := db.getCurrentBlock()
		assert.NoError(t, err)
		assert.Equal(t, newBlock.Id, block.Id)
		assert.Equal(t, newBlock.Start, block.Start)
	})
}
