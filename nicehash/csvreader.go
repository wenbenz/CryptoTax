package nicehash

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/go-nicehash/client"
)

type CsvStreamReader struct {
	Address  string
	Currency string
	iterator *common.CSVEventStream
}

func NewCsvStreamReaderFromCredentials(path string, purgeOlderReports bool) *CsvStreamReader {
	nhc, err := client.NewClientReadFrom("data/nicehash.credentials")
	if err != nil {
		log.Fatalln(err)
	}
	var reports []client.ReportMetadata
	reports, err = nhc.GetReportsList()
	if err != nil {
		log.Fatalln(err)
	}

	// purge older reports
	for i := 0; purgeOlderReports && i < len(reports); i += 1 {
		nhc.DeleteReport(reports[i].Id)
	}

	// create report if none exists
	if len(reports) == 0 {
		if err = nhc.CreateReport("ALL", "BTC", "CAD", "NONE", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "0", "0"); err != nil {
			log.Fatalln(err)
		}
	}

	// wait for report to be generated
	for reports[0].Status == 0 {
		reports, err = nhc.GetReportsList()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("waiting for report to generate...")
		time.Sleep(5 * time.Second)
	}
	nhs, err := nhc.GetReport(reports[0].Id)
	if err != nil {
		log.Fatalln(err)
	}
	return NewCsvStreamReader(nhs, "NICEHASH_ADDR")
}

func NewCsvStreamReader(r io.Reader, addr string) *CsvStreamReader {
	reader := CsvStreamReader{
		Address: addr,
	}
	reader.iterator = common.NewCsvIterator(r, reader.processLineStrategy)
	headers, err := reader.iterator.Stream.Read()
	if err != nil {
		return nil
	}
	currency := headers[2]
	currency = strings.TrimLeft(currency, "Amount (")
	currency = strings.TrimRight(currency, ")")
	reader.Currency = currency
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	return reader.iterator.Next()
}

func (reader *CsvStreamReader) processLineStrategy(s []string) ([]*common.Event, error) {
	nhEvent, err := parseRow(s, reader.Currency)
	if err != nil {
		return nil, err
	}
	event := nhEvent.toCommonEvent()
	event.Wallet = reader.Address
	return []*common.Event{event}, nil
}
