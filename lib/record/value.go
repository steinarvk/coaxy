package record

import (
	"encoding/json"
	"regexp"
	"sort"
	"sync"

	"github.com/kr/logfmt"
	"github.com/steinarvk/coaxy/lib/interfaces"
)

type stringValueRecord struct {
	mu          sync.Mutex
	value       string
	parsed      bool
	arrayParse  []interface{}
	objectParse map[string]interface{}
}

func (v *stringValueRecord) tryParse() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.parsed {
		return nil
	}

	var generic interface{}
	if err := json.Unmarshal([]byte(v.value), &generic); err != nil {
		// not JSON -- this isn't an error.
	} else {
		switch jsonValue := generic.(type) {
		case []interface{}:
			if len(jsonValue) > 0 {
				v.arrayParse = jsonValue
				v.parsed = true
			}

		case map[string]interface{}:
			if len(jsonValue) > 0 {
				v.objectParse = jsonValue
				v.parsed = true
			}
		}
	}

	if !v.parsed {
		m := map[string]interface{}{}
		err := logfmt.Unmarshal([]byte(v.value), logfmt.HandlerFunc(func(key, val []byte) error {
			if len(val) > 0 {
				m[string(key)] = string(val)
			}
			return nil
		}))

		if err == nil && len(m) > 0 {
			var hasNormalKey bool
			for k := range m {
				if isNormalLogfmtKey(k) {
					hasNormalKey = true
				}
			}
			if hasNormalKey {
				v.objectParse = m
				v.parsed = true
			}
		}
	}

	v.parsed = true

	return nil
}

func (v *stringValueRecord) AsValue() (string, error) {
	return v.value, nil
}

func (v *stringValueRecord) GetByIndex(index int) (interfaces.Record, error) {
	if err := v.tryParse(); err != nil {
		return nil, err
	}

	if v.arrayParse == nil {
		return nullRecord{}, nil
	}

	return accessJSONArrayByIndex(v.arrayParse, index)
}

func (v *stringValueRecord) GetByName(name string) (interfaces.Record, error) {
	if err := v.tryParse(); err != nil {
		return nil, err
	}

	if v.objectParse == nil {
		return nullRecord{}, nil
	}

	return accessJSONObjectByName(v.objectParse, name)
}

func (v *stringValueRecord) Indices() []int {
	if err := v.tryParse(); err != nil {
		return nil
	}

	if v.arrayParse == nil {
		return nil
	}

	rv := make([]int, len(v.arrayParse))
	for i := range rv {
		rv[i] = i
	}
	return rv
}

func (v *stringValueRecord) FieldNames() []string {
	if err := v.tryParse(); err != nil {
		return nil
	}

	if v.objectParse == nil {
		return nil
	}

	var rv []string

	for k := range v.objectParse {
		rv = append(rv, k)
	}

	sort.Strings(rv)

	return rv
}

var (
	normalLogfmtKeyRE = regexp.MustCompile(`^[a-zA-Z_]+$`)
)

func isNormalLogfmtKey(s string) bool {
	return normalLogfmtKeyRE.MatchString(s)
}
