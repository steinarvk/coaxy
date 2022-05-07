package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steinarvk/chaxy/lib/gnuplot"
	"github.com/steinarvk/chaxy/lib/plotspec"
	"github.com/steinarvk/chaxy/lib/record"
)

func init() {
	var flagOutputFilename string

	scatterCmd := &cobra.Command{
		Use:   "scatter [FIELD-X] [FIELD-Y]",
		Short: "Generate a scatterplot",
		RunE: func(cmd *cobra.Command, args []string) error {
			stream, err := record.OpenStream(os.Stdin)
			if err != nil {
				return err
			}

			if len(args) != 2 {
				return fmt.Errorf("must specify exactly two fields; got %v", args)
			}

			var columns []*record.Accessor

			for _, arg := range args {
				field, err := stream.ResolveField(arg)
				if err != nil {
					return err
				}

				columns = append(columns, field)
			}

			reader, err := stream.Select(columns)
			if err != nil {
				return err
			}

			dataChan := make(chan []string, 100)

			var readErr error
			go func() {
				readErr = reader.ForEach(func(tuple []string) error {
					dataChan <- tuple
					return nil
				})
				close(dataChan)
			}()

			if flagOutputFilename == "" {
				flagOutputFilename = "scatterplot.generated.png"
				fmt.Fprintf(os.Stderr, "warning: --output not specified; writing to %q.\n", flagOutputFilename)
			}

			terminalType, err := gnuplot.TerminalTypeFromFilename(flagOutputFilename)
			if err != nil {
				return err
			}

			spec := plotspec.Scatterplot{
				Data: dataChan,
			}
			options := gnuplot.Options{
				TerminalType:   terminalType,
				OutputFilename: flagOutputFilename,
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

			if readErr != nil {
				return readErr
			}

			if scriptErr != nil {
				return scriptErr
			}

			return nil
		},
	}

	scatterCmd.Flags().StringVar(&flagOutputFilename, "output", "", "output file to generate")

	rootCmd.AddCommand(scatterCmd)
}
