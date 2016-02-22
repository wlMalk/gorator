package database

import (
	"database/sql"
)

type Driver interface {
	Queryer
	QueryRower
	Execer
	Close() error
	Ping() error
	Prepare(string) (*sql.Stmt, error)
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}

type Transaction interface {
	Queryer
	QueryRower
	Execer
	Rollback() error
	Commit() error
	Prepare(string) (*sql.Stmt, error)
	Stmt(*sql.Stmt) *sql.Stmt
}

type Model interface {
}

type Slice interface {
}

type Queryer interface {
	Query(string, ...interface{}) (Rower, error)
}

type QueryRower interface {
	QueryRow(string, ...interface{}) (Rower, error)
}

type Execer interface {
	Exec(string, ...interface{}) (Result, error)
}

type Scanner interface {
	Scan(...interface{}) error
}

type Rower interface {
	Scanner
	Next() bool
	Columns() ([]string, error)
	Close() error
}

type Result interface {
	sql.Result
	// InsertIds() ([]interface{}, error)
}

type Tabler interface {
	Driver
	Name() string
	Columns(string) string
}

type Table struct {
	name    string
	schema  string
	columns map[string]string
}

func (t *Table) SetName(name string) {
	if t.name == "" {
		t.name = name
	}
}

func (t *Table) SetSchema(schema string) {
	if t.schema == "" {
		t.schema = schema
	}
}

func (t *Table) SetColumns(columns map[string]string) {
	if t.columns == nil || len(t.columns) == 0 {
		t.columns = columns
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Column(a string) string {
	return t.columns[a]
}

func (t *Table) Columns(a ...string) []string {
	var b []string

	for _, c := range a {
		b = append(b, t.columns[c])
	}

	return b
}

type DB struct {
	*sql.DB
	name string
}

type Tx struct {
	*sql.Tx
}

func New() *DB {
	return &DB{}
}

func (db *DB) SetName(name string) {
	db.name = name
}

func (db *DB) Query(query string, args ...interface{}) (Rower, error) {
	return db.DB.Query(query, args...)
}

func (db *DB) Begin() (Transaction, error) {
	xtx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{xtx}, err
}

func (db *DB) QueryRow(query string, args ...interface{}) (Rower, error) {
	return db.DB.Query(query, args...)
}

func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	return db.DB.Exec(query, args...)
}

func (tx *Tx) Query(query string, args ...interface{}) (Rower, error) {
	return tx.Tx.Query(query, args...)
}

func (tx *Tx) QueryRow(query string, args ...interface{}) (Rower, error) {
	return tx.Tx.Query(query, args...)
}

func (tx *Tx) Exec(query string, args ...interface{}) (Result, error) {
	return tx.Tx.Exec(query, args...)
}
