package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

const defaultDBFile = "scheduler.db"

var db *sql.DB

// File возвращает путь к файлу БД: из TODO_DBFILE или значение по умолчанию.
func File() string {
	if f := os.Getenv("TODO_DBFILE"); f != "" {
		return f
	}
	return defaultDBFile
}

// Init открывает базу данных и при отсутствии файла создаёт таблицу и индекс.
func Init(dbFile string) error {
	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		install = true
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

// DB возвращает подключение к базе данных.
func DB() *sql.DB {
	return db
}

// Close закрывает подключение к базе данных.
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
