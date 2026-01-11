package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var (
	schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(256) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`
	db *sql.DB
)

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	if install == true {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
