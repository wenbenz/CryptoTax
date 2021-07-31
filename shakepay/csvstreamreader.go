package shakepay

import (
	"io"
	"os"

	"github.com/wenbenz/CryptoTax/common"
)

type CsvStreamReader struct {
	source   string
	iterator *common.CsvIterator
}

func processLineStrategy(s []string) (*common.Event, error) {
	return processLine(s), nil
}

func NewCsvStreamReaderFromFile(csvFilepath string) (*CsvStreamReader, error) {
	var err error
	if file, err := os.Open(csvFilepath); err == nil {
		reader := NewCsvStreamReader(file)
		reader.source = csvFilepath
		return reader, nil
	}
	return nil, err
}

func NewCsvStreamReader(r io.Reader) *CsvStreamReader {
	reader := CsvStreamReader{
		iterator: common.NewCsvIterator(r, processLineStrategy),
	}
	reader.iterator.Stream.Read()
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	var event *common.Event
	var err error
	if event, err = reader.iterator.Next(); err == nil {
		event.Metadata[common.SOURCE] = reader.source
	}
	return event, err
}
