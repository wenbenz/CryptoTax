package nicehash

import (
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
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
