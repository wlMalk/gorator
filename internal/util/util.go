package util

import (
	"go/ast"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/gedex/inflector"
	"github.com/serenize/snaker"
)

func GetFuncsMap() template.FuncMap {
	return template.FuncMap{
		"snakecase": snaker.CamelToSnake,
		"plural":    inflector.Pluralize,
		"unexport":  strings.ToLower,
		"exported":  ast.IsExported,
		"primitive": Primitive,
		"numerical": Numerical,
		"lower":     strings.ToLower,
		"upper":     strings.ToTitle,
		"zeroValue": ZeroValue,
	}
}

func Plural(s string) string {
	return inflector.Pluralize(s)
}

func Snakecase(s string) string {
	return snaker.CamelToSnake(s)
}

func GetPath(path string) string {
	return strings.TrimSuffix(path, string(os.PathSeparator))
}

func GetImportPath(path string) string {
	return GetPath(strings.Replace(path, os.Getenv("GOPATH")+string(os.PathSeparator)+"src"+string(os.PathSeparator), "", 1))
}

func GetFullPath(path string, subpath string) string {
	return path + string(os.PathSeparator) + strings.Replace(subpath, "/", string(os.PathSeparator), -1)
}

func Mkdir(name string) error {
	return os.MkdirAll(name, 0777)
}

func SaveFile(p string, data []byte) error {
	return ioutil.WriteFile(p, data, 0666)
}

func Primitive(a string) bool {
	if a == "uint" || a == "uint8" || a == "uint16" || a == "uint32" || a == "uint64" ||
		a == "int" || a == "int8" || a == "int16" || a == "int32" || a == "int64" ||
		a == "float32" || a == "float64" ||
		a == "complex32" || a == "complex64" ||
		a == "byte" || a == "[]byte" ||
		a == "time.Time" ||
		a == "rune" ||
		a == "bool" ||
		a == "string" {
		return true
	}
	return false
}

func ZeroValue(a string) string {
	if a == "uint" || a == "uint8" || a == "uint16" || a == "uint32" || a == "uint64" ||
		a == "int" || a == "int8" || a == "int16" || a == "int32" || a == "int64" ||
		a == "float32" || a == "float64" ||
		a == "complex32" || a == "complex64" {
		return "0"
	}
	if a == "byte" || a == "[]byte" || a == "*time.Time" {
		return "nil"
	}
	if a == "rune" {
		return "''"
	}
	if a == "bool" {
		return "f"
	}
	if a == "string" {
		return "\"\""
	}
	return ""
}

func Numerical(a string) bool {
	if a == "uint" || a == "uint8" || a == "uint16" || a == "uint32" || a == "uint64" ||
		a == "int" || a == "int8" || a == "int16" || a == "int32" || a == "int64" ||
		a == "time.Time" ||
		a == "float32" || a == "float64" ||
		a == "complex32" || a == "complex64" {
		return true
	}
	return false
}
