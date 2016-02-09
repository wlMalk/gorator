package parser

func (c *Config) def() {
	c.defImports()
}

func (c *Config) defImports() {
	c.Imports = map[string][]map[string]string{}

	c.Imports["database"] = []map[string]string{
		map[string]string{
			"database/sql": "",
			"strings":      "",
			"fmt":          "",
		}, map[string]string{},
		map[string]string{
			"github.com/wlMalk/gorator/database": "odatabase",
		},
	}
	c.Imports["orm"] = []map[string]string{
		map[string]string{},
		map[string]string{
			c.Path + "/database":           "",
			c.Path + "/database/orm/query": "",
			c.Path + "/database/orm/model": "",
		}, map[string]string{
			"github.com/wlMalk/gorator/database": "odatabase",
		},
	}
	c.Imports["query"] = []map[string]string{
		map[string]string{
			"database/sql": "",
			"strings":      "",
			"fmt":          "",
		}, map[string]string{},
		map[string]string{
			"github.com/wlMalk/gorator/query":    "oquery",
			"github.com/wlMalk/gorator/database": "odatabase",
		},
	}
	c.Imports["model"] = []map[string]string{
		map[string]string{
			"encoding/json": "",
			"bytes":         "",
		}, map[string]string{
			c.Path + "/database/orm/query": "",
		},
		map[string]string{},
	}
	c.Imports["callback"] = []map[string]string{
		map[string]string{},
		map[string]string{},
		map[string]string{},
	}
}

func (d *Database) def() {
	d.DriverName = defaultDriver
}

func (mo *Model) def() {
	mo.SoftDelete = true

	mo.CreatedAt = true
	mo.UpdatedAt = true
	mo.DeletedAt = true

	mo.CreatedBy = true
	mo.UpdatedBy = true
	mo.DeletedBy = true

	mo.Uuid = 4

	mo.Sliced = true

	mo.defPrimaryKey()

}

func (t *Table) def() {
	t.Schema = ""
	t.IsPivot = false
}

func (r *Relation) def() {

}

func (f *Field) def() {
	f.Null = true
	f.InDB = true
	f.InEncoding = true

	f.defCallbacks()
}

func (f *Field) defCallbacks() {

}

func (mo *Model) defPrimaryKey() {
	mo.PrimaryKey = &PrimaryKey{
		Model:  mo,
		Fields: []string{"ID"},
	}
}
