package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path"
	"time"

	"github.com/kilianmandscharo/work_hours/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func GetNewTestDatabase() *DB {
	db, err := NewTestDatabase()
	if err != nil {
		log.Fatalf("ERROR: could not open test database, %v", err)
	}
	err = db.Init()
	if err != nil {
		log.Fatalf("ERROR: could not initialize test database, %v", err)
	}
	return db
}

func NewTestDatabase() (*DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory&_foreign_keys=true")
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func NewDatabase() (*DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataPath := path.Join(homeDir, ".work_hours_data")

	if _, err := os.Stat(dataPath); err != nil {
		err = os.Mkdir(dataPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", path.Join(dataPath, "data.db?_foreign_keys=true"))
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (db *DB) Init() error {
	q := `
  CREATE TABLE IF NOT EXISTS block
  (id INTEGER PRIMARY KEY ASC,
  start TEXT,
  end TEXT,
  homeoffice INTEGER)
  `
	_, err := db.db.Exec(q)
	if err != nil {
		return err
	}
	q = `
  CREATE TABLE IF NOT EXISTS pause 
  (id INTEGER PRIMARY KEY ASC, 
  start TEXT, 
  end TEXT, 
  block_id INTEGER, 
  FOREIGN KEY(block_id) REFERENCES block(id) ON DELETE CASCADE)
  `
	_, err = db.db.Exec(q)
	if err != nil {
		return err
	}

	q = `
  CREATE TABLE IF NOT EXISTS current 
  (id INTEGER PRIMARY KEY ASC, 
  current_block_id INTEGER, 
  current_pause_id INTEGER)
  `
	_, err = db.db.Exec(q)
	if err != nil {
		return err
	}

	q = `
  INSERT OR IGNORE INTO current (id, current_block_id, current_pause_id)
  VALUES (1, -1, -1)
  `
	_, err = db.db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getBlocksFromRows(rows *sql.Rows) ([]models.Block, error) {
	var blocks []models.Block
	for rows.Next() {
		var b models.Block
		err := rows.Scan(&b.Id, &b.Start, &b.End, &b.Homeoffice)
		if err != nil {
			return nil, err
		}

		pauses, err := db.GetPausesByBlockID(b.Id)
		if err != nil {
			return nil, err
		}
		b.Pauses = pauses

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func (db *DB) GetBlocksAfterStart(start string) ([]models.Block, error) {
	q := `
  SELECT * FROM block
  WHERE start > date(?)
  `
	rows, err := db.db.Query(q, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks, err := db.getBlocksFromRows(rows)

	return blocks, nil
}

func (db *DB) GetBlocksBeforeEnd(end string) ([]models.Block, error) {
	q := `
  SELECT * FROM block
  WHERE end < date(?)
  `
	rows, err := db.db.Query(q, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks, err := db.getBlocksFromRows(rows)

	return blocks, nil
}

func (db *DB) GetBlocksWithinRange(start, end string) ([]models.Block, error) {
	q := `
  SELECT * FROM block
  WHERE start > date(?) AND end < date(?)
  `
	rows, err := db.db.Query(q, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks, err := db.getBlocksFromRows(rows)

	return blocks, nil
}

func (db *DB) GetAllBlocks() ([]models.Block, error) {
	q := `
  SELECT * FROM block
  `
	rows, err := db.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks, err := db.getBlocksFromRows(rows)

	return blocks, nil
}

func (db *DB) GetPausesByBlockID(blockID int) ([]models.Pause, error) {
	q := `
  SELECT * FROM pause
  WHERE block_id = ?
  `
	rows, err := db.db.Query(q, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pauses []models.Pause
	for rows.Next() {
		var p models.Pause
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.BlockID)
		if err != nil {
			return nil, err
		}

		pauses = append(pauses, p)
	}

	return pauses, nil
}

func (db *DB) GetBlockByID(id int) (models.Block, error) {
	q := `
  SELECT * FROM block
  WHERE id = ?
  `
	row := db.db.QueryRow(q, id)
	var b models.Block
	if err := row.Scan(&b.Id, &b.Start, &b.End, &b.Homeoffice); err != nil {
		return b, err
	}
	pauses, err := db.GetPausesByBlockID(b.Id)
	if err != nil {
		return b, err
	}
	b.Pauses = pauses
	return b, nil
}

func (db *DB) GetPauseByID(id int) (models.Pause, error) {
	q := `
  SELECT * FROM pause
  WHERE id = ?
  `
	row := db.db.QueryRow(q, id)
	var p models.Pause
	if err := row.Scan(&p.Id, &p.Start, &p.End, &p.BlockID); err != nil {
		return p, err
	}
	return p, nil
}

func (db *DB) AddBlock(block models.BlockCreate) (models.Block, error) {
	var newBlock models.Block
	q := `
  INSERT INTO block (start, end, homeoffice)
  VALUES (?, ?, ?)
  `
	result, err := db.db.Exec(q, block.Start, block.End, block.Homeoffice)
	if err != nil {
		return newBlock, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return newBlock, err
	}

	for _, pause := range block.Pauses {
		newPause, err := db.AddPause(
			models.PauseCreate{
				Start:   pause.Start,
				End:     pause.End,
				BlockID: int(id)})
		if err != nil {
			return newBlock, err
		}
		newBlock.Pauses = append(newBlock.Pauses, newPause)
	}

	newBlock.Id = int(id)
	newBlock.Start = block.Start
	newBlock.End = block.End
	return newBlock, nil
}

func (db *DB) AddPause(pause models.PauseCreate) (models.Pause, error) {
	var newPause models.Pause
	s := `
  INSERT INTO pause (start, end, block_id)
  VALUES (?, ?, ?)
  `
	result, err := db.db.Exec(s, pause.Start, pause.End, pause.BlockID)
	if err != nil {
		return newPause, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return newPause, err
	}

	newPause.Id = int(id)
	newPause.Start = pause.Start
	newPause.End = pause.End
	newPause.BlockID = pause.BlockID
	return newPause, nil
}

func (db *DB) DeleteBlock(id int) (int, error) {
	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return 0, err
	}
	if id == currentBlockID {
		err := db.setCurrentBlockID(-1)
		if err != nil {
			return 0, err
		}

		err = db.setCurrentPauseID(-1)
		if err != nil {
			return 0, err
		}
	}

	q := `
  DELETE FROM block
  WHERE id = ?
  `

	result, err := db.db.Exec(q, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), err
}

func (db *DB) DeletePause(id int) (int, error) {
	currentPauseID, err := db.getCurrentPauseID()
	if err != nil {
		return 0, err
	}
	if currentPauseID == id {
		err := db.setCurrentPauseID(-1)
		if err != nil {
			return 0, err
		}
	}

	q := `
  DELETE FROM pause
  WHERE id = ?
  `
	result, err := db.db.Exec(q, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), err
}

func (db *DB) UpdateBlock(block models.Block) (int, error) {
	q := `
  UPDATE block
  SET start = ?, end = ?, homeoffice = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, block.Start, block.End, block.Homeoffice, block.Id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdateBlockStart(id int, start string) (int, error) {
	q := `
  UPDATE block
  SET start = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, start, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdateBlockEnd(id int, end string) (int, error) {
	q := `
  UPDATE block
  SET end = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, end, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdateBlockHomeoffice(id int, homeoffice bool) (int, error) {
	q := `
  UPDATE block
  SET homeoffice = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, homeoffice, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdatePause(pause models.Pause) (int, error) {
	q := `
  UPDATE pause
  SET start = ?, end = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, pause.Start, pause.End, pause.Id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdatePauseStart(id int, start string) (int, error) {
	q := `
  UPDATE pause
  SET start = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, start, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) UpdatePauseEnd(id int, end string) (int, error) {
	q := `
  UPDATE pause
  SET end = ?
  WHERE id = ?
  `
	result, err := db.db.Exec(q, end, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (db *DB) getCurrentBlockID() (int, error) {
	q := `
  SELECT current_block_id from current
  WHERE id = 1
  `

	row := db.db.QueryRow(q)

	var id int

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (db *DB) setCurrentBlockID(id int) error {
	q := `
  UPDATE current
  SET current_block_id = ?
  WHERE id = 1
  `

	_, err := db.db.Exec(q, id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) getCurrentPauseID() (int, error) {
	q := `
  SELECT current_pause_id from current
  WHERE id = 1
  `

	row := db.db.QueryRow(q)

	var id int

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (db *DB) setCurrentPauseID(id int) error {
	q := `
  UPDATE current
  SET current_pause_id = ?
  WHERE id = 1
  `

	_, err := db.db.Exec(q, id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) StartBlock(homeoffice bool) (models.Block, error) {
	var newBlock models.Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return newBlock, err
	}
	if currentBlockID != -1 {
		return newBlock, errors.New("current block already active")
	}

	block := models.BlockCreate{
		Start:      time.Now().Format(time.RFC3339),
		Homeoffice: homeoffice,
	}
	newBlock, err = db.AddBlock(block)
	if err != nil {
		return newBlock, err
	}

	err = db.setCurrentBlockID(newBlock.Id)
	if err != nil {
		return newBlock, err
	}

	return newBlock, nil
}

func (db *DB) EndBlock() (models.Block, error) {
	var block models.Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return block, err
	}
	if currentBlockID == -1 {
		return block, errors.New("no current block active")
	}

	currentPauseID, err := db.getCurrentPauseID()
	if err != nil {
		return block, err
	}
	if currentPauseID != -1 {
		return block, errors.New("pause not ended")
	}

	q := `
  UPDATE block
  SET end = ?
  WHERE id = ?
  `
	end := time.Now().Format(time.RFC3339)
	_, err = db.db.Exec(q, end, currentBlockID)
	if err != nil {
		return block, err
	}

	block, err = db.GetBlockByID(currentBlockID)
	if err != nil {
		return block, err
	}

	err = db.setCurrentBlockID(-1)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (db *DB) GetCurrentBlock() (models.Block, error) {
	var block models.Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return block, err
	}

	block, err = db.GetBlockByID(currentBlockID)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (db *DB) StartPause() (models.Pause, error) {
	var newPause models.Pause

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return newPause, err
	}
	if currentBlockID == -1 {
		return newPause, errors.New("no current block active")
	}

	currentPauseID, err := db.getCurrentPauseID()
	if err != nil {
		return newPause, err
	}
	if currentPauseID != -1 {
		return newPause, errors.New("current pause already active")
	}

	pause := models.PauseCreate{
		Start:   time.Now().Format(time.RFC3339),
		BlockID: currentBlockID,
	}
	newPause, err = db.AddPause(pause)
	if err != nil {
		return newPause, err
	}

	err = db.setCurrentPauseID(newPause.Id)
	if err != nil {
		return newPause, err
	}

	return newPause, nil
}

func (db *DB) EndPause() (models.Pause, error) {
	var pause models.Pause

	currentPauseID, err := db.getCurrentPauseID()
	if err != nil {
		return pause, err
	}
	if currentPauseID == -1 {
		return pause, errors.New("no current pause active")
	}

	q := `
  UPDATE pause
  SET end = ?
  WHERE id = ?
  `
	end := time.Now().Format(time.RFC3339)
	_, err = db.db.Exec(q, end, currentPauseID)
	if err != nil {
		return pause, err
	}

	pause, err = db.GetPauseByID(currentPauseID)
	if err != nil {
		return pause, err
	}

	err = db.setCurrentPauseID(-1)
	if err != nil {
		return pause, err
	}

	return pause, nil
}
