package matplotlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"text/template"

	"github.com/steinarvk/coaxy/lib/plotspec"
	"github.com/steinarvk/coaxy/lib/plotutil"
	"github.com/steinarvk/coaxy/lib/sniff"
	"github.com/steinarvk/coaxy/lib/timestamps"
)

type Options struct {
	OutputFilename string
	LogBins        bool
	WidthPixels    int
	HeightPixels   int
	DotsPerInch    int
}

var (
	tmpl = template.Must(template.New("matplotlib.prefix.py").Parse(`
import json
import matplotlib.pyplot as plt

def plot_data(data):
	xs = [x for (x, _) in data]
	ys = [y for (_, y) in data]

	kwargs = {}
	if {{ if .Options.LogBins }}True{{ else }}False{{ end}}:
		kwargs["bins"] = "log"
	
	kwargs["cmap"] = "inferno"

	plt.figure(
		figsize = ({{ .Options.WidthPixels }} / {{ .Options.DotsPerInch }},
			{{ .Options.HeightPixels }} / {{ .Options.DotsPerInch }}),
		dpi = {{ .Options.DotsPerInch }},
	)

	plt.hexbin(
		xs,
		ys,
		**kwargs
	)

	plt.savefig(
		json.loads(""" {{ .OutputFilenameJSON }} """),
		bbox_inches = "tight",
	)
`))
)

func Scatterplot(plot plotspec.Scatterplot, opts Options, w io.Writer) error {
	if opts.OutputFilename == "" {
		return errors.New("no output filename specified")
	}

	data, columnTypes, err := plotutil.SniffColumnTypes(plot.Data)
	if err != nil {
		return err
	}

	if len(columnTypes) != 2 {
		return fmt.Errorf("expected 2 columns; got %v", len(columnTypes))
	}

	transformers := []func(string) (string, error){
		nil,
		nil,
	}

	for i, ct := range columnTypes {
		if ct.Kind == sniff.KindTimestamp {
			transformers[i] = timestamps.NewNormalizerUnix()
		}
	}

	switch {
	case opts.WidthPixels == 0 && opts.HeightPixels == 0:
		opts.WidthPixels = 1024
		opts.HeightPixels = 768

	case opts.WidthPixels == 0:
		opts.WidthPixels = (opts.HeightPixels * 4) / 3

	case opts.HeightPixels == 0:
		opts.HeightPixels = (opts.WidthPixels * 3) / 4
	}

	if opts.DotsPerInch == 0 {
		opts.DotsPerInch = 80
	}

	filenameJSON, err := json.Marshal(opts.OutputFilename)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, struct {
		Options            Options
		OutputFilenameJSON string
	}{
		Options:            opts,
		OutputFilenameJSON: string(filenameJSON),
	}); err != nil {
		return err
	}

	fmt.Fprintf(w, "data = []\n")

	for tuple := range data {
		fmt.Fprintf(w, "data.append((")

		for i, v := range tuple {
			if transformers[i] != nil {
				nv, err := transformers[i](v)
				if err != nil {
					return err
				}
				v = nv
			}

			fmt.Fprintf(w, "%s,", v)
		}

		fmt.Fprintf(w, "))\n")
	}

	fmt.Fprintf(w, "plot_data(data)\n")

	return nil
}
