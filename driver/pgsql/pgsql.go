package pgsql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wlMalk/gorator/driver"
	"github.com/wlMalk/gorator/internal/util"
)

type Driver struct {
	driver.BaseDriver
}

func init() {
	var pgsql *Driver

	types := map[string]string{
		"bigint":    "int64",
		"int":       "int",
		"integer":   "int",
		"smallint":  "int",
		"character": "string",
		"text":      "string",
		"timestamp": "time.Time",
		"numeric":   "float64",
		"boolean":   "bool",
		"serial":    "int",
		"bigserial": "int64",
		"varchar":   "string",
	}

	tmplsDir := os.Getenv("GOPATH") + "/src/github.com/wlMalk/gorator/templates/"
	files, err := filepath.Glob(tmplsDir + "pgsql.*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	commons, err := filepath.Glob(tmplsDir + "common.*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	files = append(files, commons...)

	pgsql = &Driver{driver.New("pgsql", types, files)}
	pgsql.SetFuncs(util.GetFuncsMap())
	pgsql.Parse()

	driver.Register(pgsql)
}
