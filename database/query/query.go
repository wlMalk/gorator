package query

type Query interface {
	Name(string)
	GetName() string
	ToSql() (string, []interface{}, error)
}

func Tokenize(q Query) (t *Token, err error) {
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
			t, targs, err := q.ToSql()
			if err != nil {
				return nil, nil, err
			}
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

type Selecter interface {
	Query
	Select(...interface{})
	GetAlias() string
	SetPlaceholderFormat(PlaceholderFormat)
}

type token struct {
	id  int
	str string
}

type Token struct {
	str  string
	args []interface{}
}

func (t *Token) String() string {
	return t.str
}

func (t *Token) Token() (*token, []interface{}) {
	return &token{str: t.str}, t.args
}

func T(str string, args ...interface{}) *Token {
	return &Token{str, args}
}

type Preload struct {
	Name     string
	Q        *SelectQuery
	Preloads []*Preload
}

func P(name string, q *SelectQuery, preloads []*Preload) *Preload {
	return &Preload{name, q, preloads}
}
