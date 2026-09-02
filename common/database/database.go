package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"cth.release/common"
	_ "github.com/lib/pq"
)

func Connect(cfg *common.Config) (*sql.DB, error) {
	port, err := strconv.Atoi(cfg.Database.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid database port %q: %w", cfg.Database.Port, err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
