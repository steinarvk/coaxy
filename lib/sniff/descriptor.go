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
	FormatLogfmt = FormatType("logfmt")
	FormatJSONL  = FormatType("jsonl")
)

const (
	KindInt    = ValueKind("int")
	KindNumber = ValueKind("number")
	KindString = ValueKind("string")
)

type ValueType struct {
	Kind ValueKind
}

func (t ValueType) String() string {
	if t.Kind == "" {
		return "ValueType{}"
	}
	return string(t.Kind)
}

type Descriptor struct {
	Format      FormatType
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
		fmt.Fprintf(w, "  %q: %v\n", k, d.FieldTypes[k].String())
	}
}
