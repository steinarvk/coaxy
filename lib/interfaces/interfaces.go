package interfaces

type Record interface {
	GetByIndex(int) (Record, error)
	GetByName(string) (Record, error)
	AsValue() (string, error)
}

type Accessor interface {
	Extract(Record) (Record, error)
}
