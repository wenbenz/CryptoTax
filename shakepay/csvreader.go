package shakepay

import (
	"io"
	"os"

	"github.com/wenbenz/CryptoTax/common"
)

const SOURCE = "source"

type CsvStreamReader struct {
	BtcAddr  string
	EthAddr  string
	source   string
	iterator *common.CSVEventStream
}

func NewCsvStreamReaderFromFile(csvFilepath string) (*CsvStreamReader, error) {
	var err error
	if file, err := os.Open(csvFilepath); err == nil {
		reader := NewCsvStreamReader(file, "SHAKEPAY_BTC_ADDR", "SHAKEPAY_ETH_ADDR")
		reader.source = csvFilepath
		return reader, nil
	}
	return nil, err
}

func NewCsvStreamReader(r io.Reader, btcAddr, ethAddr string) *CsvStreamReader {
	reader := CsvStreamReader{
		BtcAddr: btcAddr,
		EthAddr: ethAddr,
	}
	reader.iterator = common.NewCsvIterator(r, reader.processLineStrategy)
	reader.iterator.Stream.Read()	// discard headers
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	return reader.iterator.Next()
}

func (r *CsvStreamReader) processLineStrategy(line []string) ([]*common.Event, error) {
	spEvent, err := parseRow(line);
	if err != nil {
		return nil, err
	}
	return spEvent.toCommonEvents(), nil
}
