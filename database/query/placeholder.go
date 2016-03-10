package query

import (
	"bytes"
	"fmt"
	"strings"
)

type PlaceholderFormat interface {
	ReplacePlaceholders(sql string) string
}

var (
	Question = questionFormat{}
	Dollar   = dollarFormat{}
)

type questionFormat struct{}

func (_ questionFormat) ReplacePlaceholders(sql string) string {
	return sql
}

type dollarFormat struct{}

func (_ dollarFormat) ReplacePlaceholders(sql string) string {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" {
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			fmt.Fprintf(buf, "$%d", i)
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String()
}

func Placeholders(count int) string {
	if count < 1 {
		return ""
	}

	return strings.Repeat(",?", count)[1:]
}
