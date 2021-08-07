package nicehash

import (
	"io"
	"strconv"
	"time"

	"github.com/wenbenz/CryptoTax/common"
)

const NICEHASH_TIME_FORMAT = "2006-01-02 15:04:05 MST"

type CsvStreamReader struct {
	Address  string
	Currency string
	iterator *common.CsvIterator
}

func NewCsvStreamReader(r io.Reader, addr, currency string) *CsvStreamReader {
	reader := CsvStreamReader{
		Address:  addr,
		Currency: currency,
	}
	reader.iterator = common.NewCsvIterator(r, reader.processLineStrategy)
	reader.iterator.Stream.Read()
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

	event := common.Event{
		Time: eventTime,
		Metadata: map[string]interface{}{
			common.SOURCE: "NiceHash report" + reader.Address, //TODO: replace with report ID
		},
	}

	switch s[1] {
	case "Deposit complete", "Hashpower mining":
		event.Type = common.DEPOSIT
		event.Debit = common.Action{}
		event.Credit = common.Action{
			Address:  reader.Address,
			Currency: reader.Currency,
			Amount:   amount,
			CadValue: cadVal,
		}
	case "Withdrawal complete":
		event.Type = common.WITHDRAW
		event.Debit = common.Action{
			Address:  reader.Address,
			Currency: reader.Currency,
			Amount:   amount,
			CadValue: cadVal,
		}
		event.Credit = common.Action{}
	case "Hashpower mining fee", "Withdrawal fee":
		event.Type = common.FEE
	case "Exchange trade":
		if s[2][0] == '-' {
			event.Type = common.EXCHANGE_SELL
		} else {
			event.Type = common.EXCHANGE_BUY
		}
	default:
		event.Type = common.UNKNOWN
	}

	return &event, nil
}
