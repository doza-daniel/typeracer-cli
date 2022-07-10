package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"play/game"
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
func (db *DB) GetTextAt(int64) game.Text {
	text, err := db.GetRandomText()
	if err != nil || text.Content == "" {
		return game.Text{
			Content: "The quick brown fox jumps over the lazy dog.",
			Source:  "Mock_Source",
			Type:    "Mock_Type",
			Author:  "Mock_Author",
		}
	}

	return text
}

// Size ... mocked for now
func (db *DB) Size() int64 {
	return 42
}

// GetRandomText ...
func (db *DB) GetRandomText() (game.Text, error) {
	id, err := db.getRandID()
	if err != nil {
		return game.Text{}, err
	}

	row := db.db.QueryRow("SELECT text, type, source, author FROM texts WHERE id = ?", id)

	var text game.Text
	switch err = row.Scan(&text.Content, &text.Type, &text.Source, &text.Author); err {
	case sql.ErrNoRows:
		return game.Text{}, fmt.Errorf("text with id %d not found", id)
	case nil:
		return text, nil
	default:
		return game.Text{}, err
	}
}

func (db *DB) getRandID() (int64, error) {
	rows, err := db.db.Query("SELECT id FROM texts ORDER BY id DESC")
	if err != nil {
		return 0, err
	}

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}

		ids = append(ids, id)
	}

	i := rand.Intn(len(ids))

	return ids[i], nil
}
