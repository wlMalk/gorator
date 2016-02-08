package generate

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/wlMalk/gorator/parser"

)

func getPackage(name string, config *parser.Config) interface{} {

	return struct {
		Name           string
		GoratorVersion string
		ConfigVersion  string
		Description    string
		Imports        []map[string]string
	}{
		Name:           name,
		GoratorVersion: config.GoratorVersion,
		ConfigVersion:  config.Version,
		Description:    packageDescriptions[name],
		Imports:        config.Imports[name],
	}

}

func getPath(path string) string {
	return strings.TrimSuffix(path, string(os.PathSeparator))
}

func getImportPath(path string) string {
	return getPath(strings.Replace(path, os.Getenv("GOPATH")+string(os.PathSeparator)+"src"+string(os.PathSeparator), "", 1))
}

func getFullPath(path string, subpath string) string {
	return path + string(os.PathSeparator) + strings.Replace(subpath, "/", string(os.PathSeparator), -1)
}

func mkdir(name string) error {
	return os.MkdirAll(name, 0777)
}

func saveFile(p string, data []byte) error {
	return ioutil.WriteFile(p, data, 0666)
}

