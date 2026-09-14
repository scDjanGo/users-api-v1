package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed users.sql
var schema string

func Open(path string) (*sql.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"

	conn, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)

	if _, err := conn.Exec(schema); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
