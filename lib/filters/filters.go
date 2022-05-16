package filters

import (
	"errors"
	"fmt"
)

var (
	errNoSuchFilter = errors.New("no such filter")
)

type FilterEntry struct {
	name           string
	noArgsFun      func(string) (string, error)
	noArgsFunMaker func() func(string) (string, error)
}

func (f *FilterEntry) NoArgs() (*Filter, error) {
	if f == nil {
		return nil, errNoSuchFilter
	}

	if f.noArgsFun != nil {
		return &Filter{f.name, f.noArgsFun}, nil
	}

	if f.noArgsFunMaker != nil {
		return &Filter{f.name, f.noArgsFunMaker()}, nil
	}

	return nil, fmt.Errorf("wrong number of arguments")
}

func Lookup(name string) *FilterEntry {
	return filterbank[name]
}
