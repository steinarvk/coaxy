package record

import "github.com/steinarvk/coaxy/lib/interfaces"

type primitiveValueRecord struct {
	value string
}

func (v primitiveValueRecord) AsValue() (string, error) {
	return v.value, nil
}

func (v primitiveValueRecord) GetByIndex(index int) (interfaces.Record, error) {
	return nil, errNotIndexable
}

func (v primitiveValueRecord) GetByName(name string) (interfaces.Record, error) {
	return nil, errNoFields
}

func (v primitiveValueRecord) Indices() []int       { return nil }
func (v primitiveValueRecord) FieldNames() []string { return nil }
