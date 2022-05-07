package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/steinarvk/chaxy/lib/record"
)

func init() {
	sniffCmd := &cobra.Command{
		Use:   "sniff",
		Short: "detect the structure of data",
		RunE: func(cmd *cobra.Command, args []string) error {
			stream, err := record.OpenStream(os.Stdin)
			if err != nil {
				return err
			}

			stream.Descriptor().Show(os.Stdout)

			return nil
		},
	}

	rootCmd.AddCommand(sniffCmd)
}
