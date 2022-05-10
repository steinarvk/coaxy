package coaxyexpr

import (
	"strings"

	"github.com/steinarvk/coaxy/lib/accessor"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

type filter struct {
	name string
}

func (f filter) FormatFilter() string {
	return ":" + f.name
}

type expression struct {
	priorChunks []*expression
	filters     []*filter
	path        record.Path
}

func (x *expression) flushFilters() {
	if len(x.filters) > 0 {
		x.priorChunks = append(x.priorChunks, &expression{
			path:    x.path,
			filters: x.filters,
		})
		x.path = nil
		x.filters = nil
	}
}

func (x *expression) addIndex(i int) {
	x.flushFilters()
	x.path = x.path.Append(record.Index(i))
}

func (x *expression) addKey(k string) {
	x.flushFilters()
	x.path = x.path.Append(record.Field(k))
}

func (x *expression) addFilter(filt string) {
	x.filters = append(x.filters, &filter{
		name: filt,
	})
}

func (x *expression) FormatExpression() string {
	var chunks []string

	for _, chunk := range x.priorChunks {
		chunkformatted := chunk.FormatExpression()
		chunks = append(chunks, chunkformatted)
	}

	if len(x.path) > 0 {
		path, err := x.path.PathExpression()
		if err != nil {
			panic(err)
		}
		chunks = append(chunks, path)
	}

	for _, filt := range x.filters {
		chunks = append(chunks, filt.FormatFilter())
	}

	return strings.Join(chunks, " | ")
}

func (x *expression) MakeAccessor() interfaces.Accessor {
	var accessors []interfaces.Accessor

	for _, comp := range x.path {
		accessors = append(accessors, comp.MakeAccessor())
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
