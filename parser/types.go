package parser

import (
	"github.com/wlMalk/gorator/driver"
)

type Config struct {
	Path           string
	Version        string
	GoratorVersion string

	Packages map[string]*Package

	Databases []*Database

	EncodingJSON          bool
	EncodingJSONOmitEmpty bool
	EncodingJSONUseStdLib bool
	EncodingXML           bool
	EncodingXMLOmitEmpty  bool
	EncodingXMLUseStdLib  bool
}

type Package struct {
	Config      *Config
	Name        string
	Description string
	Path        string
	Imports     []map[string]string
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
	Tables   []*Table
}

type Table struct {
	Name   string
	Schema string
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

	PrimaryKey *PrimaryKey

	CreatedAt bool
	UpdatedAt bool
	DeletedAt bool
	CreatedBy bool
	UpdatedBy bool
	DeletedBy bool

	Callbacks []string

	HoldOriginal bool
	SoftDelete   bool
	AllowExtra   bool

	Uuid   int
	Listed bool
	List   *List

	IsPivot bool

	EncodingJSON          bool
	EncodingJSONOmitEmpty bool
	EncodingJSONUseStdLib bool
	EncodingXML           bool
	EncodingXMLOmitEmpty  bool
	EncodingXMLUseStdLib  bool
}

type PrimaryKey struct {
	Model  *Model
	Fields []string
}

type Field struct {
	Model          *Model
	Name           string
	NameInDB       string
	Type           string
	TypeInDB       string
	NameInEncoding string

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
	InEncoding   bool
	Primitive    bool

	Where   bool
	OrderBy bool
	GroupBy bool
	Having  bool
}

type List struct {
	Model *Model
	Linq  bool
}

type Relation struct {
	Model          *Model
	Name           string
	NameInEncoding string
	OtherModelName string
	OtherModel     *Model
	Type           string
	ForeignKey     string
	LocalKey       string
	OtherKey       string
	PivotName      string
	Query          map[string]interface{}
	// following is for belongsToMany relations
	// OtherForeignKey []string
	// Key             []string
	Pivot *Model
}
