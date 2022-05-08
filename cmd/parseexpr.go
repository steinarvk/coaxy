package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/steinarvk/chaxy/lib/chaxyexpr"
)

func init() {
	parseexprCmd := &cobra.Command{
		Use:   "parse-expr",
		Short: "parse a chaxy expression",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument; got %d", len(args))
			}

			arg := args[0]
			result, err := chaxyexpr.Parse(arg)
			if err != nil {
				return err
			}

			fmt.Println(result)

			return nil
		},
	}

	rootCmd.AddCommand(parseexprCmd)
}
