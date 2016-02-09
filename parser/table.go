package parser

import (
	"fmt"
	"strings"

	"github.com/wlMalk/gorator/internal/util"
)

const (
	tableNameInDB = "nameInDB"
	tableSchema   = "schema"
)

func (t *Table) parse(name string, m map[interface{}]interface{}) error {
	t.def()
	t.Name = strings.ToLower(util.Plural(name))

	err := t.parseNameInDB(m)
	if err != nil {
		return err
	}

	err = t.parseSchema(m)
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) parseNameInDB(m map[interface{}]interface{}) error {
	avi, ok := m[tableNameInDB]
	if ok {
		av, ok := avi.(string)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' table", tableNameInDB, t.Name)
		}
		t.Name = av
	}
	return nil
}

func (t *Table) parseSchema(m map[interface{}]interface{}) error {
	avi, ok := m[tableSchema]
	if ok {
		av, ok := avi.(string)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' table", tableSchema, t.Name)
		}
		t.Schema = av
	}
	return nil
}
