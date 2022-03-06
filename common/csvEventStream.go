package common

import (
	"encoding/csv"
	"errors"
	"io"
)

const Q_SIZE = 4

//CSVEventStream is an eventStream created from:
//	Stream: a csv reader
//	Strategy: a function that takes as input a CSV row and returns an Event.
type CSVEventStream struct {
	Stream   *csv.Reader
	Strategy func([]string) ([]*Event, error)
	q []*Event
	i int
	j int
}

func NewCsvIterator(stream io.Reader, strategy func([]string) ([]*Event, error)) *CSVEventStream {
	return &CSVEventStream{
		Stream:   csv.NewReader(stream),
		Strategy: strategy,
		q: make([]*Event, Q_SIZE),
		i: 0,
		j: 0,
	}
}

func (r *CSVEventStream) Next() (*Event, error) {
	if (r.i == r.j) {
		line, err := r.Stream.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return nil, err
		}
		addToQ, err := r.Strategy(line)
		if err != nil {
			return nil, err
		}
		for k, event := range addToQ {
			index := (r.j + k) % Q_SIZE
			if (k > 0 && index == r.j) {
				return nil, errors.New("too many items returned by strategy")
			}
			r.q[index] = event
		}
		r.j = (r.j + len(addToQ)) % Q_SIZE
	}
	next := r.q[r.i]
	r.i = (r.i + 1) % Q_SIZE
	return next, nil
}
