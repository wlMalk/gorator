package driver

import (
	"fmt"
)

type Driver interface {
	Name() string
	Type(string) string
	Types() map[string]string
	Path() string
	Appended() bool
	SetFuncs() // sets templates funcs
	ExecuteDatabase(interface{}) ([]byte, error)
	ExecuteQuery(interface{}) ([]byte, error)
	ExecuteModel(interface{}) ([]byte, error)
	ExecuteTable(interface{}) ([]byte, error)
	ExecuteSlice(interface{}) ([]byte, error)
	ExecuteCallbacks(interface{}) ([]byte, error)
}

var drivers map[string]Driver = map[string]Driver{}

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
