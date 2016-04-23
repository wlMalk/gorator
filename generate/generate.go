package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/wlMalk/gorator/driver/pgsql"
	"github.com/wlMalk/gorator/internal/util"

	"github.com/wlMalk/gorator/parser"

	"github.com/fatih/color"
)

var magenta = color.New(color.FgMagenta, color.Underline).SprintFunc()
var okMsg = color.New(color.FgGreen, color.Bold).SprintFunc()("OK!")

const (
	VERSION = "0.1"
)

var tmpls *template.Template = template.New("")

var tmplsDir = ""

func init() {
	tmplsDir = util.GetFullPath(os.Getenv("GOPATH"), "src/github.com/wlMalk/gorator/templates/")
	tmplFiles, err := filepath.Glob(tmplsDir + "*/*.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	tmpls, err = tmpls.Funcs(util.GetFuncsMap()).ParseFiles(tmplFiles...)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func Generate(path string, version string, cmd string) error {
	path = util.GetPath(path)

	if version != "" {
		version = "." + version
	}

	configFiles, err := filepath.Glob(path + "/*.yml" + version)
	if err != nil {
		return fmt.Errorf("could not open config file")
	}

	if len(configFiles) == 0 {
		return fmt.Errorf("could not find any config file")
	}

	importPath := util.GetImportPath(path)
	var files [][]byte

	for _, f := range configFiles {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("could not open config file")
		}
		files = append(files, file)
	}

	fmt.Printf("parsing config file(s)...") // if show
	config, err := parser.Parse(importPath, files...)
	if err != nil {
		return err
	}
	fmt.Printf("\t%s\n", okMsg)

	return GenerateFrom(path, config, cmd)
}

func GenerateFrom(path string, config *parser.Config, cmd string) error {
	fmt.Printf("generating %s...\n", magenta(cmd)) // if show
	err := generateORM(path, config)
	if err != nil {
		return err
	}
	goFiles, err := filepath.Glob(tmplsDir + "*.go")
	if err != nil {
		return err
	}
	if len(goFiles) > 0 {
		for _, f := range goFiles {
			name := filepath.Base(f)
			if strings.HasSuffix(name, ".inc.go") || strings.HasSuffix(name, ".rep.go") {
				p := config.GetPackage(name[:len(name)-7])
				if p == nil {
					continue
				}
				nPath := p.Path + "/" + name
				if strings.HasSuffix(name, ".rep.go") {
					nPath = p.Path + "/" + name[:len(name)-7] + ".go"
				}
				err = util.MoveFile(util.GetFullPath(path, name), util.GetFullPath(path, nPath))
				if err != nil {
					return err
				}
			}
		}
	}
	fmt.Printf("happy coding!\n") // if show
	return nil
}

func generateORM(path string, config *parser.Config) error {
	for _, p := range config.Packages {
		fmt.Printf("making %s folder...", magenta(p.Path)) // if show
		err := util.Mkdir(util.GetFullPath(path, p.Path))
		if err != nil {
			return err
		}
		fmt.Printf("\t%s\n", okMsg)

		var w bytes.Buffer
		err = tmpls.ExecuteTemplate(&w, "heading", p)
		if err != nil {
			return err
		}
		for _, db := range config.Databases {
			if p.Name != "database" {
				for _, model := range db.Models {
					if db.Driver.Generate(p.Name) {
						fmt.Printf("generating code for %s package - %s model...", magenta(p.Name), magenta(model.Name)) // if show
						err = tmpls.ExecuteTemplate(&w, p.Name, model)
						if err != nil {
							return err
						}
						fmt.Printf("\t%s\n", okMsg)
					} else {
						// driver headings
						fmt.Printf("generating %s driver code for %s package - %s model...", magenta(db.Driver.Name()), magenta(p.Name), magenta(model.Name)) // if show
						err = db.Driver.Execute(&w, p.Name, model)
						if err != nil {
							return err
						}
						fmt.Printf("\t%s\n", okMsg)
					}
				}
			} else {
				if db.Driver.Generate(p.Name) {
					fmt.Printf("generating code for %s database...", magenta(db.Name)) // if show
					err = tmpls.ExecuteTemplate(&w, p.Name, db)
					if err != nil {
						return err
					}
					fmt.Printf("\t%s\n", okMsg)
				}
				// driver headings
				// check if tmp is provided
				fmt.Printf("generating %s driver code for %s database...", magenta(db.Driver.Name()), magenta(db.Name)) // if show
				err = db.Driver.Execute(&w, p.Name, db)
				if err != nil {
					return err
				}
				fmt.Printf("\t%s\n", okMsg)
			}
		}

		// fix white spaces

		fmt.Printf("formatting code for %s package...", magenta(p.Name)) // if show
		b, err := format.Source(w.Bytes())
		if err != nil {
			return err
		}
		fmt.Printf("\t%s\n", okMsg)

		fmt.Printf("saving %s package file...", magenta(p.Name)) // if show
		err = util.SaveFile(util.GetFullPath(path, p.Path+"/"+p.Name+"_gen.go"), b)
		if err != nil {
			return err
		}
		fmt.Printf("\t%s\n", okMsg)

		w.Reset()
	}
	return nil
}
