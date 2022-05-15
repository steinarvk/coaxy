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

	flagInputFilename string
}

func (d *dataProcessorCommand) registerCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&d.flagInputFilename, "input", "", "input filename (if not stdin)")
}

func (d *dataProcessorCommand) RunE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	var inputReader io.Reader

	if d.flagInputFilename == "" && len(args) > 0 {
		if fileExists(args[len(args)-1]) {
			d.flagInputFilename = args[len(args)-1]
			args = args[:len(args)-1]
		} else if fileExists(args[0]) {
			d.flagInputFilename = args[0]
			args = args[1:]
		}
	}

	if d.flagInputFilename != "" {
		f, err := os.Open(d.flagInputFilename)
		if err != nil {
			return fmt.Errorf("error opening %q: %w", d.flagInputFilename, err)
		}

		inputReader = f

		defer f.Close()
	} else {
		inputReader = os.Stdin
	}

	stream, err := coaxy.OpenStream(inputReader)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		n := stream.Descriptor().NumColumns
		if n > 1 {
			for i := 0; i < n; i++ {
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
