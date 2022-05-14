package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	extractCmd := &cobra.Command{
		Use:   "extract [FIELDS...]",
		Short: "extract columns from file",
	}

	dataproc := &dataProcessorCommand{
		processData: func(ctx context.Context, data *dataProcessorData) error {
			// TODO: should normalize to TSV _with_ headers if they exist
			// TODO: should properly handle quoting multiline strings etc
			for tuple := range data.values {
				fmt.Println(strings.Join(tuple, "\t"))
			}

			return nil
		},
	}
	dataproc.registerCommonFlags(extractCmd)
	extractCmd.RunE = dataproc.RunE

	rootCmd.AddCommand(extractCmd)
}
