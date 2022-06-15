package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	// blank
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// New ...
func New(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// DB ...
type DB struct {
	db *sql.DB
}

// GetTextAt ... should implement correct way of handing this
func (db *DB) GetTextAt(int64) string {
	text, err := db.GetRandomText()
	if err != nil || text == "" {
		return "The quick brown fox jumps over the lazy dog."
	}

	return text
}

// Size ... mocked for now
func (db *DB) Size() int64 {
	return 42
}

// GetRandomText ...
func (db *DB) GetRandomText() (string, error) {
	rows, err := db.db.Query("SELECT id FROM texts ORDER BY id DESC")
	if err != nil {
		return "", err
	}

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return "", err
		}

		ids = append(ids, id)
	}

	i := rand.Intn(len(ids))
	fmt.Println(i, ids[i])

	row := db.db.QueryRow("SELECT text FROM texts WHERE id = ?", ids[i])

	var text string
	switch err = row.Scan(&text); err {
	case sql.ErrNoRows:
		return "The quick brown fox jumps over the lazy dog.", nil
	case nil:
		return text, nil
	default:
		return "", err
	}
}
