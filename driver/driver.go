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
}

type BaseDriver struct {
	name      string
	files     []string
	types     map[string]string
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

func (d *BaseDriver) Parse() {
	template.Must(d.templates.ParseFiles(d.files...))
}

func (d *BaseDriver) Execute(w io.Writer, t string, data interface{}) error {
	return d.templates.ExecuteTemplate(w, t, data)
}

var drivers = map[string]Driver{}

func New(name string, types map[string]string, files []string) BaseDriver {
	d := BaseDriver{name, files, types, template.New("")}
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
