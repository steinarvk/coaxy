package sniff

import (
	"fmt"

	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func visitRecordValues(rec interfaces.Record, visit func(string, string) error) error {
	value, err := rec.AsValue()
	if err != nil && !record.IsNotAPrimitive(err) {
		return err
	} else {
		visit("", value)
	}

	for _, index := range rec.Indices() {
		child, err := rec.GetByIndex(index)
		if err != nil {
			return err
		}

		if err := visitRecordValues(child, func(key, value string) error {
			return visit(fmt.Sprintf("[%d]%s", index, key), value)
		}); err != nil {
			return err
		}
	}

	for _, name := range rec.FieldNames() {
		child, err := rec.GetByName(name)
		if err != nil {
			return err
		}

		if err := visitRecordValues(child, func(key, value string) error {
			return visit(fmt.Sprintf("[%q]%s", name, key), value)
		}); err != nil {
			return err
		}
	}

	return nil
}

func collectNestedValues(records []interfaces.Record) (map[string][]string, error) {
	seen := map[string]bool{}

	for i := range records {
		if err := visitRecordValues(records[i], func(key, value string) error {
			seen[key] = true
			return nil
		}); err != nil {
			return nil, err
		}
	}

	valuemap := map[string][]string{}

	for i := range records {
		seenlocal := map[string]bool{}

		if err := visitRecordValues(records[i], func(key, value string) error {
			seenlocal[key] = true
			valuemap[key] = append(valuemap[key], value)
			return nil
		}); err != nil {
			return nil, err
		}

		for k := range seen {
			if !seenlocal[k] {
				valuemap[k] = append(valuemap[k], "")
			}
		}
	}

	return valuemap, nil
}

func DetectNestedFields(records []interfaces.Record) (map[string]ValueType, error) {
	valuemap, err := collectNestedValues(records)
	if err != nil {
		return nil, err
	}

	rv := map[string]ValueType{}

	for path, values := range valuemap {
		valuetype := DetectValueType(values)
		rv[path] = valuetype
	}

	return rv, nil
}
