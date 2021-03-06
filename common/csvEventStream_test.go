package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvNext(t *testing.T) {
	stringReader := strings.NewReader("Hello,World\nWhat's,up?\n")
	reader := NewCsvIterator(stringReader, func(s []string) ([]*Event, error) {
		return []*Event{{
			TransactionType: strings.Join(s, " -> "),
		}}, nil
	})
	assertNextIs(t, reader, "Hello -> World")
	assertNextIs(t, reader, "What's -> up?")
	assertNextIsNil(t, reader)
}

func TestLineEndsWithEOF(t *testing.T) {
	stringReader := strings.NewReader("Hello,World")
	reader := NewCsvIterator(stringReader, func(s []string) ([]*Event, error) {
		return []*Event{{
			TransactionType: strings.Join(s, " -> "),
		}}, nil
	})
	assertNextIs(t, reader, "Hello -> World")
	assertNextIsNil(t, reader)
}

func assertNextIs(t *testing.T, r *CSVEventStream, s string) {
	event, err := r.Next()
	assert.Nil(t, err)
	assert.Equal(t, s, event.TransactionType)
}

func assertNextIsNil(t *testing.T, r *CSVEventStream) {
	event, err := r.Next()
	assert.Nil(t, err)
	assert.Nil(t, event)
}
