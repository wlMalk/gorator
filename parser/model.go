package parser

import (
	"fmt"
)

const (
	modelFields       = "fields"
	modelRelations    = "relations"
	modelCallbacks    = "callbacks"
	modelPrimaryKey   = "primaryKey"
	modelSoftDelete   = "softDelete"
	modelAllowExtra   = "allowExtra"
	modelTimestamps   = "timestamps"
	modelBy           = "by"
	modelHoldOriginal = "holdOriginal"
	modelSlice        = "slice"
	modelUuid         = "uuid"
)

func (mo *Model) parse(name string, m map[interface{}]interface{}) error {
	mo.def()
	mo.Name = name

	err := mo.parseRelations(m)
	if err != nil {
		return err
	}

	err = mo.parseSoftDelete(m)
	if err != nil {
		return err
	}

	err = mo.parseAllowExtra(m)
	if err != nil {
		return err
	}

	err = mo.parseTimestamps(m)
	if err != nil {
		return err
	}

	err = mo.parseBy(m)
	if err != nil {
		return err
	}

	err = mo.parseUuid(m)
	if err != nil {
		return err
	}

	err = mo.parseFields(m)
	if err != nil {
		return err
	}

	err = mo.parseCallbacks(m)
	if err != nil {
		return err
	}

	err = mo.parsePrimaryKey(m)
	if err != nil {
		return err
	}

	err = mo.parseHoldOriginal(m)
	if err != nil {
		return err
	}

	err = mo.parseSlice(m)
	if err != nil {
		return err
	}

	err = mo.finalize(m)
	if err != nil {
		return err
	}

	return nil
}

func (mo *Model) parseFields(m map[interface{}]interface{}) error {
	if _, ok := m[modelFields]; ok {
		for k, v := range mi(m[modelFields]) {
			field := &Field{}
			field.Model = mo
			err := field.parse(s(k), mi(v))
			if err != nil {
				return err
			}
			mo.Fields = append(mo.Fields, field)
		}
	} else {
		return fmt.Errorf("no '%s' found in '%s' model", modelFields, mo.Name)
	}
	return nil
}

func (mo *Model) parseRelations(m map[interface{}]interface{}) error {
	if _, ok := m[modelRelations]; ok {
		for k, v := range mi(m[modelRelations]) {
			relation := &Relation{}
			err := relation.parse(s(k), mi(v))
			if err != nil {
				return err
			}
			mo.Relations = append(mo.Relations, relation)
		}
	}
	return nil
}

var mCallbacks = []string{
	"beforeSave", "afterSave",
	"beforeSelect", "afterSelect",
	"beforeUpdate", "afterUpdate",
	"beforeInsert", "afterInsert",
	"beforeDelete", "afterDelete",
	"beforeSoftDelete", "afterSoftDelete"}

func (mo *Model) parseCallbacks(m map[interface{}]interface{}) error {
	if _, ok := m[modelCallbacks]; ok {
		s, ok := m[modelCallbacks].([]interface{}) // could be a string too
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' model", modelCallbacks, mo.Name)
		}

		for _, ci := range s {
			c, ok := ci.(string)
			if !ok {
				return fmt.Errorf("could not parse '%s' for '%s' model", modelCallbacks, mo.Name)
			}

			err := mo.checkCallback(c)
			if err != nil {
				return err
			}

			var found bool
			for _, ca := range mo.Callbacks {
				if c == ca {
					found = true
					break
				}
			}
			if !found {
				mo.Database.Config.Packages["orm"].Imports[importsInternal][mo.Database.Config.Path+"/database/orm/internal/callback"] = ""
				mo.Database.Config.addPackage(packageCallback(mo.Database.Config))
				mo.Callbacks = append(mo.Callbacks, c)
			}
		}
	}
	return nil
}

func (mo *Model) checkCallback(c string) error {
	var found bool
	for _, ca := range mCallbacks {
		if c == ca {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("'%s' callback is invalid for '%s' model", c, mo.Name)
	}

	return nil
}

func (mo *Model) parsePrimaryKey(m map[interface{}]interface{}) error {
	if _, ok := m[modelPrimaryKey]; ok {
		for _, v := range si(m[modelPrimaryKey]) {
			av := v.(string)

			foundM := false
			for _, f := range mo.Fields {
				if av == f.Name {
					foundM = true
					break
				}
			}
			if !foundM {
				return fmt.Errorf("no '%s' field found in '%s' model for '%s'", av, mo.Name, modelPrimaryKey)
			}

			foundPK := false
			for _, pk := range mo.PrimaryKey.Fields {
				if av == pk {
					foundPK = true
					break
				}
			}
			if !foundPK {
				mo.PrimaryKey.Fields = append(mo.PrimaryKey.Fields, av)
			}
		}
	}
	return nil
}

func (mo *Model) parseSoftDelete(m map[interface{}]interface{}) error {
	avi, ok := m[modelSoftDelete]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' model", modelSoftDelete, mo.Name)
		}
		mo.SoftDelete = av
	}
	return nil
}

func (mo *Model) parseAllowExtra(m map[interface{}]interface{}) error {
	avi, ok := m[modelAllowExtra]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' model", modelAllowExtra, mo.Name)
		}
		mo.AllowExtra = av
	}
	return nil
}

var mStamps = []string{
	"created",
	"updated",
	"deleted",
}

func (mo *Model) parseTimestamps(m map[interface{}]interface{}) error {
	avi, ok := m[modelTimestamps]
	if ok {
		av, ok := avi.(bool)
		if ok {
			mo.CreatedAt = av
			mo.UpdatedAt = av
			mo.DeletedAt = av
		} else {
			av, ok := avi.([]interface{})
			if !ok {
				return fmt.Errorf("could not parse '%s' for '%s' model", modelTimestamps, mo.Name)
			}

			mo.CreatedAt = false
			mo.UpdatedAt = false
			mo.DeletedAt = false

			for _, t := range av {
				ts := s(t)
				err := mo.checkTimestamp(ts)
				if err != nil {
					return err
				}

				switch ts {
				case mStamps[0]:
					mo.CreatedAt = true
				case mStamps[1]:
					mo.UpdatedAt = true
				case mStamps[2]:
					mo.DeletedAt = true
				}
			}
		}
	}
	return nil
}

func (mo *Model) checkTimestamp(c string) error {
	var found bool
	for _, ca := range mStamps {
		if c == ca {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("'%s' timestamp is invalid for '%s' model", c, mo.Name)
	}

	return nil
}

func (mo *Model) parseBy(m map[interface{}]interface{}) error {
	avi, ok := m[modelBy]
	if ok {
		av, ok := avi.(bool)
		if ok {
			mo.CreatedBy = av
			mo.UpdatedBy = av
			mo.DeletedBy = av
		} else {
			av, ok := avi.([]interface{})
			if !ok {
				return fmt.Errorf("could not parse '%s' for '%s' model", modelBy, mo.Name)
			}

			mo.CreatedBy = false
			mo.UpdatedBy = false
			mo.DeletedBy = false

			for _, b := range av {
				bs := s(b)
				err := mo.checkBy(bs)
				if err != nil {
					return err
				}

				switch bs {
				case mStamps[0]:
					mo.CreatedBy = true
				case mStamps[1]:
					mo.UpdatedBy = true
				case mStamps[2]:
					mo.DeletedBy = true
				}
			}
		}
	}
	return nil
}

func (mo *Model) checkBy(c string) error {
	var found bool
	for _, ca := range mStamps {
		if c == ca {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("'%s' stamp is invalid for '%s' model", c, mo.Name)
	}

	return nil
}

func (mo *Model) parseHoldOriginal(m map[interface{}]interface{}) error {
	avi, ok := m[modelHoldOriginal]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' model", modelHoldOriginal, mo.Name)
		}
		mo.HoldOriginal = av
	}
	return nil
}

func (mo *Model) parseSlice(m map[interface{}]interface{}) error {
	avi, ok := m[modelSlice]
	if ok {
		av, ok := avi.(bool)
		if !ok {
			return fmt.Errorf("could not parse '%s' for '%s' model", modelSlice, mo.Name)
		}
		mo.Sliced = av
	}
	return nil
}

func (mo *Model) parseUuid(m map[interface{}]interface{}) error {
	avi, ok := m[modelUuid]
	if ok {
		av, ok := avi.(int)
		if ok {
			if av <= 0 || av > 5 {
				return fmt.Errorf("'%s' version invalid for '%s' model", modelUuid, mo.Name)
			}
			mo.Uuid = av
		} else {
			av, ok := avi.(bool)
			if !ok {
				return fmt.Errorf("could not parse '%s' for '%s' model", modelUuid, mo.Name)
			}
			if av {
				mo.Uuid = defaultUuid
			} else {
				mo.Uuid = 0
			}
		}
	}
	return nil
}

func (mo *Model) finalize(m map[interface{}]interface{}) error {
	// types := mo.Database.Driver.Types()
	// for _, f := range mo.Fields {
	// 	f.Type = types[f.TypeInDB]
	// }

	return nil
}
