package parser

func (p *PrimaryKey) GetFields() []*Field {
	var fields []*Field
	for _, pk := range p.Fields {
		for _, f := range p.Model.Fields {
			if pk == f.Name {
				fields = append(fields, f)
			}
		}
	}
	return fields
}
