package record

import (
	"errors"
	"fmt"
	"sort"

	"github.com/steinarvk/chaxy/lib/chaxyvalue"
)

var (
	errNotAPrimitive = errors.New("value is not a primitive")
	errNotIndexable  = errors.New("value is not indexable")
	errNoFields      = errors.New("value does not have fields")
)

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

type record interface {
	GetByIndex(int) (record, error)
	GetByName(string) (record, error)
	AsValue() (string, error)
}

type valueRecord struct {
	value string
}

func (v valueRecord) AsValue() (string, error) {
	return v.value, nil
}

func (v valueRecord) GetByIndex(index int) (record, error) {
	return nil, errNotIndexable
}

func (v valueRecord) GetByName(name string) (record, error) {
	return nil, errNoFields
}

type tupleRecord struct {
	indexByName map[string]int
	values      []string
}

func (t tupleRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func (t tupleRecord) GetByIndex(index int) (record, error) {
	if index < 0 || index >= len(t.values) {
		return nil, errOutOfBounds{index, len(t.values)}
	}
	return valueRecord{t.values[index]}, nil
}

func (t tupleRecord) GetByName(name string) (record, error) {
	if t.indexByName != nil {
		index, ok := t.indexByName[name]
		if ok {
			return t.GetByIndex(index)
		}
	}

	return nil, errNoSuchField{field: name}
}

type jsonObjectRecord struct {
	values map[string]interface{}
}

func (j jsonObjectRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func (j jsonObjectRecord) GetByIndex(int) (record, error) {
	return nil, errNotIndexable
}

func (j jsonObjectRecord) GetByName(k string) (record, error) {
	v, ok := j.values[k]
	if !ok {
		return nil, errNoSuchField{field: k}
	}
	return jsonValueToRecord(v)
}

type jsonArrayRecord struct {
	values []interface{}
}

func (j jsonArrayRecord) AsValue() (string, error) {
	return "", errNotAPrimitive
}

func (j jsonArrayRecord) GetByIndex(index int) (record, error) {
	if index < 0 || index >= len(j.values) {
		return nil, errOutOfBounds{index, len(j.values)}
	}
	return jsonValueToRecord(j.values[index])
}

func (j jsonArrayRecord) GetByName(string) (record, error) {
	return nil, errNoFields
}

func jsonValueToRecord(value interface{}) (record, error) {
	switch v := value.(type) {
	case []interface{}:
		return jsonArrayRecord{v}, nil

	case map[string]interface{}:
		return jsonObjectRecord{v}, nil

	default:
		asString, err := chaxyvalue.JSONPrimitiveToString(value)
		if err != nil {
			return nil, err
		}

		return valueRecord{asString}, nil
	}
}
