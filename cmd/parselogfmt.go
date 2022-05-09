package cmd

import (
	"fmt"

	"github.com/kr/logfmt"
	"github.com/spf13/cobra"
)

func init() {
	parselogfmtCmd := &cobra.Command{
		Use:   "parse-logfmt",
		Short: "parse a coaxy logfmt",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument; got %d", len(args))
			}

			arg := args[0]

			err := logfmt.Unmarshal([]byte(arg), logfmt.HandlerFunc(func(key, val []byte) error {
				fmt.Printf("%q: %q\n", string(key), string(val))
				return nil
			}))
			return err
		},
	}

	rootCmd.AddCommand(parselogfmtCmd)
}
