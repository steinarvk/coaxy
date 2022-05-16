package filters

import (
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/record"
)

type Filter struct {
	s string
	f func(string) (string, error)
}

func (t *Filter) String() string {
	return t.s
}

func (t *Filter) Extract(rec interfaces.Record) (interfaces.Record, error) {
	s, err := rec.AsValue()
	if err != nil {
		return nil, err
	}

	sPrime, err := t.f(s)
	if err != nil {
		return nil, err
	}

	return record.FromString(sPrime)
}
