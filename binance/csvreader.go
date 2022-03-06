package binance

import (
	"io"

	"github.com/wenbenz/CryptoTax/common"
)

type CsvStreamReader struct {
	iterator *common.CSVEventStream
}

func NewCsvStreamReader(r io.Reader) *CsvStreamReader {
	reader := CsvStreamReader{}
	reader.iterator = common.NewCsvIterator(r, reader.processLineStrategy)
	reader.iterator.Stream.Read()	// discard the headers
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	return reader.iterator.Next()
}

func (reader *CsvStreamReader) processLineStrategy(s []string) ([]*common.Event, error) {
	binanceEvent := parseRow(s)
	event, err := binanceEvent.toCommonEvent()
	if err != nil {
		return nil, err
	}
	return []*common.Event{event}, nil
}
