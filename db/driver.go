package db

import (
	"context"
	"log/slog"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresPool(dataSource string) *pgxpool.Pool {
	dbPool, err := pgxpool.New(context.Background(), dataSource)
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		panic(err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		slog.Error("error pinging database, maybe incorrect datasource", "error", err)
		panic(err)
	}

	slog.Info("connected to database")
	return dbPool
}
