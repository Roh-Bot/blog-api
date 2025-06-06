package database

import (
	"context"
	"fmt"
	"github.com/Roh-Bot/blog-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var dbPool *pgxpool.Pool

func New(config config.Database) (pool *pgxpool.Pool, err error) {
	connString := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s database=%s sslmode=%s`,
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return
	}
	cfg.MaxConns = config.MaxConnections
	cfg.MaxConnIdleTime = time.Minute * time.Duration(config.MaxConnectionIdleTime)
	cfg.MaxConnLifetime = time.Minute * time.Duration(config.MaxConnectionLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return
	}

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	dbPool = pool
	return pool, nil
}

func Flush() {
	dbPool.Close()
}
