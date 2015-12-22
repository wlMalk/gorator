package generate

import (
	"errors"
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
	return generate(path, config)
}

func generate(path string, config *parser.Config) error {
	return nil
}
