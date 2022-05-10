package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steinarvk/coaxy/lib/coaxy"
	"github.com/steinarvk/coaxy/lib/interfaces"
)

func init() {
	normalizeCmd := &cobra.Command{
		Use:   "normalize [FIELDS...]",
		Short: "normalize data to a TSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			stream, err := coaxy.OpenStream(os.Stdin)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				n := stream.Descriptor().NumColumns
				if n > 1 {
					for i := 1; i <= n; i++ {
						args = append(args, fmt.Sprintf("%d", i))
					}
				}
			}

			if len(args) == 0 {
				return errors.New("no fields specified")
			}

			var columns []interfaces.Accessor

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

			values := make([]string, len(columns))
			for {
				err := reader.Read(values)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				// TODO: should normalize to TSV _with_ headers if they exist
				// TODO: should properly handle quoting multiline strings etc
				fmt.Println(strings.Join(values, "\t"))
			}

			return nil
		},
	}

	rootCmd.AddCommand(normalizeCmd)
}
