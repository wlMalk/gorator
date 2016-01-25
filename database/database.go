package database

type Table struct {
	name         string
	schema       string
	columns      map[string]string
	nameReady    bool
	schemaReady  bool
	columnsReady bool
}

func (t *Table) SetName(name string) {
	if !t.nameReady {
		t.name = name
		t.nameReady = true
	}
}

func (t *Table) SetSchema(schema string) {
	if !t.schemaReady {
		t.schema = schema
		t.schemaReady = true
	}
}

func (t *Table) SetColumns(columns map[string]string) {
	if !t.columnsReady {
		t.columns = columns
		t.columnsReady = true
	}
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Column(a string) string {
	return t.columns[a]
}
