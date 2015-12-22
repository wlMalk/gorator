package parser

import (
	"errors"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Path           string
	Version        string
	OrmatorVersion string

	Databases []*Database
}

type Database struct {
	Name    string
	Driver  string
	Schemas []*Schema
}

type Schema struct {
	Name   string
	Models []*Model
}

type Table struct {
	Name   string
	Schema *Schema
	Model  *Model
}

type SuperModel struct {
	Name       string
	Table      *Table
	Fields     []*Field
	Relations  []*Relation
	SoftDelete bool

	CreatedAt  bool
	ModifiedAt bool
	DeletedAt  bool
	CreatedBy  bool
	ModifiedBy bool
	DeletedBy  bool

	UUID int
}

type Model struct {
	SuperModel *SuperModel
	Name       string
	Table      *Table
	Fields     []*Field
	Relations  []*Relation
	SoftDelete bool

	PrimaryKey *PrimaryKey

	CreatedAt  bool
	ModifiedAt bool
	DeletedAt  bool
	CreatedBy  bool
	ModifiedBy bool
	DeletedBy  bool

	Callbacks []string

	UUID        int
	IsSliceable bool
	Slice       *Slice
}

type PivotModel struct {
	Name       string
	Table      *Table
	Fields     []*Field
	Relations  []*Relation
	SoftDelete bool

	PrimaryKey *PrimaryKey

	CreatedAt  bool
	ModifiedAt bool
	DeletedAt  bool
	CreatedBy  bool
	ModifiedBy bool
	DeletedBy  bool

	UUID        int
	IsSliceable bool
	Slice       *Slice
}

type PrimaryKey struct {
	Fields []*Field
}

type Field struct {
	Model  *Model
	Name   string
	NameDB string
	Type   string
	TypeDB string

	Default interface{}

	Validations []string

	Callbacks []string

	IsArray        bool
	IsNull         bool
	IsUnique       bool
	IsNumeric      bool
	IsIncrementing bool
	IsInDB         bool

	Where   bool
	OrderBy bool
	GroupBy bool
	Having  bool
}

type Slice struct {
	Model *Model
	Linq  bool
}

type Relation struct {
	Model          *Model
	Name           string
	otherModelName string
	OtherModel     *Model
	Type           string
	ForeignKey     string
	LocalKey       string
	OtherKey       string
	PivotTable     *PivotModel
}

func defaultDatabase() *Database {
	return &Database{
		Driver: "pgsql",
	}
}

func defaultSchema() *Schema {
	return &Schema{}
}

func defaultModel() *Model {
	return &Model{
		SoftDelete: true,

		CreatedAt:  true,
		ModifiedAt: true,
		DeletedAt:  true,

		CreatedBy:  true,
		ModifiedBy: true,
		DeletedBy:  true,

		UUID: 4,

		IsSliceable: true,
	}
}

func defaultField() *Field {
	return &Field{}
}

func defaultFieldCallbacks() []string {
	return nil
}

func defaultPrimaryKeyField() *Field {
	return &Field{
		Name: "ID",
	}
}

func defaultPrimaryKey() *PrimaryKey {
	return &PrimaryKey{
		Fields: []*Field{defaultPrimaryKeyField()},
	}
}

func Parse(b []byte) (*Config, error) {
	config := map[string]interface{}{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, errors.New("could not unmarshal config file")
	}
	return parse(config)
}

func parse(m map[string]interface{}) (*Config, error) {
	c := &Config{}

	if _, ok := m["databases"]; ok {
		for k, v := range m["databases"].(map[interface{}]interface{}) {
			database, err := parseDatabase(k.(string), v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			c.Databases = append(c.Databases, database)
		}
	} else {
		return nil, errors.New("no databases found in config file")
	}

	return c, nil
}

func parseDatabases(m map[interface{}]interface{}) ([]*Database, error) {
	var databases []*Database

	return databases, nil
}

func parseDatabase(name string, m map[interface{}]interface{}) (*Database, error) {
	database := defaultDatabase()
	database.Name = name
	if _, ok := m["schemas"]; ok {
		for k, v := range m["schemas"].(map[interface{}]interface{}) {
			schema, err := parseSchema(k.(string), v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			database.Schemas = append(database.Schemas, schema)
		}
	}
	if len(database.Schemas) == 0 {
		schema := defaultSchema()
		schema.Name = "public"
		if _, ok := m["models"]; ok {
			for k, v := range m["models"].(map[interface{}]interface{}) {
				model, err := parseModel(k.(string), v.(map[interface{}]interface{}))
				if err != nil {
					return nil, err
				}
				schema.Models = append(schema.Models, model)
			}
			if len(schema.Models) == 0 {
				return nil, errors.New("no models or schemas found in '" + database.Name + "' database")
			}
			database.Schemas = append(database.Schemas, schema)
		}
	}
	return database, nil
}

func parseSchema(name string, m map[interface{}]interface{}) (*Schema, error) {
	schema := defaultSchema()
	schema.Name = name
	if _, ok := m["models"]; ok {
		for k, v := range m["models"].(map[interface{}]interface{}) {
			model, err := parseModel(k.(string), v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			schema.Models = append(schema.Models, model)
		}
	} else {
		return nil, errors.New("no models found in '" + schema.Name + "' schema")
	}
	return schema, nil
}

func parseModel(name string, m map[interface{}]interface{}) (*Model, error) {
	model := defaultModel()
	model.Name = name

	// fields
	if _, ok := m["fields"]; ok {
		for k, v := range m["fields"].(map[interface{}]interface{}) {
			field, err := parseField(k.(string), v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			model.Fields = append(model.Fields, field)
		}
	} else {
		return nil, errors.New("no fields found in '" + model.Name + "' model")
	}

	// relations
	if _, ok := m["relations"]; ok {
		for k, v := range m["relations"].(map[interface{}]interface{}) {
			relation, err := parseRelation(k.(string), v.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			model.Relations = append(model.Relations, relation)
		}
	}

	// callbacks
	if _, ok := m["callbacks"]; ok {
		callbacks, ok := parseModelCallbacks(m["callbacks"])
		if !ok {
			return nil, errors.New("could not parse 'callbacks' for '" + model.Name + "' model")
		}
		model.Callbacks = callbacks
	}
	return model, nil
}

func finalizeModel(model *Model, m map[string]interface{}) error {
	return nil
}

func parseTable(name string, m map[string]interface{}) (*Table, error) {
	return nil, nil
}

func parseField(name string, m map[interface{}]interface{}) (*Field, error) {
	field := defaultField()
	field.Name = name

	// type
	tyi, ok := m["type"]
	if !ok {
		return nil, errors.New("'type' is not defined for '" + field.Name + "' field")
	}
	ty, ok := tyi.(string)
	if !ok {
		return nil, errors.New("could not parse 'type' for '" + field.Name + "' field into string")
	}
	field.TypeDB = ty

	// default
	def, ok := m["default"]
	if ok {
		field.Default = def
	}

	// null
	nuli, ok := m["null"]
	if ok {
		nul, ok := nuli.(bool)
		if !ok {
			return nil, errors.New("could not parse 'null' for '" + field.Name + "' field into bool")
		}
		field.IsNull = nul
	}

	// unique
	unii, ok := m["unique"]
	if ok {
		uni, ok := unii.(bool)
		if !ok {
			return nil, errors.New("could not parse 'unique' for '" + field.Name + "' field into bool")
		}
		field.IsUnique = uni
	}

	// inDB
	inDBi, ok := m["inDB"]
	if ok {
		inDB, ok := inDBi.(bool)
		if !ok {
			return nil, errors.New("could not parse 'inDB' for '" + field.Name + "' field into bool")
		}
		field.IsInDB = inDB
	}

	// nameInDB
	nDBi, ok := m["nameInDB"]
	if ok {
		nDB, ok := nDBi.(string)
		if !ok {
			return nil, errors.New("could not parse 'nameInDB' for '" + field.Name + "' field into string")
		}
		field.NameDB = nDB
	}

	// --- sql shortcuts
	// orderBy
	ordBi, ok := m["orderBy"]
	if ok {
		ordB, ok := ordBi.(bool)
		if !ok {
			return nil, errors.New("could not parse 'orderBy' for '" + field.Name + "' field into bool")
		}
		field.OrderBy = ordB
	}

	// groupBy
	grpBi, ok := m["groupBy"]
	if ok {
		grpB, ok := grpBi.(bool)
		if !ok {
			return nil, errors.New("could not parse 'groupBy' for '" + field.Name + "' field into bool")
		}
		field.GroupBy = grpB
	}

	// callbacks
	if _, ok := m["callbacks"]; ok {
		callbacks, ok := parseFieldCallbacks(m["callbacks"])
		if !ok {
			return nil, errors.New("could not parse 'callbacks' for '" + field.Name + "' field")
		}
		field.Callbacks = callbacks
	}

	// parse field wheres

	return field, nil
}

func parseModelCallbacks(i interface{}) ([]string, bool) {
	s, ok := i.([]interface{})
	if !ok {
		return nil, false
	}
	var callbacks []string
	for _, ci := range s {
		c, ok := ci.(string)
		if !ok {
			return nil, false
		}

		if !checkModelCallback(c) {
			continue
		}

		var found bool
		for _, ca := range callbacks {
			if c == ca {
				found = true
				break
			}
		}
		if !found {
			callbacks = append(callbacks, c)
		}
	}
	return callbacks, true
}

var modelCalbacks = []string{
	"beforeSave", "afterSave",
	"beforeUpdate", "afterUpdate",
	"beforeInsert", "afterInsert",
	"beforeDelete", "afterDelete",
	"beforeSoftDelete", "afterSoftDelete"}

func checkModelCallback(c string) bool {
	var found bool
	for _, ca := range modelCalbacks {
		if c == ca {
			found = true
			break
		}
	}
	return found
}

func parseFieldCallbacks(i interface{}) ([]string, bool) {
	s, ok := i.([]interface{})
	if !ok {
		return nil, false
	}
	var callbacks []string
	for _, ci := range s {
		c, ok := ci.(string)
		if !ok {
			return nil, false
		}

		if !checkFieldCallback(c) {
			continue
		}

		var found bool
		for _, ca := range callbacks {
			if c == ca {
				found = true
				break
			}
		}
		if !found {
			callbacks = append(callbacks, c)
		}
	}
	return callbacks, true
}

var fieldCalbacks = []string{
	"beforeSet", "afterSet",
	"beforeGet", "afterGet"}

func checkFieldCallback(c string) bool {
	var found bool
	for _, ca := range fieldCalbacks {
		if c == ca {
			found = true
			break
		}
	}
	return found
}

func parseRelation(name string, m map[interface{}]interface{}) (*Relation, error) {
	return nil, nil
}
