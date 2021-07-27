package common

type CSVReader interface {
	ReadAll() ([]Event, error)
}
