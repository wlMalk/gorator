package parser

func (c *Config) GetSchemas() (s []*Schema) {
	for _, d := range c.Databases {
		s = append(s, d.GetSchemas()...)
	}
	return
}

func (d *Database) GetSchemas() (s []*Schema) {
	for _, m := range d.Models {
		if m.Table.Schema != "" {
			ind := -1
			for i, sc := range s {
				if sc.Name == m.Table.Schema {
					ind = i
					break
				}
			}
			if ind == -1 {
				s = append(s, &Schema{
					Database: d,
					Name:     m.Table.Schema,
					Tables:   []*Table{m.Table},
				})
			} else {
				s[ind].Tables = append(s[ind].Tables, m.Table)
			}
		}
	}
	return
}

func (c *Config) GetModel(name string) *Model {
	for _, m := range c.GetAllModels() { // reconsider this
		if m.Name == name {
			return m
		}
	}
	return nil
}

func (c *Config) GetAllModels() (s []*Model) {
	for _, d := range c.Databases {
		s = append(s, d.Models...)
	}
	return
}

func (d *Database) GetAllModels() (s []*Model) {
	return d.Models
}

func (c *Config) GetModels() (s []*Model) {
	for _, d := range c.Databases {
		s = append(s, d.GetModels()...)
	}
	return
}

func (d *Database) GetModels() (s []*Model) {
	for _, m := range d.Models {
		if !m.IsPivot {
			s = append(s, m)
		}
	}
	return
}

func (c *Config) GetPivotModels() (s []*Model) {
	for _, d := range c.Databases {
		s = append(s, d.GetPivotModels()...)
	}
	return
}

func (d *Database) GetPivotModels() (s []*Model) {
	for _, m := range d.Models {
		if m.IsPivot {
			s = append(s, m)
		}
	}
	return
}

func (c *Config) GetAllTables() (s []*Table) {
	for _, d := range c.Databases {
		s = append(s, d.GetAllTables()...)
	}
	return
}

func (d *Database) GetAllTables() (s []*Table) {
	for _, m := range d.Models {
		s = append(s, m.Table)
	}
	return
}

func (c *Config) GetTables() (s []*Table) {
	for _, d := range c.Databases {
		s = append(s, d.GetTables()...)
	}
	return
}

func (d *Database) GetTables() (s []*Table) {
	for _, m := range d.Models {
		if !m.IsPivot {
			s = append(s, m.Table)
		}
	}
	return
}

func (c *Config) GetPivotTables() (s []*Table) {
	for _, d := range c.Databases {
		s = append(s, d.GetPivotTables()...)
	}
	return
}

func (d *Database) GetPivotTables() (s []*Table) {
	for _, m := range d.Models {
		if m.IsPivot {
			s = append(s, m.Table)
		}
	}
	return
}

func (c *Config) GetPrimaryKeys() (s []*PrimaryKey) {
	for _, d := range c.Databases {
		s = append(s, d.GetPrimaryKeys()...)
	}
	return
}

func (d *Database) GetPrimaryKeys() (s []*PrimaryKey) {
	for _, m := range d.Models {
		s = append(s, m.PrimaryKey)
	}
	return
}

func (c *Config) GetDrivers() (s []string) {
	for _, d := range c.Databases {
		s = append(s, d.DriverName)
	}
	return
}

func (c *Config) GetPackage(name string) (p *Package) {
	for _, pa := range c.Packages {
		if pa.Name == name {
			p = pa
		}
	}
	return
}
