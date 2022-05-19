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
	makerI         func(int) (func(string) (string, error), error)
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

func (f *FilterEntry) WithArgsInt(n int) (*Filter, error) {
	if f == nil {
		return nil, errNoSuchFilter
	}

	if f.makerI == nil {
		return nil, fmt.Errorf("wrong number or type of arguments")
	}

	callback, err := f.makerI(n)
	if err != nil {
		return nil, err
	}

	fullName := fmt.Sprintf("%s[%d]", f.name, n)

	return &Filter{fullName, callback}, nil
}

func Lookup(name string) *FilterEntry {
	return filterbank[name]
}
