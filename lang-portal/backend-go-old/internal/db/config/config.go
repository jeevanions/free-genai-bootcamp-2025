package config

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Config struct {
	DatabasePath string
}

func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
