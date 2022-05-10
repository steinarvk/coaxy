package record

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/steinarvk/coaxy/lib/accessor"
	"github.com/steinarvk/coaxy/lib/interfaces"
)

var unquotedFieldRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z_0-9]*$`)

func fieldMustBeQuoted(s string) bool {
	return !unquotedFieldRE.MatchString(s)
}

type PathComponent struct {
	Index *int
	Field *string
}

func (c PathComponent) MakeAccessor() interfaces.Accessor {
	switch {
	case c.Index != nil:
		return accessor.AtIndex(*c.Index)

	case c.Field != nil:
		return accessor.AtField(*c.Field)

	default:
		return nil
	}
}

func Index(i int) PathComponent {
	return PathComponent{Index: &i}
}

func Field(name string) PathComponent {
	return PathComponent{Field: &name}
}

type Path []PathComponent

func (p Path) PathExpression() (string, error) {
	var rv string

	for _, pc := range p {
		switch {
		case pc.Index != nil:
			n := *pc.Index

			if rv == "" {
				rv = strconv.Itoa(n)
			} else {
				rv += "[" + strconv.Itoa(n) + "]"
			}

		case pc.Field != nil:
			name := *pc.Field

			if fieldMustBeQuoted(name) {
				rv += fmt.Sprintf("[%q]", name)
			} else {
				if rv == "" {
					rv = name
				} else {
					rv += "." + name
				}
			}

		default:
			return "", fmt.Errorf("invalid path component: %v", pc)
		}
	}

	return rv, nil
}

func MakePath(pc ...PathComponent) Path {
	return Path(pc)
}

func (p Path) Prepend(pc ...PathComponent) Path {
	return Path(append(pc, p...))
}

func (p Path) Append(pc ...PathComponent) Path {
	return Path(append(p, pc...))
}
