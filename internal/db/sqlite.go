package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Enable WAL
	if _, err := conn.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		return nil, err
	}

	// Optional but recommended
	if _, err := conn.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	s := &SQLite{conn: conn}

	// if err := s.initSchema(); err != nil {
	// 	return nil, err
	// }

	return s, nil
}

func (s *SQLite) Conn() *sql.DB {
	return s.conn
}

func (s *SQLite) Close() error {
	return s.conn.Close()
}

// func (s *SQLite) initSchema() error {
// 	schema := `
// 	CREATE TABLE IF NOT EXISTS routes (
// 		id TEXT PRIMARY KEY,
// 		origin_id TEXT NOT NULL,
// 		destination_id TEXT NOT NULL,
// 		distance_meters INTEGER NOT NULL,
// 		duration_seconds INTEGER NOT NULL,
// 		recorded_at DATETIME NOT NULL
// 	);
// 	`
// 	_, err := s.conn.Exec(schema)
// 	return err
// }
