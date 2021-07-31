package common

import (
	"encoding/csv"
	"io"
)

type CsvIterator struct {
	Stream   *csv.Reader
	Strategy func([]string) (*Event, error)
}

func NewCsvIterator(stream io.Reader, strategy func([]string) (*Event, error)) *CsvIterator {
	return &CsvIterator{
		Stream:   csv.NewReader(stream),
		Strategy: strategy,
	}
}

func (r *CsvIterator) Next() (*Event, error) {
	line, err := r.Stream.Read()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return nil, err
	}
	return r.Strategy(line)
}
