package generate

import (
	"errors"
	"os"
	"io/ioutil"

	"github.com/wlMalk/ormator/generate/internal/parser"
)

const (
	VERSION = "0.1"
)

func GenerateFromFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("could not open config file")
	}
	return Generate(getPath(path), file)
}

func Generate(path string, b []byte) error {
	config, err := parser.Parse(b)
	if err != nil {
		return err
	}
	return generateORM(path, config)
}

func generateORM(path string, config *parser.Config) error {
	mkdir(path + "database")
	mkdir(path + "database" + string(os.PathSeparator) + "orm")
	mkdir(path + "database" + string(os.PathSeparator) + "orm" + string(os.PathSeparator) + "callback")
	mkdir(path + "database" + string(os.PathSeparator) + "orm" + string(os.PathSeparator) + "table")
	mkdir(path + "database" + string(os.PathSeparator) + "orm" + string(os.PathSeparator) + "query")
	mkdir(path + "database" + string(os.PathSeparator) + "orm" + string(os.PathSeparator) + "model")
	mkdir(path + "database" + string(os.PathSeparator) + "orm" + string(os.PathSeparator) + "model" + string(os.PathSeparator) + "slice")
	return nil
}
