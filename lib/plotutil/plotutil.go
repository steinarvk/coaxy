package plotutil

import "github.com/steinarvk/coaxy/lib/sniff"

func limitedTupleTee(inputCh <-chan []string, limit int, fullOut, limitedOut chan<- []string) {
	for tuple := range inputCh {
		if limit > 0 {
			limitedOut <- tuple
			limit--
			if limit == 0 {
				close(limitedOut)
			}
		}

		fullOut <- tuple
	}

	if limit > 0 {
		close(limitedOut)
	}

	close(fullOut)
}

func detectColumnTypes(ch <-chan []string) ([]sniff.ValueType, error) {
	var values [][]string

	for tuple := range ch {
		if values == nil {
			values = make([][]string, len(tuple))
		}

		for i, value := range tuple {
			values[i] = append(values[i], value)
		}
	}

	rv := make([]sniff.ValueType, len(values))

	for i, columnvalues := range values {
		rv[i] = sniff.DetectValueType(columnvalues)
	}

	return rv, nil
}

func SniffColumnTypes(data <-chan []string) (<-chan []string, []sniff.ValueType, error) {
	n := 1000
	limitedChan := make(chan []string, 1000)
	unlimitedChan := make(chan []string, 1000)

	go limitedTupleTee(data, n, unlimitedChan, limitedChan)

	columnTypes, err := detectColumnTypes(limitedChan)
	if err != nil {
		return nil, nil, err
	}

	return unlimitedChan, columnTypes, nil
}
