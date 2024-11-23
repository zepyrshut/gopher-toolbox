package db

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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

	if err := d.Ping(); err != nil {
		slog.Error("error pinging database", "error", err)
		return nil, err
	}

	return d, nil
}
