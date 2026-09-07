package database

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"cth.release/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *common.Config) (*pgxpool.Pool, error) {
	port, err := strconv.Atoi(cfg.Database.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid database port %q: %w", cfg.Database.Port, err)
	}

	dsn := (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Database.User, cfg.Database.Pass),
		Host:   fmt.Sprintf("%s:%d", cfg.Database.Host, port),
		Path:   cfg.Database.Name,
	}).String()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
