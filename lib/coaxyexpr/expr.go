package coaxyexpr

import (
	"fmt"
	"strings"

	"github.com/steinarvk/coaxy/lib/accessor"
	"github.com/steinarvk/coaxy/lib/interfaces"
)

type filter struct {
	name string
}

type component struct {
	index  *int
	key    *string
	filter *filter
}

func (c component) stringExpr() string {
	switch {
	case c.index != nil:
		return fmt.Sprintf("%d", *c.index)

	case c.key != nil:
		if fieldMustBeQuoted(*c.key) {
			return fmt.Sprintf("[%q]", *c.key)
		}
		return *c.key

	default:
		panic(fmt.Errorf("invalid accessor: %v", c))
	}
}

func (c component) makeAccessor() (interfaces.Accessor, error) {
	switch {
	case c.index != nil:
		return accessor.AtIndex(*c.index), nil

	case c.key != nil:
		return accessor.AtField(*c.key), nil

	default:
		return nil, fmt.Errorf("invalid accessor: %v", c)
	}
}

type expression struct {
	components []component
}

func (x *expression) addIndex(i int) {
	x.components = append(x.components, component{
		index: &i,
	})
}

func (x *expression) addKey(k string) {
	x.components = append(x.components, component{
		key: &k,
	})
}

func (x *expression) addSimpleFilter(name string) {
	x.components = append(x.components, component{
		filter: &filter{
			name: name,
		},
	})
}

func (x *expression) FormatExpression() string {
	var xs []string

	for _, c := range x.components {
		xs = append(xs, c.stringExpr())
	}

	return strings.Join(xs, ",")
}

func (x *expression) MakeAccessor() interfaces.Accessor {
	var accessors []interfaces.Accessor

	for _, comp := range x.components {
		acc, err := comp.makeAccessor()
		if err != nil {
			return accessor.Error{err}
		}
		accessors = append(accessors, acc)
	}

	return accessor.Func(func(rec interfaces.Record) (interfaces.Record, error) {
		for _, acc := range accessors {
			newrec, err := acc.Extract(rec)
			if err != nil {
				return nil, err
			}
			rec = newrec
		}
		return rec, nil
	})
}
