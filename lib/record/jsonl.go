package record

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func makeJSONLRecordReader(r io.Reader) func() (record, error) {
	scanner := bufio.NewScanner(r)

	lineno := 0

	return func() (record, error) {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			lineno++
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			var generic interface{}

			if err := json.Unmarshal([]byte(line), &generic); err != nil {
				return nil, fmt.Errorf("failed to parse JSON on line %d: %w", lineno, err)
			}

			return jsonValueToRecord(generic)
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return nil, io.EOF
	}
}
