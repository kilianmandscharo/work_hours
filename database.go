package main

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func newTestDatabase() (*DB, error) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory&_foreign_keys=true")
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func newDatabase() (*DB, error) {
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

func (db *DB) init() error {
	q := `
  CREATE TABLE IF NOT EXISTS block
  (id INTEGER PRIMARY KEY ASC,
  start TEXT,
  end TEXT)
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

	//_, err = db.db.Exec("INSERT OR IGNORE INTO ui (id, list_order) VALUES(1, '')")
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

func (db *DB) close() error {
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getAllBlocks() ([]Block, error) {
	q := `
  SELECT * FROM block
  `
	rows, err := db.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var b Block
		err = rows.Scan(&b.Id, &b.Start, &b.End)
		if err != nil {
			return nil, err
		}

		pauses, err := db.getPausesByBlockID(b.Id)
		if err != nil {
			return nil, err
		}
		b.Pauses = pauses

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func (db *DB) getPausesByBlockID(blockID int) ([]Pause, error) {
	q := `
  SELECT * FROM pause
  WHERE block_id = ?
  `
	rows, err := db.db.Query(q, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pauses []Pause
	for rows.Next() {
		var p Pause
		err = rows.Scan(&p.Id, &p.Start, &p.End, &p.BlockID)
		if err != nil {
			return nil, err
		}

		pauses = append(pauses, p)
	}

	return pauses, nil
}

func (db *DB) getBlockByID(id int) (Block, error) {
	q := `
  SELECT * FROM block
  WHERE id = ?
  `
	row := db.db.QueryRow(q, id)
	var b Block
	if err := row.Scan(&b.Id, &b.Start, &b.End); err != nil {
		return b, err
	}
	pauses, err := db.getPausesByBlockID(b.Id)
	if err != nil {
		return b, err
	}
	b.Pauses = pauses
	return b, nil
}

func (db *DB) getPauseByID(id int) (Pause, error) {
	q := `
  SELECT * FROM pause
  WHERE id = ?
  `
	row := db.db.QueryRow(q, id)
	var p Pause
	if err := row.Scan(&p.Id, &p.Start, &p.End, &p.BlockID); err != nil {
		return p, err
	}
	return p, nil
}

func (db *DB) addBlock(block BlockCreate) (Block, error) {
	var newBlock Block
	q := `
  INSERT INTO block (start, end)
  VALUES (?, ?)
  `
	result, err := db.db.Exec(q, block.Start, block.End)
	if err != nil {
		return newBlock, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return newBlock, err
	}

	for _, pause := range block.Pauses {
		newPause, err := db.addPause(
			PauseCreate{
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

func (db *DB) addPause(pause PauseCreate) (Pause, error) {
	var newPause Pause
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

func (db *DB) deleteBlock(id int) error {
	q := `
  DELETE FROM block
  WHERE id = ?
  `
	_, err := db.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) deletePause(id int) error {
	q := `
  DELETE FROM pause
  WHERE id = ?
  `
	_, err := db.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updateBlock(block Block) error {
	q := `
  UPDATE block
  SET start = ?, end = ?
  WHERE id = ?
  `
	_, err := db.db.Exec(q, block.Start, block.End, block.Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) updatePause(pause Pause) error {
	q := `
  UPDATE pause
  SET start = ?, end = ?
  WHERE id = ?
  `
	_, err := db.db.Exec(q, pause.Start, pause.End, pause.Id)
	if err != nil {
		return err
	}
	return nil
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

func (db *DB) startBlock() (Block, error) {
	var newBlock Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return newBlock, err
	}
	if currentBlockID != -1 {
		return newBlock, errors.New("current block already active")
	}

	block := BlockCreate{
		Start: time.Now().Format(time.RFC3339),
	}
	newBlock, err = db.addBlock(block)
	if err != nil {
		return newBlock, err
	}

	err = db.setCurrentBlockID(newBlock.Id)
	if err != nil {
		return newBlock, err
	}

	return newBlock, nil
}

func (db *DB) endBlock() (Block, error) {
	var block Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return block, err
	}
	if currentBlockID == -1 {
		return block, errors.New("no current block active")
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

	block, err = db.getBlockByID(currentBlockID)
	if err != nil {
		return block, err
	}

	err = db.setCurrentBlockID(-1)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (db *DB) getCurrentBlock() (Block, error) {
	var block Block

	currentBlockID, err := db.getCurrentBlockID()
	if err != nil {
		return block, err
	}

	block, err = db.getBlockByID(currentBlockID)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (db *DB) startPause() (Pause, error) {
  var newPause Pause

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

	pause := PauseCreate{
		Start:   time.Now().Format(time.RFC3339),
		BlockID: currentBlockID,
	}
	newPause, err = db.addPause(pause)
	if err != nil {
		return newPause, err
	}

	err = db.setCurrentPauseID(newPause.Id)
	if err != nil {
		return newPause, err
	}

	return newPause, nil
}

func (db *DB) endPause() (Pause, error) {
	var pause Pause

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

	pause, err = db.getPauseByID(currentPauseID)
	if err != nil {
		return pause, err
	}

	err = db.setCurrentPauseID(-1)
	if err != nil {
		return pause, err
	}

	return pause, nil
}
