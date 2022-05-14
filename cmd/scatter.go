package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steinarvk/coaxy/lib/gnuplot"
	"github.com/steinarvk/coaxy/lib/plotspec"
)

func init() {
	var flagOutputFilename string
	var flagShowScript bool

	scatterCmd := &cobra.Command{
		Use:   "scatter [FIELD-X] [FIELD-Y]",
		Short: "Generate a scatterplot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
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

			terminalType, err := gnuplot.TerminalTypeFromFilename(flagOutputFilename)
			if err != nil {
				return err
			}

			spec := plotspec.Scatterplot{
				Data: data.values,
			}
			options := gnuplot.Options{
				TerminalType:   terminalType,
				OutputFilename: flagOutputFilename,
			}

			if flagShowScript {
				return gnuplot.Scatterplot(spec, options, os.Stdout)
			}

			runner := gnuplot.Subprocess()

			var scriptErr error
			go func() {
				scriptErr = gnuplot.Scatterplot(spec, options, runner.Stdin)
				runner.Stdin.Close()
			}()

			if err := runner.Run(); err != nil {
				return err
			}

			if scriptErr != nil {
				return scriptErr
			}

			return nil
		},
	}
	dataproc.registerCommonFlags(scatterCmd)
	scatterCmd.RunE = dataproc.RunE

	scatterCmd.Flags().StringVar(&flagOutputFilename, "output", "", "output file to generate")
	scatterCmd.Flags().BoolVar(&flagShowScript, "show-script", false, "show raw script")

	rootCmd.AddCommand(scatterCmd)
}
