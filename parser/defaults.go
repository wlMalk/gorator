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
}

func (r *Relation) def() {

}

func (f *Field) def() {
	f.InDB = true
	f.defCallbacks()
}

func (f *Field) defCallbacks() {

}

func (m *Model) defPrimaryKey() {
	m.PrimaryKey = &PrimaryKey{
		Fields: []string{"ID"},
	}
}
