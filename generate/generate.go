package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/wlMalk/gorator/driver/pgsql"

	"github.com/wlMalk/gorator/parser"
)

const (
	VERSION = "0.1"
)

var tmpls *template.Template = template.New("")

var ormTmplsMap map[string]string = map[string]string{
	"database": "database",
	"orm":      "database/orm",
	"callback": "database/orm/callback",
	"query":    "database/orm/query",
	"model":    "database/orm/model",
}
var ormDirs []string = []string{
	"database", "database/orm", "database/orm/callback",
	"database/orm/query", "database/orm/model",
}

func init() {
	tmplsDir := getFullPath(os.Getenv("GOPATH"), "src/github.com/wlMalk/gorator/templates/")
	tmplFiles, err := filepath.Glob(tmplsDir + "*/*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	tmpls, err = tmpls.Funcs(getFuncsMap()).ParseFiles(tmplFiles...)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func Generate(path string) error {
	path = getPath(path)

	configFiles, err := filepath.Glob(path + "/*.yml")
	if err != nil {
		return fmt.Errorf("could not open config file")
	}

	importPath := getImportPath(path)
	var files [][]byte

	for _, f := range configFiles {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("could not open config file")
		}
		files = append(files, file)
	}

	config, err := parser.Parse(importPath, files...)
	if err != nil {
		return err
	}

	// a, _ := json.Marshal(config)
	// fmt.Println(string(a))
	return generateORM(path, config)
}

func generateORM(path string, config *parser.Config) error {
	for _, d := range ormDirs {
		mkdir(getFullPath(path, d))
	}

	var w bytes.Buffer
	for t, dir := range ormTmplsMap {
		for _, db := range config.Databases {
			err := tmpls.ExecuteTemplate(&w, "heading", getPackage(t, config))
			if err != nil {
				return err
			}

			if db.Driver.Generate(t) {
				err = tmpls.ExecuteTemplate(&w, t, config)
				if err != nil {
					return err
				}
			}
			// driver headings
			err = db.Driver.Execute(&w, t, config)
			if err != nil {
				return err
			}
		}

		// fix white spaces

		b, err := format.Source(w.Bytes())
		if err != nil {
			return err
		}

		err = saveFile(getFullPath(path, dir+"/"+t+"_gen.go"), b)
		if err != nil {
			return err
		}

		w.Reset()
	}
	return nil
}
