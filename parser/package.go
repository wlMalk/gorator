package parser

func packageDatabase(c *Config) *Package {
	p := &Package{
		Config:      c,
		Name:        "database",
		Description: "",
		Path:        "database",
		Imports: []map[string]string{
			map[string]string{
				"database/sql": "",
				"fmt":          "",
			},
			map[string]string{},
			map[string]string{
				"github.com/wlMalk/gorator/database": "odatabase",
			},
		},
	}
	return p
}

func packageORM(c *Config) *Package {
	p := &Package{
		Config:      c,
		Name:        "orm",
		Description: "",
		Path:        "database/orm",
		Imports: []map[string]string{
			map[string]string{},
			map[string]string{
				c.Path + "/database/orm/query": "",
				c.Path + "/database/orm/model": "",
			}, map[string]string{
				"github.com/wlMalk/gorator/database":       "odatabase",
				"github.com/wlMalk/gorator/database/query": "oquery",
			},
		},
	}
	return p
}

func packageQuery(c *Config) *Package {
	p := &Package{
		Config:      c,
		Name:        "query",
		Description: "",
		Path:        "database/orm/query",
		Imports: []map[string]string{
			map[string]string{
				"strings": "",
			}, map[string]string{
				c.Path + "/database": "",
			},
			map[string]string{
				"github.com/wlMalk/gorator/database/query": "oquery",
				"github.com/wlMalk/gorator/database":       "odatabase",
			},
		},
	}
	return p
}

func packageModel(c *Config) *Package {
	p := &Package{
		Config:      c,
		Name:        "model",
		Description: "",
		Path:        "database/orm/model",
		Imports: []map[string]string{
			map[string]string{
				// "encoding/json": "gojson",
				"bytes":   "",
				"fmt":     "",
				"strings": "",
			}, map[string]string{
				c.Path + "/database/orm/query": "_",
			},
			map[string]string{
				"github.com/wlMalk/gorator/database": "odatabase",
				"github.com/wlMalk/json":             "",
			},
		},
	}
	return p
}

func packageCallback(c *Config) *Package {
	p := &Package{
		Config:      c,
		Name:        "callback",
		Description: "",
		Path:        "database/orm/internal/callback",
		Imports: []map[string]string{
			map[string]string{},
			map[string]string{},
			map[string]string{
				"github.com/wlMalk/gorator/database/query": "oquery",
				"github.com/wlMalk/gorator/database":       "odatabase",
			},
		},
	}
	return p
}
