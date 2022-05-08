package chaxyexpr

// Code generated by peg lib/chaxyexpr/grammar.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleExpr
	ruleRoot
	rulePathComponent
	ruleRootField
	ruleRootIndex
	ruleBracketedRootIndex
	ruleBracketedRootField
	ruleFieldAccess
	ruleBracketedFieldAccess
	ruleBracketedIndexAccess
	ruleIndexAccess
	ruleDOT
	ruleLBRACKET
	ruleRBRACKET
	ruleDOLLAR
	ruleDQUOTE
	ruleBACKSLASH
	ruleStringLiteral
	ruleStringChar
	ruleEscape
	ruleInteger
	ruleIdentifier
	ruleIdentifierInitialChar
	ruleIdentifierContinuedChar
	ruleEND
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	rulePegText
	ruleAction8
	ruleAction9
	ruleAction10
)

var rul3s = [...]string{
	"Unknown",
	"Expr",
	"Root",
	"PathComponent",
	"RootField",
	"RootIndex",
	"BracketedRootIndex",
	"BracketedRootField",
	"FieldAccess",
	"BracketedFieldAccess",
	"BracketedIndexAccess",
	"IndexAccess",
	"DOT",
	"LBRACKET",
	"RBRACKET",
	"DOLLAR",
	"DQUOTE",
	"BACKSLASH",
	"StringLiteral",
	"StringChar",
	"Escape",
	"Integer",
	"Identifier",
	"IdentifierInitialChar",
	"IdentifierContinuedChar",
	"END",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"PegText",
	"Action8",
	"Action9",
	"Action10",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type parser struct {
	expr          Expr
	identifier    string
	stringLiteral string
	integer       int

	Buffer string
	buffer []rune
	rules  [38]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *parser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *parser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *parser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *parser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *parser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func (p *parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.expr = rootField(p.identifier)

		case ruleAction1:

			p.expr = rootIndex(p.integer)

		case ruleAction2:

			p.expr = rootIndex(p.integer)

		case ruleAction3:

			p.expr = rootField(p.stringLiteral)

		case ruleAction4:

			p.expr = fieldAccess{p.expr, p.identifier}

		case ruleAction5:

			p.expr = fieldAccess{p.expr, p.stringLiteral}

		case ruleAction6:

			p.expr = indexAccess{p.expr, p.integer}

		case ruleAction7:

			p.expr = indexAccess{p.expr, p.integer}

		case ruleAction8:

			s, _ := strconv.Unquote(text)
			p.stringLiteral = s

		case ruleAction9:

			n, _ := strconv.Atoi(text)
			p.integer = n

		case ruleAction10:

			p.identifier = text

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func Pretty(pretty bool) func(*parser) error {
	return func(p *parser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*parser) error {
	return func(p *parser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *parser) Init(options ...func(*parser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Expr <- <(Root PathComponent* END)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleRoot]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[rulePathComponent]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				if !_rules[ruleEND]() {
					goto l0
				}
				add(ruleExpr, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Root <- <(RootField / RootIndex / BracketedRootIndex / BracketedRootField)> */
		func() bool {
			position4, tokenIndex4 := position, tokenIndex
			{
				position5 := position
				{
					position6, tokenIndex6 := position, tokenIndex
					if !_rules[ruleRootField]() {
						goto l7
					}
					goto l6
				l7:
					position, tokenIndex = position6, tokenIndex6
					if !_rules[ruleRootIndex]() {
						goto l8
					}
					goto l6
				l8:
					position, tokenIndex = position6, tokenIndex6
					if !_rules[ruleBracketedRootIndex]() {
						goto l9
					}
					goto l6
				l9:
					position, tokenIndex = position6, tokenIndex6
					if !_rules[ruleBracketedRootField]() {
						goto l4
					}
				}
			l6:
				add(ruleRoot, position5)
			}
			return true
		l4:
			position, tokenIndex = position4, tokenIndex4
			return false
		},
		/* 2 PathComponent <- <(FieldAccess / IndexAccess / BracketedFieldAccess / BracketedIndexAccess)> */
		func() bool {
			position10, tokenIndex10 := position, tokenIndex
			{
				position11 := position
				{
					position12, tokenIndex12 := position, tokenIndex
					if !_rules[ruleFieldAccess]() {
						goto l13
					}
					goto l12
				l13:
					position, tokenIndex = position12, tokenIndex12
					if !_rules[ruleIndexAccess]() {
						goto l14
					}
					goto l12
				l14:
					position, tokenIndex = position12, tokenIndex12
					if !_rules[ruleBracketedFieldAccess]() {
						goto l15
					}
					goto l12
				l15:
					position, tokenIndex = position12, tokenIndex12
					if !_rules[ruleBracketedIndexAccess]() {
						goto l10
					}
				}
			l12:
				add(rulePathComponent, position11)
			}
			return true
		l10:
			position, tokenIndex = position10, tokenIndex10
			return false
		},
		/* 3 RootField <- <((DOLLAR DOT)? Identifier Action0)> */
		func() bool {
			position16, tokenIndex16 := position, tokenIndex
			{
				position17 := position
				{
					position18, tokenIndex18 := position, tokenIndex
					if !_rules[ruleDOLLAR]() {
						goto l18
					}
					if !_rules[ruleDOT]() {
						goto l18
					}
					goto l19
				l18:
					position, tokenIndex = position18, tokenIndex18
				}
			l19:
				if !_rules[ruleIdentifier]() {
					goto l16
				}
				if !_rules[ruleAction0]() {
					goto l16
				}
				add(ruleRootField, position17)
			}
			return true
		l16:
			position, tokenIndex = position16, tokenIndex16
			return false
		},
		/* 4 RootIndex <- <((DOLLAR DOT)? Integer Action1)> */
		func() bool {
			position20, tokenIndex20 := position, tokenIndex
			{
				position21 := position
				{
					position22, tokenIndex22 := position, tokenIndex
					if !_rules[ruleDOLLAR]() {
						goto l22
					}
					if !_rules[ruleDOT]() {
						goto l22
					}
					goto l23
				l22:
					position, tokenIndex = position22, tokenIndex22
				}
			l23:
				if !_rules[ruleInteger]() {
					goto l20
				}
				if !_rules[ruleAction1]() {
					goto l20
				}
				add(ruleRootIndex, position21)
			}
			return true
		l20:
			position, tokenIndex = position20, tokenIndex20
			return false
		},
		/* 5 BracketedRootIndex <- <(DOLLAR LBRACKET Integer RBRACKET Action2)> */
		func() bool {
			position24, tokenIndex24 := position, tokenIndex
			{
				position25 := position
				if !_rules[ruleDOLLAR]() {
					goto l24
				}
				if !_rules[ruleLBRACKET]() {
					goto l24
				}
				if !_rules[ruleInteger]() {
					goto l24
				}
				if !_rules[ruleRBRACKET]() {
					goto l24
				}
				if !_rules[ruleAction2]() {
					goto l24
				}
				add(ruleBracketedRootIndex, position25)
			}
			return true
		l24:
			position, tokenIndex = position24, tokenIndex24
			return false
		},
		/* 6 BracketedRootField <- <(DOLLAR LBRACKET StringLiteral RBRACKET Action3)> */
		func() bool {
			position26, tokenIndex26 := position, tokenIndex
			{
				position27 := position
				if !_rules[ruleDOLLAR]() {
					goto l26
				}
				if !_rules[ruleLBRACKET]() {
					goto l26
				}
				if !_rules[ruleStringLiteral]() {
					goto l26
				}
				if !_rules[ruleRBRACKET]() {
					goto l26
				}
				if !_rules[ruleAction3]() {
					goto l26
				}
				add(ruleBracketedRootField, position27)
			}
			return true
		l26:
			position, tokenIndex = position26, tokenIndex26
			return false
		},
		/* 7 FieldAccess <- <(DOT Identifier Action4)> */
		func() bool {
			position28, tokenIndex28 := position, tokenIndex
			{
				position29 := position
				if !_rules[ruleDOT]() {
					goto l28
				}
				if !_rules[ruleIdentifier]() {
					goto l28
				}
				if !_rules[ruleAction4]() {
					goto l28
				}
				add(ruleFieldAccess, position29)
			}
			return true
		l28:
			position, tokenIndex = position28, tokenIndex28
			return false
		},
		/* 8 BracketedFieldAccess <- <(LBRACKET StringLiteral RBRACKET Action5)> */
		func() bool {
			position30, tokenIndex30 := position, tokenIndex
			{
				position31 := position
				if !_rules[ruleLBRACKET]() {
					goto l30
				}
				if !_rules[ruleStringLiteral]() {
					goto l30
				}
				if !_rules[ruleRBRACKET]() {
					goto l30
				}
				if !_rules[ruleAction5]() {
					goto l30
				}
				add(ruleBracketedFieldAccess, position31)
			}
			return true
		l30:
			position, tokenIndex = position30, tokenIndex30
			return false
		},
		/* 9 BracketedIndexAccess <- <(LBRACKET Integer RBRACKET Action6)> */
		func() bool {
			position32, tokenIndex32 := position, tokenIndex
			{
				position33 := position
				if !_rules[ruleLBRACKET]() {
					goto l32
				}
				if !_rules[ruleInteger]() {
					goto l32
				}
				if !_rules[ruleRBRACKET]() {
					goto l32
				}
				if !_rules[ruleAction6]() {
					goto l32
				}
				add(ruleBracketedIndexAccess, position33)
			}
			return true
		l32:
			position, tokenIndex = position32, tokenIndex32
			return false
		},
		/* 10 IndexAccess <- <(DOT Integer Action7)> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				if !_rules[ruleDOT]() {
					goto l34
				}
				if !_rules[ruleInteger]() {
					goto l34
				}
				if !_rules[ruleAction7]() {
					goto l34
				}
				add(ruleIndexAccess, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 11 DOT <- <'.'> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				if buffer[position] != rune('.') {
					goto l36
				}
				position++
				add(ruleDOT, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 12 LBRACKET <- <'['> */
		func() bool {
			position38, tokenIndex38 := position, tokenIndex
			{
				position39 := position
				if buffer[position] != rune('[') {
					goto l38
				}
				position++
				add(ruleLBRACKET, position39)
			}
			return true
		l38:
			position, tokenIndex = position38, tokenIndex38
			return false
		},
		/* 13 RBRACKET <- <']'> */
		func() bool {
			position40, tokenIndex40 := position, tokenIndex
			{
				position41 := position
				if buffer[position] != rune(']') {
					goto l40
				}
				position++
				add(ruleRBRACKET, position41)
			}
			return true
		l40:
			position, tokenIndex = position40, tokenIndex40
			return false
		},
		/* 14 DOLLAR <- <'$'> */
		func() bool {
			position42, tokenIndex42 := position, tokenIndex
			{
				position43 := position
				if buffer[position] != rune('$') {
					goto l42
				}
				position++
				add(ruleDOLLAR, position43)
			}
			return true
		l42:
			position, tokenIndex = position42, tokenIndex42
			return false
		},
		/* 15 DQUOTE <- <'"'> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				if buffer[position] != rune('"') {
					goto l44
				}
				position++
				add(ruleDQUOTE, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 16 BACKSLASH <- <'\\'> */
		func() bool {
			position46, tokenIndex46 := position, tokenIndex
			{
				position47 := position
				if buffer[position] != rune('\\') {
					goto l46
				}
				position++
				add(ruleBACKSLASH, position47)
			}
			return true
		l46:
			position, tokenIndex = position46, tokenIndex46
			return false
		},
		/* 17 StringLiteral <- <(<(DQUOTE StringChar* DQUOTE)> Action8)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				{
					position50 := position
					if !_rules[ruleDQUOTE]() {
						goto l48
					}
				l51:
					{
						position52, tokenIndex52 := position, tokenIndex
						if !_rules[ruleStringChar]() {
							goto l52
						}
						goto l51
					l52:
						position, tokenIndex = position52, tokenIndex52
					}
					if !_rules[ruleDQUOTE]() {
						goto l48
					}
					add(rulePegText, position50)
				}
				if !_rules[ruleAction8]() {
					goto l48
				}
				add(ruleStringLiteral, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 18 StringChar <- <(Escape / (!('"' / '\n' / '\\') .))> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				{
					position55, tokenIndex55 := position, tokenIndex
					if !_rules[ruleEscape]() {
						goto l56
					}
					goto l55
				l56:
					position, tokenIndex = position55, tokenIndex55
					{
						position57, tokenIndex57 := position, tokenIndex
						{
							position58, tokenIndex58 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l59
							}
							position++
							goto l58
						l59:
							position, tokenIndex = position58, tokenIndex58
							if buffer[position] != rune('\n') {
								goto l60
							}
							position++
							goto l58
						l60:
							position, tokenIndex = position58, tokenIndex58
							if buffer[position] != rune('\\') {
								goto l57
							}
							position++
						}
					l58:
						goto l53
					l57:
						position, tokenIndex = position57, tokenIndex57
					}
					if !matchDot() {
						goto l53
					}
				}
			l55:
				add(ruleStringChar, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 19 Escape <- <(BACKSLASH ('"' / '\n' / '\\'))> */
		func() bool {
			position61, tokenIndex61 := position, tokenIndex
			{
				position62 := position
				if !_rules[ruleBACKSLASH]() {
					goto l61
				}
				{
					position63, tokenIndex63 := position, tokenIndex
					if buffer[position] != rune('"') {
						goto l64
					}
					position++
					goto l63
				l64:
					position, tokenIndex = position63, tokenIndex63
					if buffer[position] != rune('\n') {
						goto l65
					}
					position++
					goto l63
				l65:
					position, tokenIndex = position63, tokenIndex63
					if buffer[position] != rune('\\') {
						goto l61
					}
					position++
				}
			l63:
				add(ruleEscape, position62)
			}
			return true
		l61:
			position, tokenIndex = position61, tokenIndex61
			return false
		},
		/* 20 Integer <- <(<('0' / ([1-9] [0-9]*))> Action9)> */
		func() bool {
			position66, tokenIndex66 := position, tokenIndex
			{
				position67 := position
				{
					position68 := position
					{
						position69, tokenIndex69 := position, tokenIndex
						if buffer[position] != rune('0') {
							goto l70
						}
						position++
						goto l69
					l70:
						position, tokenIndex = position69, tokenIndex69
						if c := buffer[position]; c < rune('1') || c > rune('9') {
							goto l66
						}
						position++
					l71:
						{
							position72, tokenIndex72 := position, tokenIndex
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l72
							}
							position++
							goto l71
						l72:
							position, tokenIndex = position72, tokenIndex72
						}
					}
				l69:
					add(rulePegText, position68)
				}
				if !_rules[ruleAction9]() {
					goto l66
				}
				add(ruleInteger, position67)
			}
			return true
		l66:
			position, tokenIndex = position66, tokenIndex66
			return false
		},
		/* 21 Identifier <- <(<(IdentifierInitialChar IdentifierContinuedChar*)> Action10)> */
		func() bool {
			position73, tokenIndex73 := position, tokenIndex
			{
				position74 := position
				{
					position75 := position
					if !_rules[ruleIdentifierInitialChar]() {
						goto l73
					}
				l76:
					{
						position77, tokenIndex77 := position, tokenIndex
						if !_rules[ruleIdentifierContinuedChar]() {
							goto l77
						}
						goto l76
					l77:
						position, tokenIndex = position77, tokenIndex77
					}
					add(rulePegText, position75)
				}
				if !_rules[ruleAction10]() {
					goto l73
				}
				add(ruleIdentifier, position74)
			}
			return true
		l73:
			position, tokenIndex = position73, tokenIndex73
			return false
		},
		/* 22 IdentifierInitialChar <- <([a-z] / [A-Z])> */
		func() bool {
			position78, tokenIndex78 := position, tokenIndex
			{
				position79 := position
				{
					position80, tokenIndex80 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex = position80, tokenIndex80
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l78
					}
					position++
				}
			l80:
				add(ruleIdentifierInitialChar, position79)
			}
			return true
		l78:
			position, tokenIndex = position78, tokenIndex78
			return false
		},
		/* 23 IdentifierContinuedChar <- <([a-z] / [A-Z] / '_' / [0-9])> */
		func() bool {
			position82, tokenIndex82 := position, tokenIndex
			{
				position83 := position
				{
					position84, tokenIndex84 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l85
					}
					position++
					goto l84
				l85:
					position, tokenIndex = position84, tokenIndex84
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l86
					}
					position++
					goto l84
				l86:
					position, tokenIndex = position84, tokenIndex84
					if buffer[position] != rune('_') {
						goto l87
					}
					position++
					goto l84
				l87:
					position, tokenIndex = position84, tokenIndex84
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l82
					}
					position++
				}
			l84:
				add(ruleIdentifierContinuedChar, position83)
			}
			return true
		l82:
			position, tokenIndex = position82, tokenIndex82
			return false
		},
		/* 24 END <- <!.> */
		func() bool {
			position88, tokenIndex88 := position, tokenIndex
			{
				position89 := position
				{
					position90, tokenIndex90 := position, tokenIndex
					if !matchDot() {
						goto l90
					}
					goto l88
				l90:
					position, tokenIndex = position90, tokenIndex90
				}
				add(ruleEND, position89)
			}
			return true
		l88:
			position, tokenIndex = position88, tokenIndex88
			return false
		},
		/* 26 Action0 <- <{
		  p.expr = rootField(p.identifier)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 27 Action1 <- <{
		  p.expr = rootIndex(p.integer)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 28 Action2 <- <{
		  p.expr = rootIndex(p.integer)
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 29 Action3 <- <{
		  p.expr = rootField(p.stringLiteral)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 30 Action4 <- <{
		  p.expr = fieldAccess{p.expr, p.identifier}
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 31 Action5 <- <{
		  p.expr = fieldAccess{p.expr, p.stringLiteral}
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 32 Action6 <- <{
		  p.expr = indexAccess{p.expr, p.integer}
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 33 Action7 <- <{
		  p.expr = indexAccess{p.expr, p.integer}
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		nil,
		/* 35 Action8 <- <{
		  s, _ := strconv.Unquote(text)
		  p.stringLiteral = s
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 36 Action9 <- <{
		  n, _ := strconv.Atoi(text)
		  p.integer = n
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 37 Action10 <- <{
		  p.identifier = text
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
	}
	p.rules = _rules
	return nil
}
