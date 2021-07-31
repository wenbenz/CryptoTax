package nicehash

import (
	"io"

	"github.com/wenbenz/CryptoTax/common"
)

type CsvStreamReader struct {
	iterator *common.CsvIterator
}

func processLineStrategy(s []string) (*common.Event, error) {
	return nil, nil
}

func NewCsvStreamReader(r io.Reader) *CsvStreamReader {
	reader := CsvStreamReader{
		iterator: common.NewCsvIterator(r, processLineStrategy),
	}
	reader.iterator.Stream.Read()
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	return reader.iterator.Next()
}
