package shakepay

import "github.com/wenbenz/CryptoTax/common"

type CSVReader struct {
	path string
}

func NewCSVReader(path string) *CSVReader {
	return &CSVReader{
		path: path,
	}
}

func (r *CSVReader) ReadAll() ([]common.Event, error) {
	return GetEventsFromCSV(r.path)
}
