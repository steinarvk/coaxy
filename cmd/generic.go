package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/steinarvk/coaxy/lib/coaxy"
	"github.com/steinarvk/coaxy/lib/interfaces"
	"github.com/steinarvk/coaxy/lib/sniff"
)

type dataProcessorData struct {
	descriptor *sniff.Descriptor

	columnNames []string
	columnTypes []sniff.ValueType

	values <-chan []string
}

type dataProcessorCommand struct {
	processData func(context.Context, *dataProcessorData) error
}

func (d *dataProcessorCommand) registerCommonFlags(cmd *cobra.Command) {
}

func (d *dataProcessorCommand) RunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

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
	var columnNames []string

	for _, arg := range args {
		field, err := stream.ResolveField(arg)
		if err != nil {
			return err
		}

		columnNames = append(columnNames, arg)

		columns = append(columns, field)
	}

	data := &dataProcessorData{
		descriptor:  stream.Descriptor(),
		columnNames: columnNames,
	}

	reader, err := stream.Select(columns)
	if err != nil {
		return err
	}

	ch := make(chan []string, 100)

	var readErr error
	go func() {
		defer close(ch)

		for {
			values := make([]string, len(columns))
			err := reader.Read(values)
			if err == io.EOF {
				break
			}
			if err != nil {
				readErr = err
				break
			}

			ch <- values
		}
	}()

	data.values = ch

	if err := d.processData(ctx, data); err != nil {
		return err
	}

	return readErr
}
