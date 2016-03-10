package query

import (
	"strings"
)

func Eq(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " = " + acol
		return c
	}
	if arg == nil {
		c.str = col + " IS NULL"
		return c
	}
	c.str = col + " = ?"
	c.args = append(c.args, arg)
	return c
}

func NotEq(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " <> " + acol
		return c
	}
	if arg == nil {
		c.str = col + " IS NOT NULL"
		return c
	}
	c.str = col + " <> ?"
	c.args = append(c.args, arg)
	return c
}

func Gt(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " > " + acol
		return c
	}
	if arg == nil {
		return c
	}
	c.str = col + " > ?"
	c.args = append(c.args, arg)
	return c
}

func GtEq(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " >= " + acol
		return c
	}
	if arg == nil {
		return c
	}
	c.str = col + " >= ?"
	c.args = append(c.args, arg)
	return c
}

func Lt(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " < " + acol
		return c
	}
	if arg == nil {
		return c
	}
	c.str = col + " < ?"
	c.args = append(c.args, arg)
	return c
}

func LtEq(col string, acol string, arg interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	if len(acol) > 0 {
		c.str = col + " <= " + acol
		return c
	}
	if arg == nil {
		return c
	}
	c.str = col + " <= ?"
	c.args = append(c.args, arg)
	return c
}

func In(col string, args ...interface{}) *Token {
	c := &Token{}
	if len(col) == 0 {
		return c
	}
	l := len(args)
	if l == 0 {
		return c
	}
	c.str = col + " IN (" + strings.Repeat("?,", l)[:l-1] + ")"
	c.args = append(c.args, args...)
	return c
}

type Sqler interface {
	ToSql() (string, []interface{})
}

func InQuery(col string, q Query) *Token {
	c := &Token{}
	if q == nil {
		return c
	}
	sql, args, _ := q.ToSql()
	c.str = col + " IN (" + sql + ")"
	c.args = append(c.args, args...)
	return c
}

func And(con ...*Token) *Token {
	c := &Token{}
	if len(con) == 0 {
		return c
	}
	var strs []string
	for i := range con {
		strs = append(strs, con[i].str)
		c.args = append(c.args, con[i].args...)
	}
	str := strings.Join(strs, " AND ")
	c.str = "(" + str + ")"
	return c
}

func Or(con ...*Token) *Token {
	c := &Token{}
	if len(con) == 0 {
		return c
	}
	var strs []string
	for i := range con {
		strs = append(strs, con[i].str)
		c.args = append(c.args, con[i].args...)
	}
	str := strings.Join(strs, " OR ")
	c.str = "(" + str + ")"
	return c
}
