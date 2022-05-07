package plotspec

type Scatterplot struct {
	ColumnNames []string
	Data        <-chan []string
}
