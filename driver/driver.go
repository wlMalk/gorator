package driver

import (
	"fmt"
	"io"
	"text/template"
)

type Driver interface {
	Name() string
	Type(string) string
	Types() map[string]string
	SetFuncs(template.FuncMap) // sets templates funcs
	Execute(io.Writer, string, interface{}) error
	ExecuteDatabase(io.Writer, interface{}) error
	ExecuteQuery(io.Writer, interface{}) error
	ExecuteModel(io.Writer, interface{}) error
	ExecuteCallbacks(io.Writer, interface{}) error
	Generate(string) bool
	GenerateQuery() bool
	GenerateDatabase() bool
	GenerateModel() bool
	GenerateCallbacks() bool
}

type BaseDriver struct {
	name      string
	files     []string
	types     map[string]string
	generate  map[string]bool
	templates *template.Template
}

func (d *BaseDriver) Name() string {
	return d.name
}

func (d *BaseDriver) Type(t string) string {
	return d.types[t]
}

func (d *BaseDriver) Types() map[string]string {
	return d.types
}

func (d *BaseDriver) Appended() bool {
	return true
}

func (d *BaseDriver) SetFuncs(f template.FuncMap) { // sets templates funcs
	d.templates = d.templates.Funcs(f)
}

func (d *BaseDriver) Parse() { // sets templates funcs
	d.templates.ParseFiles(d.files...)
}

func (d *BaseDriver) ExecuteDatabase(w io.Writer, data interface{}) error {
	return d.templates.ExecuteTemplate(w, "database", data)
}

func (d *BaseDriver) ExecuteModel(w io.Writer, data interface{}) error {
	return d.templates.ExecuteTemplate(w, "model", data)
}

func (d *BaseDriver) ExecuteQuery(w io.Writer, data interface{}) error {
	return d.templates.ExecuteTemplate(w, "query", data)
}

func (d *BaseDriver) ExecuteCallbacks(w io.Writer, data interface{}) error {
	return d.templates.ExecuteTemplate(w, "callback", data)
}

func (d *BaseDriver) Execute(w io.Writer, t string, data interface{}) error {
	return d.templates.ExecuteTemplate(w, t, data)
}

func (d *BaseDriver) GenerateDatabase() bool {
	return d.generate["database"]
}

func (d *BaseDriver) GenerateModel() bool {
	return d.generate["model"]
}

func (d *BaseDriver) GenerateQuery() bool {
	return d.generate["query"]
}

func (d *BaseDriver) GenerateCallbacks() bool {
	return d.generate["callback"]
}

func (d *BaseDriver) Generate(t string) bool {
	return d.generate[t]
}

var drivers map[string]Driver = map[string]Driver{}

func New(name string, types map[string]string, generate map[string]bool, files []string) BaseDriver {
	d := BaseDriver{name, files, types, generate, template.New("")}
	d.templates.ParseFiles(d.files...)
	return d
}

func Register(d Driver) {
	if d == nil {
		return
	}
	drivers[d.Name()] = d
}

func Get(name string) (Driver, error) {
	d, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("could not find '%s' driver", name)
	}
	return d, nil
}
