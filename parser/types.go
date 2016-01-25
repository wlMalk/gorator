package parser

import (
	"github.com/wlMalk/ormator/driver"
)

type Config struct {
	Path           string
	Version        string
	OrmatorVersion string

	Databases []*Database
}

type Database struct {
	Config     *Config
	Name       string
	DriverName string
	Driver     driver.Driver
	Models     []*Model
	Tables     []*Table
}

type Schema struct {
	Name     string
	Database *Database
	Models   []*Table
}

type Table struct {
	Name   string
	Schema *Schema
	Model  *Model

	IsPivot bool
}

type SuperModel struct {
	Name       string
	Table      *Table
	Fields     []*Field
	Relations  []*Relation
	SoftDelete bool

	CreatedAt bool
	UpdatedAt bool
	DeletedAt bool
	CreatedBy bool
	UpdatedBy bool
	DeletedBy bool

	Uuid int
}

type Model struct {
	Database *Database

	SuperModel *SuperModel
	Name       string
	Table      *Table
	Fields     []*Field
	Relations  []*Relation
	SoftDelete bool

	PrimaryKey *PrimaryKey

	CreatedAt bool
	UpdatedAt bool
	DeletedAt bool
	CreatedBy bool
	UpdatedBy bool
	DeletedBy bool

	Callbacks []string

	HoldOriginal bool

	Uuid   int
	Sliced bool
	Slice  *Slice

	IsPivot bool
}

type PrimaryKey struct {
	Model  *Model
	Fields []string
}

type Field struct {
	Model    *Model
	Name     string
	NameInDB string
	Type     string
	TypeInDB string

	Default interface{}

	Validations []string

	Callbacks []string

	Array        bool
	Null         bool
	Unique       bool
	Numeric      bool
	Incrementing bool
	InDB         bool
	Exported     bool

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
	Model      *Model
	Name       string
	OtherModel *Model
	Type       string
	ForeignKey string
	LocalKey   string
	OtherKey   string
	Pivot      *Model
}
