package parser

func (c *Config) def() {
	c.defPackages()
}

func (c *Config) defPackages() {
	c.Packages = map[string]*Package{}
	c.addPackage(packageDatabase(c))
	c.addPackage(packageORM(c))
	c.addPackage(packageQuery(c))
	c.addPackage(packageModel(c))
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

	mo.Listed = true

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
	f.Exported = true
	f.InEncoding = true

	f.defCallbacks()
}

func (f *Field) defCallbacks() {

}

func (mo *Model) defPrimaryKey() {
	mo.PrimaryKey = &PrimaryKey{
		Model:  mo,
		Fields: []string{"Id"},
	}
}
