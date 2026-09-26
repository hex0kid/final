package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	install := false
	if _, err := os.Stat(dbFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return err
		}
	}

	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			_ = db.Close()
			return err
		}
	}
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
