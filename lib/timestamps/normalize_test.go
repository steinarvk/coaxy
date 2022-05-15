package timestamps

import (
	"testing"
)

func TestSimpleNormalizing(t *testing.T) {
	testcases := []string{
		"1652654573",
		"1652654573000",
		"1652654573000000",
		"1652654573000000000",
		"1652654573.0",
		"1652654573000.0",
		"1652654573000000.0",
		"1652654573000000000.0",
		"2022-05-15T22:42:53Z",
		"2022-05-15 22:42:53Z",
		"2022-05-15 22:42:53.0Z",
		"2022-05-15 22:42:53.000000Z",
		"2022-05-15 22:42:53.000000+00:00",
		"2022-05-15T22:42:53Z",
		"2022-05-15T22:42:53.0Z",
		"2022-05-15T22:42:53.000000Z",
		"2022-05-15T22:42:53.000000+00:00",
		"May 15 2022 10:42:53:000PM",
		"May 15 2022 10:42:53.000PM",
		"May 15 2022 10:42:53PM",
		"May 15 2022 22:42:53",
	}

	want := "2022-05-15T22:42:53Z"

	norm := NewNormalizerISO()

	for _, tc := range testcases {
		got, err := norm(tc)
		if err != nil {
			t.Errorf("normalize(%q) = err: %v", tc, err)
			continue
		}

		if got != want {
			t.Errorf("normalize(%q) = %q want %q", tc, got, want)
		}
	}
}

func TestRecognizeTimestamps(t *testing.T) {
	testcases := []string{
		"Sep  1 1999 07:43:33.590PM",
		"Sep  1 1999 07:43:33:590PM",
	}

	norm := NewNormalizerISO()

	for _, tc := range testcases {
		_, err := norm(tc)
		if err != nil {
			t.Errorf("normalize(%q) = err: %v", tc, err)
		}
	}
}
