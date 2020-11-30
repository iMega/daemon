package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

type fakerDB struct{}

const errFaker = "it is a faker database"

func (f *fakerDB) Begin() (*sql.Tx, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) Close() error {
	return nil
}

func (f *fakerDB) Conn(context.Context) (*sql.Conn, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) Driver() driver.Driver {
	return nil
}

func (f *fakerDB) Exec(string, ...interface{}) (sql.Result, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) Ping() error {
	return errors.New(errFaker)
}

func (f *fakerDB) PingContext(context.Context) error {
	return errors.New(errFaker)
}

func (f *fakerDB) Prepare(string) (*sql.Stmt, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) Query(string, ...interface{}) (*sql.Rows, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	return nil, errors.New(errFaker)
}

func (f *fakerDB) QueryRow(string, ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (f *fakerDB) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (f *fakerDB) SetConnMaxIdleTime(time.Duration) {}

func (f *fakerDB) SetConnMaxLifetime(time.Duration) {}

func (f *fakerDB) SetMaxIdleConns(int) {}

func (f *fakerDB) SetMaxOpenConns(int) {}

func (f *fakerDB) Stats() sql.DBStats {
	return sql.DBStats{}
}
