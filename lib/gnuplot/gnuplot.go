package gnuplot

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/steinarvk/coaxy/lib/plotspec"
)

type Options struct {
	TerminalType   string
	OutputFilename string
}

func Scatterplot(plot plotspec.Scatterplot, opts Options, w io.Writer) error {
	if opts.TerminalType == "" {
		return errors.New("no output mode specified")
	}

	if opts.OutputFilename == "" {
		return errors.New("no output filename specified")
	}

	fmt.Fprintf(w, "$data << END_OF_DATA\n")

	for tuple := range plot.Data {
		fmt.Fprintf(w, strings.Join(tuple, " ")+"\n")
	}

	fmt.Fprintf(w, "END_OF_DATA\n")

	fmt.Fprintf(w, "set terminal %s\n", opts.TerminalType)
	fmt.Fprintf(w, "set output %q\n", opts.OutputFilename)

	fmt.Fprintf(w, "plot $data using 1:2 with points")

	return nil
}
