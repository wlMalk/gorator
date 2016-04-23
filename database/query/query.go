package query

import (
	"strings"
)

type Query interface {
	Name(string)
	GetName() string
	ToSql() (string, []interface{}, error)
}

func Tokenize(q Query) (t *Token, err error) {
	a, ok := q.(Aliased)
	if ok {
		t.isQuery = true
		t.alias = a.GetAlias()
	}
	t.str, t.args, err = q.ToSql()
	return
}

func TokenizeAll(c ...interface{}) ([]interface{}, [][]interface{}, error) {
	var args [][]interface{}
	for i, a := range c {
		switch a.(type) {
		case string:
			continue
		case *Token:
			at := a.(*Token)
			t, targs := at.Token()
			c[i] = t
			args = append(args, targs)
		case Query:
			q := a.(Query)
			at, err := Tokenize(q)
			if err != nil {
				return nil, nil, err
			}
			t, targs := at.Token()
			c[i] = t
			args = append(args, targs)
		}
	}
	return c, args, nil
}

type Extra map[string]interface{}

func (e Extra) Get(k string) interface{} {
	return e[k]
}

type Aliased interface {
	Query
	GetAlias() string
}

func Set(k string, v interface{}) *Token {
	return T(k+" = ?", v)
}

func Values(v ...interface{}) *Token {
	return T("(?"+strings.Repeat(",?", len(v)-1)+")", v...)
}

type token struct {
	id      int
	isQuery bool
	alias   string
	str     string
}

type Token struct {
	str     string
	args    []interface{}
	isQuery bool
	alias   string
}

func (t *Token) String() string {
	return t.str
}

func (t *Token) Token() (*token, []interface{}) {
	return &token{str: t.str, isQuery: t.isQuery, alias: t.alias}, t.args
}

func T(str string, args ...interface{}) *Token {
	return &Token{str, args, false, ""}
}

type Preload struct {
	Name     string
	Q        *SelectQuery
	Preloads []*Preload
}

func P(name string, q *SelectQuery, preloads []*Preload) *Preload {
	return &Preload{name, q, preloads}
}
