package nicehash

import (
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/go-nicehash/client"
)

const (
	NICEHASH_TIME_FORMAT = "2006-01-02 15:04:05 MST"
	PURPOSE              = "purpose"
	SOURCE               = "source"
)

type CsvStreamReader struct {
	Address  string
	Currency string
	iterator *common.CsvIterator
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

func (reader *CsvStreamReader) processLineStrategy(s []string) (*common.Event, error) {
	var err error
	var eventTime time.Time
	var amount, cadVal float64

	if eventTime, err = time.Parse(NICEHASH_TIME_FORMAT, s[0]); err != nil {
		return nil, err
	}
	if amount, err = strconv.ParseFloat(s[2], 64); err != nil {
		return nil, err
	}
	if cadVal, err = strconv.ParseFloat(s[4], 64); err != nil {
		return nil, err
	}

	purpose := s[1]
	event := common.Event{
		Time:   eventTime,
		Debit:  common.Action{},
		Credit: common.Action{},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: purpose,
		},
	}
	data := common.Action{
		Address:  reader.Address,
		Currency: reader.Currency,
		Amount:   math.Abs(amount),
		CadValue: math.Abs(cadVal),
	}

	switch purpose {
	case "Deposit complete", "Hashpower mining":
		event.Type = common.DEPOSIT
		event.Credit = data
	case "Withdrawal complete":
		// TODO fetch withdrawal destination address.
		event.Type = common.WITHDRAW
		event.Debit = data
	case "Hashpower mining fee", "Withdrawal fee", "Exchange fee":
		event.Type = common.FEE
		event.Debit = data
	case "Exchange trade":
		if s[2][0] == '-' {
			event.Type = common.EXCHANGE_SELL
			event.Debit = data
		} else {
			event.Type = common.EXCHANGE_BUY
			event.Credit = data
		}
	default:
		event.Type = common.UNKNOWN
	}

	return &event, nil
}
