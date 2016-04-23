package query

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var selects map[string]*SelectQuery = map[string]*SelectQuery{}
var selectsMutex *sync.Mutex = &sync.Mutex{}

func GetSelectQuery(name string) (*SelectQuery, bool) {
	selectsMutex.Lock()
	q, ok := selects[name]
	selectsMutex.Unlock()
	return q, ok
}

func SetSelectQuery(q *SelectQuery) error {
	if q == nil {
		return fmt.Errorf("query is nil")
	} else if q.name == "" {
		return fmt.Errorf("query must have a name")
	}

	selectsMutex.Lock()
	selects[q.name] = q
	selectsMutex.Unlock()

	return nil
}

func As(a interface{}, as string) *Token {
	c := &Token{}
	switch b := a.(type) {
	case *Token:
		c.str = b.str + " AS " + as
		c.args = b.args
	case string:
		c.str = b + " AS " + as
	}
	return c
}

func Over(a interface{}, over interface{}) *Token {
	c := &Token{}
	switch b := a.(type) {
	case *Token:
		c.str = b.str + " OVER"
		c.args = append(c.args, b.args...)
	case string:
		c.str = b + " OVER"
	}
	switch b := over.(type) {
	case *Token:
		c.str = "(" + b.str + ")"
		c.args = append(c.args, b.args...)
	case Query:
		sql, args, err := b.ToSql()
		if err != nil {
			return nil
		}
		c.str = "(" + sql + ")"
		c.args = append(c.args, args...)
	case string:
		c.str = "(" + b + ")"
	}
	return c
}

func OverWindow(a interface{}, window string) *Token {
	c := &Token{}
	switch b := a.(type) {
	case *Token:
		c.str = b.str + " OVER " + window
		c.args = b.args
	case string:
		c.str = b + " OVER " + window
	}
	return c
}

type SelectQuery struct {
	locked bool

	numTokens         int
	name              string
	alias             string
	distinct          bool
	prefixes          []interface{}
	suffixes          []interface{}
	columns           []interface{}
	distincton        []interface{}
	from              []interface{}
	where             []interface{}
	having            []interface{}
	join              []*JoinToken
	with              []interface{}
	window            []interface{}
	union             []interface{}
	unionall          []interface{}
	intersect         []interface{}
	intersectall      []interface{}
	limit             uint64
	offset            uint64
	orderby           []string
	groupby           []string
	placeholderFormat PlaceholderFormat
}

func Select(c ...string) *SelectQuery {
	q := &SelectQuery{
		placeholderFormat: Question,
	}
	if len(c) == 0 {
		q.columns = []interface{}{}
	}
	for _, a := range c {
		q.columns = append(q.columns, a)
	}
	return q
}

func (s *SelectQuery) Name(name string) {
	if s.name == "" && name != "" {
		s.name = name
	}
}

func (s *SelectQuery) GetName() string {
	return s.name
}

func (s *SelectQuery) Columns(c ...interface{}) {
	if len(c) == 0 {
		s.columns = []interface{}{}
	}
	for _, a := range c {
		switch a.(type) {
		case string:
			s.columns = append(s.columns, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.columns = append(s.columns, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Prefix(c ...interface{}) {
	if len(c) == 0 {
		s.prefixes = []interface{}{}
	}
	for _, a := range c {
		switch a.(type) {
		case string:
			s.prefixes = append(s.prefixes, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.prefixes = append(s.prefixes, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Suffix(c ...interface{}) {
	if len(c) == 0 {
		s.suffixes = []interface{}{}
	}
	for _, a := range c {
		switch a.(type) {
		case string:
			s.suffixes = append(s.suffixes, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.suffixes = append(s.suffixes, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Distinct() {
	s.distinct = true
}

func (s *SelectQuery) DistinctOn(c ...interface{}) {
	s.distinct = true
	for _, a := range c {
		switch a.(type) {
		case string:
			s.distincton = append(s.distincton, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.distincton = append(s.distincton, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) From(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.from = append(s.from, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.from = append(s.from, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Where(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.where = append(s.where, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.where = append(s.where, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Having(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.having = append(s.having, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.having = append(s.having, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Joins(c ...*JoinToken) {
	for _, a := range c {
		t, ok := (a.join).(*token)
		if ok {
			t.id = s.numTokens
			s.numTokens++
			a.join = t
		}
		t, ok = (a.on).(*token)
		if ok {
			t.id = s.numTokens
			s.numTokens++
			a.on = t
		}
		s.join = append(s.join, a)
	}
}

func (s *SelectQuery) Union(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.union = append(s.union, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.union = append(s.union, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) UnionAll(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.unionall = append(s.unionall, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.unionall = append(s.unionall, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Intersect(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.intersect = append(s.intersect, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.intersect = append(s.intersect, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) IntersectAll(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.intersectall = append(s.intersectall, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.intersectall = append(s.intersectall, t)
			s.numTokens++
		}
	}
}

// func (s *SelectQuery) Join(j interface{}, on interface{}) {
// 	s.Joins(Join(j, on))
// }

// func (s *SelectQuery) LeftJoin(j interface{}, on interface{}) {
// 	s.Joins(LeftJoin(j, on))
// }
//
// func (s *SelectQuery) RightJoin(j interface{}, on interface{}) {
// 	s.Joins(RightJoin(j, on))
// }

func (s *SelectQuery) Alias(a string) {
	s.alias = a
}

func (s *SelectQuery) GetAlias() string {
	return s.alias
}

func (s *SelectQuery) With(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.with = append(s.with, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.with = append(s.with, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Window(c ...interface{}) {
	for _, a := range c {
		switch a.(type) {
		case string:
			s.window = append(s.window, a)
		case *token:
			t := a.(*token)
			t.id = s.numTokens
			s.window = append(s.window, t)
			s.numTokens++
		}
	}
}

func (s *SelectQuery) Limit(l uint64) {
	s.limit = l
}

func (s *SelectQuery) Offset(o uint64) {
	s.offset = o
}

func (s *SelectQuery) OrderBy(o ...string) {
	s.orderby = append(s.orderby, o...)
}

func (s *SelectQuery) OrderByDesc(o ...string) {
	for i := range o {
		o[i] = "-" + o[i]
	}
	s.OrderBy(o...)
}

func (s *SelectQuery) GroupBy(o ...string) {
	s.groupby = append(s.groupby, o...)
}

func (s *SelectQuery) PlaceholderFormat(f PlaceholderFormat) {
	s.placeholderFormat = f
}

func (s *SelectQuery) SetFromQuery(q *SelectQuery) {
	s.columns = q.columns
	s.limit = q.limit
	s.alias = q.alias
	s.offset = q.offset
	s.orderby = q.orderby
	s.groupby = q.groupby
	s.where = q.where
	s.having = q.having
	s.join = q.join
	s.with = q.with
	s.distinct = q.distinct
	s.distincton = q.distincton
	s.prefixes = q.prefixes
	s.suffixes = q.suffixes
	s.placeholderFormat = q.placeholderFormat
}

// TODO: Copy

func (s *SelectQuery) ToSql(uargs [][]interface{}) (sqlStr string, args []interface{}, err error) {

	sql := &bytes.Buffer{}

	if len(s.prefixes) > 0 {
		for i, e := range s.prefixes {
			if i > 0 {
				_, err = sql.WriteString(" ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
		_, err = sql.WriteString(" ")
	}

	if len(s.with) > 0 {
		_, err = sql.WriteString("WITH ")
		if err != nil {
			return
		}
		for i, e := range s.with {
			if i > 0 && i < len(s.with) {
				_, err = sql.WriteString(", ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				if t.isQuery {
					// normalize t.SetPlaceholderFormat(s.placeholderFormat)
					_, err = sql.WriteString(t.alias)
					if err != nil {
						return
					}
					_, err = sql.WriteString(" AS (")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
					_, err = sql.WriteString(")")
					if err != nil {
						return
					}
				} else {
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
		_, err = sql.WriteString(" ")
		if err != nil {
			return
		}
	}

	_, err = sql.WriteString("SELECT ")
	if err != nil {
		return
	}

	if s.distinct {
		_, err = sql.WriteString("DISTINCT ")
		if err != nil {
			return
		}
		if len(s.distincton) != 0 {
			_, err = sql.WriteString("ON (")
			if err != nil {
				return
			}
			for i, e := range s.distincton {
				if i > 0 && i < len(s.distincton) {
					_, err = sql.WriteString(", ")
					if err != nil {
						return
					}
				}
				switch t := e.(type) {
				case *token:
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
					args = append(args, uargs[t.id]...)
				case string:
					_, err = sql.WriteString(t)
					if err != nil {
						return
					}
				}
			}
			_, err = sql.WriteString(") ")
			if err != nil {
				return
			}
		}
	}

	if len(s.columns) > 0 {
		for i, e := range s.columns {
			if i > 0 && i < len(s.columns) {
				_, err = sql.WriteString(", ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
	} else {
		_, err = sql.WriteString("*")
		if err != nil {
			return
		}
	}

	if len(s.from) > 0 {
		_, err = sql.WriteString(" FROM ")
		if err != nil {
			return
		}

		for i, e := range s.from {
			if i > 0 && i < len(s.from) {
				_, err = sql.WriteString(", ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				if t.isQuery {
					// normalize t.SetPlaceholderFormat(s.placeholderFormat)
					_, err = sql.WriteString("(")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
					_, err = sql.WriteString(") ")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.alias)
					if err != nil {
						return
					}
				} else {
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
	}

	if len(s.join) > 0 {
		_, err = sql.WriteString(" ")
		if err != nil {
			return
		}
		for i, e := range s.join {
			if i > 0 && i < len(s.join) {
				_, err = sql.WriteString(" ")
				if err != nil {
					return
				}
			}
			switch e.jtype {
			case JOIN:
				_, err = sql.WriteString("JOIN ")
				if err != nil {
					return
				}
			case LEFTJOIN:
				_, err = sql.WriteString("LEFT JOIN ")
				if err != nil {
					return
				}
			case RIGHTJOIN:
				_, err = sql.WriteString("RIGHT JOIN ")
				if err != nil {
					return
				}
			}
			switch t := e.join.(type) {
			case *token:
				if t.isQuery {
					// normalize t.SetPlaceholderFormat(s.placeholderFormat)
					_, err = sql.WriteString("(")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
					_, err = sql.WriteString(") ")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.alias)
					if err != nil {
						return
					}
				} else {
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
			_, err = sql.WriteString(" ON ")
			if err != nil {
				return
			}
			switch t := e.on.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
	}

	if len(s.where) > 0 {
		_, err = sql.WriteString(" WHERE ")
		if err != nil {
			return
		}
		for i, e := range s.where {
			if i > 0 {
				_, err = sql.WriteString(" AND ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
	}

	if len(s.groupby) > 0 {
		_, err = sql.WriteString(" GROUP BY ")
		if err != nil {
			return
		}
		_, err = sql.WriteString(strings.Join(s.groupby, ", "))
		if err != nil {
			return
		}
	}

	if len(s.having) > 0 {
		_, err = sql.WriteString(" HAVING ")
		if err != nil {
			return
		}
		for i, e := range s.having {
			if i > 0 && i < len(s.having) {
				_, err = sql.WriteString(" AND ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
	}

	if len(s.window) > 0 {
		_, err = sql.WriteString(" WINDOW ")
		if err != nil {
			return
		}
		for i, e := range s.window {
			if i > 0 && i < len(s.window) {
				_, err = sql.WriteString(", ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				if t.isQuery {
					// normalize t.SetPlaceholderFormat(s.placeholderFormat)
					_, err = sql.WriteString(t.alias)
					if err != nil {
						return
					}
					_, err = sql.WriteString(" AS (")
					if err != nil {
						return
					}
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
					_, err = sql.WriteString(")")
					if err != nil {
						return
					}
				} else {
					_, err = sql.WriteString(t.str)
					if err != nil {
						return
					}
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
		_, err = sql.WriteString(" ")
		if err != nil {
			return
		}
	}

	if len(s.orderby) > 0 {
		_, err = sql.WriteString(" ORDER BY ")
		if err != nil {
			return
		}
		_, err = sql.WriteString(strings.Join(s.orderby, ", "))
		if err != nil {
			return
		}
	}

	if s.limit > 0 {
		_, err = sql.WriteString(" LIMIT ")
		if err != nil {
			return
		}
		_, err = sql.WriteString(strconv.FormatUint(s.limit, 10))
		if err != nil {
			return
		}
	}

	if s.offset > 0 {
		_, err = sql.WriteString(" OFFSET ")
		if err != nil {
			return
		}
		_, err = sql.WriteString(strconv.FormatUint(s.offset, 10))
		if err != nil {
			return
		}
	}

	if len(s.suffixes) > 0 {
		for i, e := range s.suffixes {
			if i > 0 {
				_, err = sql.WriteString(" ")
				if err != nil {
					return
				}
			}
			switch t := e.(type) {
			case *token:
				_, err = sql.WriteString(t.str)
				if err != nil {
					return
				}
				args = append(args, uargs[t.id]...)
			case string:
				_, err = sql.WriteString(t)
				if err != nil {
					return
				}
			}
		}
		_, err = sql.WriteString(" ")
		if err != nil {
			return
		}
	}

	sqlStr = s.placeholderFormat.ReplacePlaceholders(sql.String())

	return
}
