package coaxyexpr

import (
	"fmt"
	"strings"

	"github.com/steinarvk/coaxy/lib/accessor"
	"github.com/steinarvk/coaxy/lib/filters"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

type expression struct {
	priorChunks []*expression
	filters     []*filters.Filter
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

func (x *expression) addFilterI(name string, intArg int) {
	filt, err := filters.Lookup(name).WithArgsInt(intArg)
	if err != nil {
		panic(fmt.Errorf("error parsing filter %q: %v", name, err))
	}

	x.filters = append(x.filters, filt)
}

func (x *expression) addSimpleFilter(name string) {
	filt, err := filters.Lookup(name).NoArgs()
	if err != nil {
		panic(fmt.Errorf("error parsing filter %q: %v", name, err))
	}

	x.filters = append(x.filters, filt)
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
		chunks = append(chunks, filt.String())
	}

	return strings.Join(chunks, " | ")
}

func (x *expression) MakeAccessor() interfaces.Accessor {
	var accessors []interfaces.Accessor

	for _, priorChunk := range x.priorChunks {
		accessors = append(accessors, priorChunk.MakeAccessor())
	}

	for _, comp := range x.path {
		accessors = append(accessors, comp.MakeAccessor())
	}

	for _, filt := range x.filters {
		accessors = append(accessors, filt)
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
