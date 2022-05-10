package record

import (
	"errors"
	"fmt"
	"sort"

	"github.com/steinarvk/coaxy/lib/coaxyvalue"
	"github.com/steinarvk/coaxy/lib/interfaces"
)

var (
	errNotAPrimitive = errors.New("value is not a primitive")
	errNotIndexable  = errors.New("value is not indexable")
	errNoFields      = errors.New("value does not have fields")
)

func IsNotAPrimitive(err error) bool {
	return err == errNotAPrimitive
}

type errNoSuchField struct {
	field  string
	fields []string
}

func (e errNoSuchField) Error() string {
	if len(e.fields) > 0 {
		sort.Strings(e.fields)
		return fmt.Sprintf("no such field: %q (valid fields: %v)", e.field, e.fields)
	}
	return fmt.Sprintf("no such field: %q", e.field)
}

type errOutOfBounds struct {
	index  int
	length int
}

func (e errOutOfBounds) Error() string {
	return fmt.Sprintf("index access [%d] out of bounds (length %d)", e.index, e.length)
}

type nullRecord struct{}

func (_ nullRecord) AsValue() (string, error)                    { return "", nil }
func (_ nullRecord) GetByIndex(int) (interfaces.Record, error)   { return nullRecord{}, nil }
func (_ nullRecord) GetByName(string) (interfaces.Record, error) { return nullRecord{}, nil }
func (_ nullRecord) Indices() []int                              { return nil }
func (_ nullRecord) FieldNames() []string                        { return nil }

type tupleRecord struct {
	indexByName map[string]int
	values      []string
}

func (t tupleRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func (t tupleRecord) GetByIndex(index int) (interfaces.Record, error) {
	if index < 0 || index >= len(t.values) {
		return nil, errOutOfBounds{index, len(t.values)}
	}
	return FromString(t.values[index])
}

func (t tupleRecord) GetByName(name string) (interfaces.Record, error) {
	if t.indexByName != nil {
		index, ok := t.indexByName[name]
		if ok {
			return t.GetByIndex(index)
		}
	}

	return nil, errNoSuchField{field: name}
}

func (t tupleRecord) Indices() []int {
	namedFields := map[int]bool{}
	if t.indexByName != nil {
		for _, index := range t.indexByName {
			namedFields[index] = true
		}
	}

	var rv []int
	for i := 0; i < len(t.values); i++ {
		if !namedFields[i] {
			rv = append(rv, i)
		}
	}

	return rv
}

func (t tupleRecord) FieldNames() []string {
	var rv []string

	if t.indexByName != nil {
		for k, _ := range t.indexByName {
			rv = append(rv, k)
		}
	}

	sort.Strings(rv)
	return rv
}

type jsonObjectRecord struct {
	values map[string]interface{}
}

func (j jsonObjectRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func (j jsonObjectRecord) GetByIndex(int) (interfaces.Record, error) {
	return nil, errNotIndexable
}

func accessJSONObjectByName(m map[string]interface{}, key string) (interfaces.Record, error) {
	v, ok := m[key]
	if !ok {
		return nullRecord{}, nil
	}
	return FromJSONValue(v)
}

func (j jsonObjectRecord) GetByName(k string) (interfaces.Record, error) {
	return accessJSONObjectByName(j.values, k)
}

func (j jsonObjectRecord) Indices() []int { return nil }
func (j jsonObjectRecord) FieldNames() []string {
	var rv []string
	for k, _ := range j.values {
		rv = append(rv, k)
	}
	sort.Strings(rv)
	return rv
}

type jsonArrayRecord struct {
	values []interface{}
}

func (j jsonArrayRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func accessJSONArrayByIndex(values []interface{}, index int) (interfaces.Record, error) {
	if index < 0 {
		return nil, errOutOfBounds{index, len(values)}
	}
	if index >= len(values) {
		return nullRecord{}, nil
	}
	return FromJSONValue(values[index])
}

func (j jsonArrayRecord) GetByIndex(index int) (interfaces.Record, error) {
	return accessJSONArrayByIndex(j.values, index)
}

func (j jsonArrayRecord) GetByName(string) (interfaces.Record, error) {
	return nil, errNoFields
}

func (j jsonArrayRecord) Indices() []int {
	rv := make([]int, len(j.values))
	for i := 0; i < len(j.values); i++ {
		rv[i] = i
	}
	return rv
}

func (j jsonArrayRecord) FieldNames() []string { return nil }

func FromJSONValue(value interface{}) (interfaces.Record, error) {
	switch v := value.(type) {
	case []interface{}:
		return jsonArrayRecord{v}, nil

	case map[string]interface{}:
		return jsonObjectRecord{v}, nil

	case string:
		return FromString(v)

	default:
		asString, err := coaxyvalue.JSONPrimitiveToString(value)
		if err != nil {
			return nil, err
		}

		return primitiveValueRecord{asString}, nil
	}
}

func FromString(value string) (interfaces.Record, error) {
	if value == "" {
		return primitiveValueRecord{""}, nil
	}

	return &stringValueRecord{value: value}, nil
}

func FromStringTuple(values []string, names map[string]int) (interfaces.Record, error) {
	return tupleRecord{
		indexByName: names,
		values:      values,
	}, nil
}

func StringTupleWithNames(rec interfaces.Record, names map[string]int) (interfaces.Record, error) {
	tr, ok := rec.(tupleRecord)
	if !ok {
		return nil, fmt.Errorf("not a string tuple")
	}

	tr.indexByName = names
	return tr, nil
}
