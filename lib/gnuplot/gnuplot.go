package gnuplot

import (
	"errors"
	"fmt"
	"io"

	"github.com/steinarvk/coaxy/lib/plotspec"
	"github.com/steinarvk/coaxy/lib/plotutil"
	"github.com/steinarvk/coaxy/lib/sniff"
	"github.com/steinarvk/coaxy/lib/timestamps"
)

type Options struct {
	TerminalType   string
	OutputFilename string
	Width          int
	Height         int
}

func Scatterplot(plot plotspec.Scatterplot, opts Options, w io.Writer) error {
	if opts.OutputFilename == "" {
		return errors.New("no output filename specified")
	}

	if opts.TerminalType == "" {
		terminalType, err := TerminalTypeFromFilename(opts.OutputFilename)
		if err != nil {
			return err
		}
		opts.TerminalType = terminalType
	}

	data, columnTypes, err := plotutil.SniffColumnTypes(plot.Data)
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

	for tuple := range data {
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
