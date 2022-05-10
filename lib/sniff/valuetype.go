package sniff

import (
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

func visitRecordValues(rec interfaces.Record, visit func(record.Path, string) error) error {
	value, err := rec.AsValue()
	if err != nil && !record.IsNotAPrimitive(err) {
		return err
	} else {
		visit(nil, value)
	}

	for _, index := range rec.Indices() {
		child, err := rec.GetByIndex(index)
		if err != nil {
			return err
		}

		if err := visitRecordValues(child, func(key record.Path, value string) error {
			return visit(key.Prepend(record.Index(index)), value)
		}); err != nil {
			return err
		}
	}

	for _, name := range rec.FieldNames() {
		child, err := rec.GetByName(name)
		if err != nil {
			return err
		}

		if err := visitRecordValues(child, func(key record.Path, value string) error {
			return visit(key.Prepend(record.Field(name)), value)
		}); err != nil {
			return err
		}
	}

	return nil
}

func collectNestedValues(records []interfaces.Record) (map[string][]string, error) {
	seen := map[string]bool{}

	for i := range records {
		if err := visitRecordValues(records[i], func(key record.Path, value string) error {
			path, err := key.PathExpression()
			if err != nil {
				return err
			}

			seen[path] = true
			return nil
		}); err != nil {
			return nil, err
		}
	}

	valuemap := map[string][]string{}

	for i := range records {
		seenlocal := map[string]bool{}

		if err := visitRecordValues(records[i], func(key record.Path, value string) error {
			path, err := key.PathExpression()
			if err != nil {
				return err
			}

			seenlocal[path] = true
			valuemap[path] = append(valuemap[path], value)
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
