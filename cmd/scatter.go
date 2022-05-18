package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/steinarvk/coaxy/lib/gnuplot"
	"github.com/steinarvk/coaxy/lib/matplotlib"
	"github.com/steinarvk/coaxy/lib/plotspec"
)

func init() {
	var flagOutputFilename string
	var flagShowScript bool
	var flagLogBins bool
	var flagWidth int
	var flagHeight int
	var flagEngine string

	scatterCmd := &cobra.Command{
		Use:   "scatter [FIELD-X] [FIELD-Y]",
		Short: "Generate a scatterplot",
	}

	genericPlotter := func(prep func(io.Writer) error, runner func(func(io.Writer) error) error) error {
		if flagShowScript {
			return prep(os.Stdout)
		}

		return runner(prep)
	}

	plotWithGnuplot := func(spec plotspec.Scatterplot) error {
		options := gnuplot.Options{
			OutputFilename: flagOutputFilename,
			Width:          flagWidth,
			Height:         flagHeight,
		}

		return genericPlotter(
			func(w io.Writer) error {
				return gnuplot.Scatterplot(spec, options, w)
			},
			gnuplot.RunSubprocess,
		)
	}

	plotWithMatplotlib := func(spec plotspec.Scatterplot) error {
		options := matplotlib.Options{
			OutputFilename: flagOutputFilename,
			LogBins:        flagLogBins,
		}

		return genericPlotter(
			func(w io.Writer) error {
				return matplotlib.Scatterplot(spec, options, w)
			},
			matplotlib.RunSubprocess,
		)
	}

	dataproc := &dataProcessorCommand{
		processData: func(ctx context.Context, data *dataProcessorData) error {
			if len(data.columnNames) != 2 {
				return fmt.Errorf("expected 2 columns; got %v", len(data.columnNames))
			}

			if flagOutputFilename == "" {
				flagOutputFilename = "scatterplot.generated.png"
				fmt.Fprintf(os.Stderr, "warning: --output not specified; writing to %q.\n", flagOutputFilename)
			}

			spec := plotspec.Scatterplot{
				Data: data.values,
			}

			var err error

			switch flagEngine {
			case "gnuplot":
				err = plotWithGnuplot(spec)

			case "matplotlib":
				err = plotWithMatplotlib(spec)

			default:
				return fmt.Errorf("no such plotting engine %q", flagEngine)
			}

			return err
		},
	}
	dataproc.registerCommonFlags(scatterCmd)
	scatterCmd.RunE = dataproc.RunE

	scatterCmd.Flags().StringVar(&flagOutputFilename, "output", "", "output file to generate")
	scatterCmd.Flags().StringVar(&flagEngine, "engine", "gnuplot", "engine to use for plotting")
	scatterCmd.Flags().BoolVar(&flagShowScript, "show-script", false, "show raw script")
	scatterCmd.Flags().BoolVar(&flagLogBins, "log-bins", false, "enable logarithmic binning for binned graphs")
	scatterCmd.Flags().IntVar(&flagWidth, "width", 0, "width (pixels) of output graphic")
	scatterCmd.Flags().IntVar(&flagHeight, "height", 0, "height (pixels) of output graphic")

	rootCmd.AddCommand(scatterCmd)
}
