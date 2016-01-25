package generate

import (
	"io/ioutil"
	"os"
	"strings"
)

func getPath(path string) string {
	return strings.Replace(path, "config.yml", "", 1)
}

func getImportPath(path string) string {
	return strings.Replace(strings.Replace(path, os.Getenv("GOPATH")+string(os.PathSeparator)+"src"+string(os.PathSeparator), "", 1), string(os.PathSeparator)+"config.yml", "", 1)
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

