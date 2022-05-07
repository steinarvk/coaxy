package sniff

import "bytes"

func fullLines(data []byte) []string {
	buf := bytes.NewBuffer(data)

	var rv []string

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break // discard incomplete line
		}

		rv = append(rv, line)
	}

	return rv
}
