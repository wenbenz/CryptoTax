package shakepay

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
)

const SOURCE = "source"

type CsvStreamReader struct {
	BtcAddr  string
	EthAddr  string
	source   string
	iterator *common.CsvIterator
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
	reader.iterator.Stream.Read()
	return &reader
}

func (reader *CsvStreamReader) Next() (*common.Event, error) {
	return reader.iterator.Next()
}

func (r *CsvStreamReader) processLineStrategy(line []string) (*common.Event, error) {
	category := getCategory(line)
	timestamp, err := getTimestamp(line)
	if err != nil {
		return nil, err
	}
	debit, err := getDebit(line, category)
	if err != nil {
		return nil, err
	}
	credit, err := getCredit(line, category)
	if err != nil {
		return nil, err
	}
	event := common.Event{
		Time:   timestamp,
		Type:   category,
		Debit:  debit,
		Credit: credit,
		Metadata: map[string]interface{}{
			SOURCE: r.source,
		},
	}
	externalAddr := line[9]
	switch category {
	case common.DEPOSIT:
		event.Debit.Address = externalAddr
	case common.WITHDRAW:
		switch event.Credit.Currency {
		case common.BTC:
			event.Debit.Address = r.BtcAddr
		case common.ETH:
			event.Debit.Address = r.EthAddr
		default:
			return nil, fmt.Errorf("unknown Shakepay currency '%s'", event.Credit.Currency)
		}
		event.Credit.Address = externalAddr
	case common.BUY:
		switch event.Credit.Currency {
		case common.BTC:
			event.Credit.Address = r.BtcAddr
		case common.ETH:
			event.Credit.Address = r.EthAddr
		default:
			return nil, fmt.Errorf("unknown Shakepay currency '%s'", event.Credit.Currency)
		}
	case common.SELL:
		switch event.Debit.Currency {
		case common.BTC:
			event.Debit.Address = r.BtcAddr
		case common.ETH:
			event.Debit.Address = r.EthAddr
		default:
			return nil, fmt.Errorf("unknown Shakepay currency '%s'", event.Debit.Currency)
		}
	default:
		return nil, fmt.Errorf("unhandled Shakepay transaction category '%s'", category)
	}
	return &event, nil
}

//TODO: Missing functionality:
//	- common.TRANSFER crypto into shakepay wallet
//	- exchange to fiat
//	- fiat withdrawal
func getDebit(row []string, category string) (common.Action, error) {
	if category == common.DEPOSIT {
		return common.Action{}, nil
	}
	amount, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return common.Action{}, err
	}
	currency := strings.ToUpper(row[3])
	rate, err := getRate(row, currency)
	if err != nil {
		return common.Action{}, err
	}
	return common.Action{Currency: currency, Amount: amount, CadValue: amount * rate}, nil
}

//TODO: Missing functionality:
//	- common.TRANSFER crypto into shakepay wallet
//	- exchange to fiat
//	- fiat withdrawal
func getCredit(row []string, category string) (common.Action, error) {
	//TODO: handle this error
	if category == common.WITHDRAW {
		return getDebit(row, category)
	}
	amount, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return common.Action{}, err
	}
	currency := strings.ToUpper(row[5])
	rate, err := getRate(row, currency)
	if err != nil {
		return common.Action{}, err
	}
	return common.Action{Currency: currency, Amount: amount, CadValue: amount * rate}, nil
}

const SHAKEPAY_TIME_FORMAT = "2006-01-02T15:04:05-07"

func getCategory(row []string) string {
	cat := row[0]
	if cat == "crypto cashout" {
		return common.WITHDRAW
	} else if cat == "purchase/sale" {
		if row[7] == "purchase" {
			return common.BUY
		} else {
			return common.SELL
		}
	} else if cat == "fiat funding" {
		return common.DEPOSIT
	}

	//TODO: implement and test other Shakepay types.
	//	currently lack sample from own records.
	return common.UNKNOWN
}

func getTimestamp(row []string) (time.Time, error) {
	timestamp := row[1]
	return time.Parse(SHAKEPAY_TIME_FORMAT, timestamp)
}

func getRate(row []string, currency string) (float64, error) {
	if currency == common.CAD {
		return 1., nil
	}
	return strconv.ParseFloat("0"+row[6]+row[8], 64)
}
