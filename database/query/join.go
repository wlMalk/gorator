package query

const (
	JOIN = iota
	LEFTJOIN
	RIGHTJOIN
)

type JoinToken struct {
	jtype int
	join  interface{}
	on    interface{}
}

func join(jtype int, j interface{}, on interface{}) *JoinToken {
	jt := &JoinToken{}
	if j == nil || on == nil {
		return jt
	}
	switch j.(type) {
	case string, *Token, Query:
		jt.join = j
	}
	switch on.(type) {
	case string, *Token:
		jt.on = on
	}
	jt.jtype = jtype
	return jt
}

func JT(jt *JoinToken) (j *JoinToken, a []interface{}, b []interface{}, err error) {
	if jt == nil {
		return nil, nil, nil, nil
	}
	switch (jt.join).(type) {
	case *Token:
		at := jt.join.(*Token)
		jt.join, a = at.Token()
	case Query:
		q := jt.join.(Query)
		at, err := Tokenize(q)
		if err != nil {
			return nil, nil, nil, err
		}
		jt.join, a = at.Token()
	}
	switch (jt.on).(type) {
	case *Token:
		at := jt.join.(*Token)
		jt.on, b = at.Token()
	}
	j = jt
	return
}

func Join(j interface{}, on interface{}) *JoinToken {
	return join(JOIN, j, on)
}

func LeftJoin(j interface{}, on interface{}) *JoinToken {
	return join(LEFTJOIN, j, on)
}

func RightJoin(j interface{}, on interface{}) *JoinToken {
	return join(RIGHTJOIN, j, on)
}
