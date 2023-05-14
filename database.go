package main

import (
	"database/sql"
	"os"
	"path"

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
