package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Open(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}

	// Important pragmas
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA busy_timeout=5000;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatal(err)
		}
	}

	return db
}
