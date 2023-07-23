package database

import (
	"testing"

	"github.com/kilianmandscharo/work_hours/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitialCurrentIDs(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, -1, currentBlockID)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, -1, currentPauseID)
}

func TestSetCurrentBlockID(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	testID := 10

	err := db.setCurrentBlockID(testID)
	assert.NoError(t, err)
	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, testID, currentBlockID)
}

func TestSetCurrentPauseID(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	testID := 10

	err := db.setCurrentPauseID(testID)
	assert.NoError(t, err)
	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, testID, currentPauseID)
}

func TestPragma(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()
	row := db.db.QueryRow("PRAGMA foreign_keys")
	assert.NotNil(t, row)
	var foreignKeys int
	if err := row.Scan(&foreignKeys); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, 1, foreignKeys)
}

func TestAddBlockWithPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	b, err := db.AddBlock(utils.TestBlockCreate())
	assert.NoError(t, err)
	utils.AssertTestBlock(t, b)
	assert.Equal(t, 1, len(b.Pauses))
	utils.AssertTestPause(t, b.Pauses[0])
}

func TestAddBlockWithoutPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	b, err := db.AddBlock(utils.TestBlockCreateWithoutPause())
	assert.NoError(t, err)
	utils.AssertTestBlock(t, b)
	assert.Equal(t, 0, len(b.Pauses))
}

func TestAddPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	_, err := db.AddPause(utils.TestPauseCreate())
	assert.Error(t, err)

	db.AddBlock(utils.TestBlockCreateWithoutPause())

	p, err := db.AddPause(utils.TestPauseCreate())
	assert.NoError(t, err)
	utils.AssertTestPause(t, p)
}

func TestGetBlockByID(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	_, err := db.GetBlockByID(utils.BID)
	assert.Error(t, err)

	db.AddBlock(utils.TestBlockCreate())

	b, err := db.GetBlockByID(utils.BID)
	assert.NoError(t, err)
	utils.AssertTestBlock(t, b)
	assert.Equal(t, 1, len(b.Pauses))
	utils.AssertTestPause(t, b.Pauses[0])
}

func TestGetPauseByID(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	_, err := db.GetPauseByID(utils.PID)
	assert.Error(t, err)

	db.AddBlock(utils.TestBlockCreate())

	p, err := db.GetPauseByID(utils.PID)
	assert.NoError(t, err)
	utils.AssertTestPause(t, p)
}

func TestGetBlocksWithinRange(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	testBlocks := utils.CreateRangeTestBlocks()

	for _, block := range testBlocks {
		db.AddBlock(block)
	}

	testCases := []struct {
		start  string
		end    string
		length int
		id     int
	}{
		{
			start:  "2023-01-01T07:00:00Z",
			end:    "2023-01-31T07:00:00Z",
			length: 0,
			id:     -1,
		},
		{
			start:  "2023-05-01T07:00:00Z",
			end:    "2023-05-31T07:00:00Z",
			length: 1,
			id:     1,
		},
		{
			start:  "2023-06-01T07:00:00Z",
			end:    "2023-06-30T07:00:00Z",
			length: 1,
			id:     2,
		},
		{
			start:  "2023-07-01T07:00:00Z",
			end:    "2023-07-31T07:00:00Z",
			length: 1,
			id:     3,
		},
		{
			start:  "2023-05-01T07:00:00Z",
			end:    "2023-07-31T07:00:00Z",
			length: 3,
			id:     -1,
		},
	}

	for _, testCase := range testCases {
		blocks, err := db.GetBlocksWithinRange(
			testCase.start,
			testCase.end,
		)
		assert.NoError(t, err)
		assert.Equal(t, testCase.length, len(blocks))

		if testCase.id > 0 && testCase.length > 0 {
			assert.Equal(t, testCase.id, blocks[0].Id)
		}
	}
}

func TestGetBlocksAfterStart(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	testBlocks := utils.CreateRangeTestBlocks()

	for _, block := range testBlocks {
		db.AddBlock(block)
	}

	testCases := []struct {
		start  string
		length int
	}{
		{
			start:  "2023-08-01T07:00:00Z",
			length: 0,
		},
		{
			start:  "2023-07-01T07:00:00Z",
			length: 1,
		},
		{
			start:  "2023-06-01T07:00:00Z",
			length: 2,
		},
		{
			start:  "2023-05-01T07:00:00Z",
			length: 3,
		},
	}

	for _, testCase := range testCases {
		blocks, err := db.GetBlocksAfterStart(
			testCase.start,
		)
		assert.NoError(t, err)
		assert.Equal(t, testCase.length, len(blocks))
	}
}

func TestGetBlocksBeforeEnd(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	testBlocks := utils.CreateRangeTestBlocks()

	for _, block := range testBlocks {
		db.AddBlock(block)
	}

	testCases := []struct {
		start  string
		length int
	}{
		{
			start:  "2023-05-01T07:00:00Z",
			length: 0,
		},
		{
			start:  "2023-05-31T07:00:00Z",
			length: 1,
		},
		{
			start:  "2023-06-30T07:00:00Z",
			length: 2,
		},
		{
			start:  "2023-07-31T07:00:00Z",
			length: 3,
		},
	}

	for _, testCase := range testCases {
		blocks, err := db.GetBlocksBeforeEnd(
			testCase.start,
		)
		assert.NoError(t, err)
		assert.Equal(t, testCase.length, len(blocks))
	}
}

func TestGetAllBlocks(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	blocks, err := db.GetAllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(blocks))

	db.AddBlock(utils.TestBlockCreate())

	blocks, err = db.GetAllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(blocks))
	utils.AssertTestBlock(t, blocks[0])
	assert.Equal(t, 1, len(blocks[0].Pauses))
	utils.AssertTestPause(t, blocks[0].Pauses[0])
}

func TestDeleteBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.DeleteBlock(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.DeleteBlock(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)
	_, err = db.GetPauseByID(utils.PID)
	assert.Error(t, err)
	_, err = db.GetBlockByID(utils.BID)
	assert.Error(t, err)
}

func TestDeleteCurrentBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	_, err := db.StartBlock(false)
	assert.NoError(t, err)
	_, err = db.StartPause()
	assert.NoError(t, err)

	rowsAffected, err := db.DeleteBlock(utils.BID)
	assert.Equal(t, rowsAffected, 1)

	currentBlockID, err := db.getCurrentBlockID()
	assert.NoError(t, err)
	assert.Equal(t, currentBlockID, -1)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, currentPauseID, -1)
}

func TestDeletePause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.DeletePause(utils.PID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.DeletePause(utils.PID)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)
	_, err = db.GetPauseByID(utils.PID)
	assert.Error(t, err)
}

func TestDeleteCurrentPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	_, err := db.StartBlock(false)
	assert.NoError(t, err)
	_, err = db.StartPause()
	assert.NoError(t, err)

	rowsAffected, err := db.DeletePause(utils.PID)
	assert.Equal(t, rowsAffected, 1)

	currentPauseID, err := db.getCurrentPauseID()
	assert.NoError(t, err)
	assert.Equal(t, currentPauseID, -1)
}

func TestUpdateBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdateBlock(utils.TestBlockUpdated())
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdateBlock(utils.TestBlockUpdated())
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	b, err := db.GetBlockByID(utils.BID)
	assert.NoError(t, err)
	utils.AssertTestBlockUpdated(t, b)
}

func TestUpdateBlockStart(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdateBlockStart(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdateBlockStart(utils.BID, utils.BStartUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.GetBlockByID(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, block.Start, utils.BStartUpdated)
}

func TestUpdateBlockEnd(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdateBlockEnd(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdateBlockEnd(utils.BID, utils.BEndUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.GetBlockByID(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, block.End, utils.BEndUpdated)
}

func TestUpdateBlockHomeoffice(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdateBlockHomeoffice(10, true)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdateBlockHomeoffice(utils.BID, utils.BHomeofficeUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	block, err := db.GetBlockByID(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, block.Homeoffice, utils.BHomeofficeUpdated)
}

func TestUpdatePause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdatePause(utils.TestPauseUpdated())
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdatePause(utils.TestPauseUpdated())
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	p, err := db.GetPauseByID(utils.PID)
	utils.AssertTestPauseUpdated(t, p)
}

func TestUpdatePauseStart(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdatePauseStart(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdatePauseStart(utils.BID, utils.PStartUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	pause, err := db.GetPauseByID(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, pause.Start, utils.PStartUpdated)
}

func TestUpdatePauseEnd(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	rowsAffected, err := db.UpdatePauseEnd(10, "test")
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 0)

	db.AddBlock(utils.TestBlockCreate())

	rowsAffected, err = db.UpdatePauseEnd(utils.BID, utils.PEndUpdated)
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, 1)

	pause, err := db.GetPauseByID(utils.BID)
	assert.NoError(t, err)
	assert.Equal(t, pause.End, utils.PEndUpdated)
}

func TestStartBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	t.Run("start successful", func(t *testing.T) {
		_, err := db.StartBlock(false)
		assert.NoError(t, err)
		currentBlockID, err := db.getCurrentBlockID()
		assert.NoError(t, err)
		assert.NotEqual(t, -1, currentBlockID)
	})

	t.Run("block already active", func(t *testing.T) {
		_, err := db.StartBlock(false)
		assert.Error(t, err)
	})
}

func TestEndBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.EndBlock()
		assert.Error(t, err)
	})

	t.Run("end successful", func(t *testing.T) {
		newBlock, err := db.StartBlock(false)
		assert.NoError(t, err)
		block, err := db.EndBlock()
		assert.NoError(t, err)
		assert.Equal(t, newBlock.Id, block.Id)
		assert.Equal(t, newBlock.Start, block.Start)
		currentBlockID, err := db.getCurrentBlockID()
		assert.NoError(t, err)
		assert.Equal(t, -1, currentBlockID)
	})
}

func TestStartPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.StartPause()
		assert.Error(t, err)
	})

	t.Run("start successful", func(t *testing.T) {
		_, err := db.StartBlock(false)
		assert.NoError(t, err)
		_, err = db.StartPause()
		assert.NoError(t, err)
		currentPauseID, err := db.getCurrentPauseID()
		assert.NoError(t, err)
		assert.NotEqual(t, -1, currentPauseID)
	})

	t.Run("pause already active", func(t *testing.T) {
		_, err := db.StartPause()
		assert.Error(t, err)
	})
}

func TestEndPause(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	t.Run("no pause active", func(t *testing.T) {
		_, err := db.StartBlock(false)
		assert.NoError(t, err)
		_, err = db.EndPause()
		assert.Error(t, err)
	})

	t.Run("end successful", func(t *testing.T) {
		_, err := db.StartPause()
		assert.NoError(t, err)
		_, err = db.EndPause()
		assert.NoError(t, err)
		currentPauseID, err := db.getCurrentPauseID()
		assert.NoError(t, err)
		assert.Equal(t, -1, currentPauseID)
	})

}

func TestGetCurrentBlock(t *testing.T) {
	db := GetNewTestDatabase()
	defer db.Close()

	t.Run("no block active", func(t *testing.T) {
		_, err := db.GetCurrentBlock()
		assert.Error(t, err)
	})

	t.Run("get successful", func(t *testing.T) {
		newBlock, err := db.StartBlock(false)
		assert.NoError(t, err)
		block, err := db.GetCurrentBlock()
		assert.NoError(t, err)
		assert.Equal(t, newBlock.Id, block.Id)
		assert.Equal(t, newBlock.Start, block.Start)
	})
}
