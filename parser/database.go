package parser

import (
	"fmt"

	"github.com/wlMalk/gorator/driver"
)

const (
	databaseDriver = "driver"
	databaseModels = "models"
)

func (d *Database) parse(name string, m map[interface{}]interface{}) error {
	d.def()
	d.Name = name

	err := d.parseDriver(m)
	if err != nil {
		return err
	}

	err = d.parseModels(m)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) parseDriver(m map[interface{}]interface{}) error {
	dr, err := driver.Get(d.DriverName)
	if err != nil {
		return err
	}

	d.Driver = dr

	return nil
}

func (d *Database) parseModels(m map[interface{}]interface{}) error {
	if _, ok := m[databaseModels]; ok {
		for k, v := range mi(m[databaseModels]) {
			table := &Table{}
			err := table.parse(s(k), mi(v))
			if err != nil {
				return err
			}
			model := &Model{}
			model.Table = table
			table.Model = model
			model.Database = d
			err = model.parse(s(k), mi(v))
			if err != nil {
				return err
			}
			d.Models = append(d.Models, model)
		}
	} else {
		return fmt.Errorf("no '%s' found in '%s' database", databaseModels, d.Name)
	}
	return nil
}
