package chaxyexpr

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/steinarvk/chaxy/lib/accessor"
	"github.com/steinarvk/chaxy/lib/interfaces"
)

var unquotedFieldRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z_0-9]*$`)

func fieldMustBeQuoted(s string) bool {
	return !unquotedFieldRE.MatchString(s)
}

type Expr interface {
	FormatExpression() string
	MakeAccessor() interfaces.Accessor
}

type rootIndex int

func (s rootIndex) FormatExpression() string {
	return fmt.Sprintf("%d", int(s))
}

func (s rootIndex) MakeAccessor() interfaces.Accessor {
	return accessor.AtIndex(int(s))
}

type rootField string

func (s rootField) FormatExpression() string {
	str := string(s)
	if fieldMustBeQuoted(str) {
		quoted := strconv.Quote(str)
		return fmt.Sprintf("$[%s]", quoted)
	}
	return str
}

func (s rootField) MakeAccessor() interfaces.Accessor {
	return accessor.AtField(string(s))
}

type fieldAccess struct {
	expr  Expr
	field string
}

func (a fieldAccess) FormatExpression() string {
	prior := a.expr.FormatExpression()

	if fieldMustBeQuoted(a.field) {
		quoted := strconv.Quote(a.field)
		return fmt.Sprintf("%s[%s]", prior, quoted)
	}

	return prior + "." + a.field
}

func (a fieldAccess) MakeAccessor() interfaces.Accessor {
	return accessor.Chain([]interfaces.Accessor{
		a.expr.MakeAccessor(),
		accessor.AtField(a.field),
	})
}

type indexAccess struct {
	expr  Expr
	index int
}

func (a indexAccess) MakeAccessor() interfaces.Accessor {
	return accessor.Chain([]interfaces.Accessor{
		a.expr.MakeAccessor(),
		accessor.AtIndex(a.index),
	})
}

func (a indexAccess) FormatExpression() string {
	return a.expr.FormatExpression() + fmt.Sprintf("[%d]", a.index)
}

func Parse(s string) (Expr, error) {
	parser := &parser{
		Buffer: s,
	}

	parser.Init()

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing expression %q: %w", s, err)
	}

	parser.Execute()

	return parser.expr, nil
}
