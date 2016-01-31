package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/wlMalk/ormator/driver"
	_ "github.com/wlMalk/ormator/driver/pgsql"

	"github.com/wlMalk/ormator/parser"
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
	"table":    "database/orm/table",
	"model":    "database/orm/model",
	"slice":    "database/orm/slice",
}
var ormDirs []string = []string{
	"database", "database/orm", "database/orm/callback",
	"database/orm/query", "database/orm/table",
	"database/orm/model", "database/orm/slice",
}

func init() {
	tmplsDir := getFullPath(os.Getenv("GOPATH"), "src/github.com/wlMalk/ormator/templates/")
	tmplFiles, err := filepath.Glob(tmplsDir + "*/*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	tf, err := filepath.Glob(tmplsDir + "*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	tmplFiles = append(tmplFiles, tf...)
	tmpls, err = tmpls.Funcs(getFuncsMap()).ParseFiles(tmplFiles...)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func GenerateFromFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not open config file")
	}
	return Generate(getPath(path), file)
}

func Generate(path string, b []byte) error {
	config, err := parser.Parse(b)
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
		err := tmpls.ExecuteTemplate(&w, t, config)
		if err != nil {
			return err
		}
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
