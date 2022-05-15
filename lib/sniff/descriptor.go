package sniff

import (
	"fmt"
	"io"
	"sort"
)

type FormatType string
type ValueKind string

const (
	FormatJSON   = FormatType("json")
	FormatYAML   = FormatType("yaml")
	FormatCSV    = FormatType("csv")
	FormatTSV    = FormatType("tsv")
	FormatSSV    = FormatType("ssv")
	FormatShell  = FormatType("shell")
	FormatLogfmt = FormatType("logfmt")
	FormatJSONL  = FormatType("jsonl")
)

const (
	KindInt       = ValueKind("int")
	KindNumber    = ValueKind("number")
	KindString    = ValueKind("string")
	KindNull      = ValueKind("null")
	KindDate      = ValueKind("date")
	KindTimestamp = ValueKind("timestamp")
)

type ValueType struct {
	Kind     ValueKind
	Optional bool
	Format   string
}

func (t ValueType) IsNull() bool {
	return t.Kind == KindNull
}

func (t ValueType) String() string {
	if t.Kind == "" {
		return "ValueType{}"
	}

	rv := string(t.Kind)

	if t.Format != "" {
		rv = fmt.Sprintf("%s[%s]", rv, t.Format)
	}

	if t.Optional {
		rv = fmt.Sprintf("Optional[%s]", rv)
	}

	return rv
}

type Descriptor struct {
	Format      FormatType
	TupleBased  bool
	HasHeader   bool
	NumColumns  int
	ColumnNames []string
	FieldTypes  map[string]ValueType
}

func (d *Descriptor) Show(w io.Writer) {
	fmt.Fprintf(w, "%-20s %s\n", "Format:", d.Format)
	fmt.Fprintf(w, "%-20s %v\n", "Header:", d.HasHeader)

	if d.NumColumns > 0 {
		fmt.Fprintf(w, "%-20s %v\n", "Number of columns:", d.NumColumns)
	}

	if len(d.ColumnNames) > 0 {
		fmt.Fprintf(w, "%-20s %v\n", "Column names:", d.ColumnNames)
	}

	fmt.Fprintf(w, "Fields:\n")

	var fieldNames []string
	for k, _ := range d.FieldTypes {
		fieldNames = append(fieldNames, k)
	}
	sort.Strings(fieldNames)
	for _, k := range fieldNames {
		ft := d.FieldTypes[k]
		if !ft.IsNull() {
			fmt.Fprintf(w, "  %q: %v\n", k, ft.String())
		}
	}
}
