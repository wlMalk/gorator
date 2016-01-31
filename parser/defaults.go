package parser

func (c *Config) def() {

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
	f.InDB = true
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
