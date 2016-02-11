package parser

import (
	"fmt"
)

const (
	relationType         = "type"
	relationModel        = "model"
	relationForeignKey   = "foreignKey"
	relationOtherKey     = "otherKey"
	relationDefaultQuery = "defaultQuery"
	relationThrough      = "through"
	relationAutoPreload  = "autoPreload"
	relationLoadable     = "loadable"
	relationPreloadable  = "preloadable"
	relationPivotTable   = "pivotTable"
	relationCallbacks    = "callbacks"
)

var rTypes = []string{
	"hasOne", "belongsTo",
	"hasMany", "belongsToMany"}

func (r *Relation) parse(name string, m map[interface{}]interface{}) error {
	r.def()
	r.Name = name
	r.NameInEncoding = name

	err := r.parseType(m)
	if err != nil {
		return err
	}

	err = r.parseModel(m)
	if err != nil {
		return err
	}

	err = r.parseForeignKey(m)
	if err != nil {
		return err
	}

	err = r.parseOtherKey(m) //conditional
	if err != nil {
		return err
	}

	err = r.parseDefaultQuery(m)
	if err != nil {
		return err
	}

	return nil
}

func (r *Relation) checkType(t string) error {
	var found bool
	for _, ta := range rTypes {
		if t == ta {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("'%s' type is invalid for '%s' relation", t, r.Name)
	}

	return nil
}

func (r *Relation) parseType(m map[interface{}]interface{}) error {
	avi, ok := m[relationType]
	if !ok {
		return fmt.Errorf("'%s' is not defined for '%s' relation", relationType, r.Name)
	}
	av, ok := avi.(string)
	if !ok {
		return fmt.Errorf("could not parse '%s' for '%s' relation", relationType, r.Name)
	}
	err := r.checkType(av)
	if err != nil {
		return err
	}
	r.Type = av
	return nil
}

func (r *Relation) parseModel(m map[interface{}]interface{}) error {
	// has to check for model is valid in finalize
	avi, ok := m[relationModel]
	if !ok {
		return fmt.Errorf("'%s' is not defined for '%s' relation", relationModel, r.Name)
	}
	av, ok := avi.(string)
	if !ok {
		return fmt.Errorf("could not parse '%s' for '%s' relation", relationModel, r.Name)
	}
	r.OtherModelName = av
	return nil
}

func (r *Relation) parseForeignKey(m map[interface{}]interface{}) error {
	avi, ok := m[relationForeignKey]
	if !ok {
		// generate name
		return nil
	}
	av, ok := avi.(string)
	if !ok {
		return fmt.Errorf("could not parse '%s' for '%s' relation", relationForeignKey, r.Name)
	}
	r.ForeignKey = av
	return nil
}

func (r *Relation) parseOtherKey(m map[interface{}]interface{}) error {
	// has to check for key in model is valid in finalize
	avi, ok := m[relationOtherKey]
	if !ok {
		// generate name
		return nil
	}
	av, ok := avi.(string)
	if !ok {
		return fmt.Errorf("could not parse '%s' for '%s' relation", relationOtherKey, r.Name)
	}
	r.OtherKey = av
	return nil
}

func (r *Relation) parseDefaultQuery(m map[interface{}]interface{}) error {
	return nil
}
