package sniff

import (
	"bytes"
	"errors"
)

// Things to detect:
//   - CSV/TSV/SSV with header
//   - CSV/TSV/SSV without header
//   - JSON
//   - YAML
//   - logfmt

var (
	errAutodetectFailed = errors.New("failed to autodetect format")
)

func Sniff(data []byte) (*Descriptor, error) {
	lines := fullLines(data)
	if len(lines) > 0 {
		desc, err := sniffLines(lines)
		if err != nil {
			return nil, err
		}
		if desc != nil {
			return desc, nil
		}
	}

	desc, err := jsonArrayParser(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if desc != nil {
		return desc, nil
	}

	return nil, errAutodetectFailed
}
