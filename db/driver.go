package db

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

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

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = time.Minute * 5

func NewMySQL(dataSource string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dataSource)
	if err != nil {
		slog.Error("error connecting to database", "error", err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	err = testDB(d)
	if err != nil {
		slog.Error("error pinging database", "error", err)
		return nil, err
	}

	return d, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
