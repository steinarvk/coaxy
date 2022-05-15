package gnuplot

import (
	"errors"
	"fmt"
	"io"

	"github.com/steinarvk/coaxy/lib/plotspec"
	"github.com/steinarvk/coaxy/lib/sniff"
	"github.com/steinarvk/coaxy/lib/timestamps"
)

type Options struct {
	TerminalType   string
	OutputFilename string
	Width          int
	Height         int
}

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

func Scatterplot(plot plotspec.Scatterplot, opts Options, w io.Writer) error {
	if opts.TerminalType == "" {
		return errors.New("no output mode specified")
	}

	if opts.OutputFilename == "" {
		return errors.New("no output filename specified")
	}

	n := 1000
	limitedChan := make(chan []string, 1000)
	unlimitedChan := make(chan []string, 1000)

	go limitedTupleTee(plot.Data, n, unlimitedChan, limitedChan)

	columnTypes, err := detectColumnTypes(limitedChan)
	if err != nil {
		return err
	}

	if len(columnTypes) != 2 {
		return fmt.Errorf("expected 2 columns; got %v", len(columnTypes))
	}

	letters := []string{"x", "y"}

	transformers := []func(string) (string, error){
		nil,
		nil,
	}

	for i, letter := range letters {
		if columnTypes[i].Kind == sniff.KindDate {
			fmt.Fprintf(w, "set %sdata time\n", letter)
			fmt.Fprintf(w, "set timefmt %q\n", "%Y-%m-%d")
		}
		if columnTypes[i].Kind == sniff.KindTimestamp {
			transformers[i] = timestamps.NewNormalizerISO()
			fmt.Fprintf(w, "set %sdata time\n", letter)
			fmt.Fprintf(w, "set timefmt %q\n", "%Y-%m-%dT%H:%M:%S")
		}
	}

	fmt.Fprintf(w, "$data << END_OF_DATA\n")

	for tuple := range unlimitedChan {
		for i, v := range tuple {
			if transformers[i] != nil {
				nv, err := transformers[i](v)
				if err != nil {
					return err
				}
				v = nv
			}

			if (i + 1) == len(tuple) {
				fmt.Fprintf(w, "%s\n", v)
			} else {
				fmt.Fprintf(w, "%s ", v)
			}
		}
	}

	fmt.Fprintf(w, "END_OF_DATA\n")

	if opts.Width != 0 || opts.Height != 0 {
		width := opts.Width
		height := opts.Height

		if width == 0 {
			width = height
		}

		fmt.Fprintf(w, "set terminal %s size %d, %d\n", opts.TerminalType, width, height)
	} else {
		fmt.Fprintf(w, "set terminal %s\n", opts.TerminalType)
	}
	fmt.Fprintf(w, "set output %q\n", opts.OutputFilename)

	fmt.Fprintf(w, "plot $data using 1:2 with points")

	return nil
}
