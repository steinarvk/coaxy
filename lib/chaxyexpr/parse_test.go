package chaxyexpr

import "testing"

func TestValidParses(t *testing.T) {
	testcases := []string{
		"hello",
		"hello_world",
		"$[1]",
		"0",
		"1",
		"0.1.2",
		"hello.world.hi",
		"hello.world.0.hi",
		`foo.bar.baz["quux"][1]`,
		`$[2].foo.bar.baz["quux"][1]`,
		`2.foo.bar.baz["quux"][1]`,
		`$["foo"].foo.bar.baz["quux"][1]`,
		`$.foo.foo.bar.baz["quux"][1]`,
		`foo.bar.baz["name.with.dots"][1]`,
		`foo.bar.baz["name with spaces"][1]`,
	}

	for _, testcase := range testcases {
		expr, err := Parse(testcase)
		if err != nil {
			t.Errorf("Parse(%q) = err: %v", testcase, err)
			continue
		}

		canonical := expr.FormatExpression()
		reparsed, err := Parse(canonical)
		if err != nil {
			t.Errorf("Parse(%q) [reformat of %q] = err: %v", canonical, testcase, err)
			continue
		}

		canonicalAgain := reparsed.FormatExpression()

		if canonical != canonicalAgain {
			t.Errorf("did not reach stability: %q => %q => %q", testcase, canonical, canonicalAgain)
		}
	}
}
