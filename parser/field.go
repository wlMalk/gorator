package parser

import (
	"fmt"
)

var fieldAttrs = []string{"type", "default", "null", "unique", "inDB", "nameInDB", "orderBy", "groupBy", "callbacks"}

func (f *Field) parse(name string, m map[interface{}]interface{}) error {
	f.def()
	f.Name = name

	err := f.parseDefault(m)
	if err != nil {
		return err
	}

	err = f.parseType(m)
	if err != nil {
		return err
	}

	err = f.parseNull(m)
	if err != nil {
		return err
	}

	err = f.parseUnique(m)
	if err != nil {
		return err
	}

	err = f.parseInDB(m)
	if err != nil {
		return err
	}

	err = f.parseNameInDB(m)
	if err != nil {
		return err
	}

	err = f.parseWhere(m)
	if err != nil {
		return err
	}

	err = f.parseHaving(m)
	if err != nil {
		return err
	}

	err = f.parseOrderBy(m)
	if err != nil {
		return err
	}

	err = f.parseGroupBy(m)
	if err != nil {
		return err
	}

	err = f.parseCallbacks(m)
	if err != nil {
		return err
	}

	for _, a := range fieldAttrs {
		switch a {
		case "default":
			def, ok := m["default"]
			if ok {
				f.Default = def
			}
		}
	}

	// parse field wheres

	return nil
}

const (
	fieldDefault   = "default"
	fieldType      = "type"
	fieldNull      = "null"
	fieldUnique    = "unique"
	fieldInDB      = "inDB"
	fieldNameInDB  = "nameInDB"
	fieldWhere     = "where"
	fieldHaving    = "having"
	fieldOrderBy   = "orderBy"
	fieldGroupBy   = "groupBy"
	fieldCallbacks = "callbacks"
)

var fCallbacks = []string{
	"beforeSet", "afterSet",
	"beforeGet", "afterGet"}

func (f *Field) checkCallback(c string) error {
	var found bool
	for _, ca := range fCallbacks {
		if c == ca {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("'%s' callback is invalid for '%s' field", c, f.Name)
	}

	return nil
}

func (f *Field) parseDefault(m map[interface{}]interface{}) error {
	return nil
}

func (f *Field) parseType(m map[interface{}]interface{}) error {
	avi, ok := m[fieldType]
	if !ok {
		return fmt.Errorf("'%s' is not defined for '%s' field", fieldType, f.Name)
	}
	av, ok := avi.(string)
	if !ok {
		return fmt.Errorf("could not parse '%s' for '%s' field", fieldType, f.Name)
	}
	f.Type = av
	return nil
}

func (f *Field) parseNull(m map[interface{}]interface{}) error {
	avi, ok := m[fieldNull]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldNull, f.Name)
		}
		f.Null = av
	}
	return nil
}

func (f *Field) parseUnique(m map[interface{}]interface{}) error {
	avi, ok := m[fieldUnique]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldUnique, f.Name)
		}
		f.Unique = av
	}
	return nil
}

func (f *Field) parseInDB(m map[interface{}]interface{}) error {
	avi, ok := m[fieldInDB]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldInDB, f.Name)
		}
		f.InDB = av
	}
	return nil
}

func (f *Field) parseNameInDB(m map[interface{}]interface{}) error {
	avi, ok := m[fieldNameInDB]
	if ok {
		av, ok := avi.(string)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldNameInDB, f.Name)
		}
		f.NameInDB = av
	}
	return nil
}

func (f *Field) parseWhere(m map[interface{}]interface{}) error {
	avi, ok := m[fieldWhere]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldWhere, f.Name)
		}
		f.Where = av
	}
	return nil
}

func (f *Field) parseHaving(m map[interface{}]interface{}) error {
	avi, ok := m[fieldHaving]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldHaving, f.Name)
		}
		f.Having = av
	}
	return nil
}

func (f *Field) parseOrderBy(m map[interface{}]interface{}) error {
	avi, ok := m[fieldOrderBy]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldOrderBy, f.Name)
		}
		f.OrderBy = av
	}
	return nil
}

func (f *Field) parseGroupBy(m map[interface{}]interface{}) error {
	avi, ok := m[fieldGroupBy]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldGroupBy, f.Name)
		}
		f.GroupBy = av
	}
	return nil
}

func (f *Field) parseCallbacks(m map[interface{}]interface{}) error {
	if _, ok := m[fieldCallbacks]; ok {
		s, ok := m[fieldCallbacks].([]interface{}) // could be a string too
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' field", fieldCallbacks, f.Name)
		}

		for _, ci := range s {
			c, ok := ci.(string)
			if !ok {
				return fmt.Errorf("could not parse '%s' for '%s' field", fieldCallbacks, f.Name)
			}

			err := f.checkCallback(c)
			if err != nil {
				return err
			}

			var found bool
			for _, ca := range f.Callbacks {
				if c == ca {
					found = true
					break
				}
			}
			if !found {
				f.Callbacks = append(f.Callbacks, c)
			}
		}
	}
	return nil
}
