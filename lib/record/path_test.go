package record

import "testing"

func TestPathExpressions(t *testing.T) {
	testcases := map[string]Path{
		"foo.bar.baz": MakePath(Field("foo"), Field("bar"), Field("baz")),
		"foo[42].baz": MakePath(Field("foo"), Index(42), Field("baz")),
		"42.foo[0]":   MakePath(Index(42), Field("foo"), Index(0)),
	}

	for want, path := range testcases {
		got, err := path.PathExpression()
		if err != nil {
			t.Errorf("(%v).PathExpression() = err: %v", path, err)
			continue
		}

		if got != want {
			t.Errorf("(%v).PathExpression() = %q want %q", path, got, want)
		}
	}
}
