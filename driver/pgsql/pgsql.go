package pgsql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wlMalk/gorator/driver"
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

	generate := map[string]bool{
		"database": true,
		"orm":      true,
		"query":    true,
		"model":    true,
		"callback": true,
	}

	tmplsDir := os.Getenv("GOPATH") + "/src/github.com/wlMalk/gorator/driver/pgsql/templates/"
	files, err := filepath.Glob(tmplsDir + "*/*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	pgsql = &Driver{driver.New("pgsql", types, generate, files)}

	driver.Register(pgsql)
}
