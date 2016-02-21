package database

import (
	"database/sql"
)

type Driver interface {
	Queryer
	QueryRower
	Execer
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
	InsertIds() ([]interface{}, error)
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

type DB struct {
	*sql.DB
	name string
}

func New() *DB {
	return &DB{}
}

func (db *DB) SetName(name string) {
	db.name = name
}
