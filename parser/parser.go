package parser

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	defaultDriver = "pgsql"
	defaultUuid   = 4
)

func s(s interface{}) string {
	return s.(string)
}

func mi(m interface{}) map[interface{}]interface{} {
	return m.(map[interface{}]interface{})
}

func Parse(b []byte) (*Config, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("could not parse config file")
			os.Exit(0)
		}
	}()

	configMap := map[string]interface{}{}
	err := yaml.Unmarshal(b, &configMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file")
	}
	config := &Config{}
	err = config.parse(configMap)
	return config, err
}

const (
	configDatabases = "databases"
)

func (c *Config) parse(m map[string]interface{}) error {
	c.def()

	err := c.parseDatabases(m)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) parseDatabases(m map[string]interface{}) error {
	if _, ok := m[configDatabases]; ok {
		for k, v := range mi(m[configDatabases]) {
			database := &Database{}
			database.Config = c
			err := database.parse(s(k), mi(v))
			if err != nil {
				return err
			}
			c.Databases = append(c.Databases, database)
		}
	} else {
		return fmt.Errorf("no '%s' found in config file", configDatabases)
	}
	return nil
}
