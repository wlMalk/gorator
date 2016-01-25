package parser

var relationTypes = []string{
	"hasOne", "belongsTo",
	"hasMany", "belongsToMany",
	"through"}

func (r *Relation) parse(name string, m map[interface{}]interface{}) error {
	r.def()
	r.Name = name

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

func (r *Relation) parseType(m map[interface{}]interface{}) error {
	return nil
}

func (r *Relation) parseModel(m map[interface{}]interface{}) error {
	return nil
}

func (r *Relation) parseForeignKey(m map[interface{}]interface{}) error {
	return nil
}

func (r *Relation) parseOtherKey(m map[interface{}]interface{}) error {
	return nil
}

func (r *Relation) parseDefaultQuery(m map[interface{}]interface{}) error {
	return nil
}
