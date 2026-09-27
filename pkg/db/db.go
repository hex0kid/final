package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return err
	}
	if _, err = db.Exec(schema); err != nil {
		_ = db.Close()
		return err
	}
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
