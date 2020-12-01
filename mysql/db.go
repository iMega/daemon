package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"
)

// DB is a sql interface.
type DB interface {
	Begin() (*sql.Tx, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	Ping() error
	PingContext(context.Context) error
	Prepare(string) (*sql.Stmt, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	SetConnMaxIdleTime(time.Duration)
	SetConnMaxLifetime(time.Duration)
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
	Stats() sql.DBStats
}
