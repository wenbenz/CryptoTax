package shakepay

import (
	"strconv"
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/CryptoTax/tools"
)

//TODO: parameterize this.
const (
	SHAKEPAY_BTC_ADDR = "BTC_ADDR"
	SHAKEPAY_ETH_ADDR = "ETH_ADDR"
)

const SHAKEPAY_TIME_FORMAT = "2006-01-02T15:04:05-07"

func GetEventsFromCSV(path string) ([]common.Event, error) {
	lines, err := tools.ReadCsv(path)
	if err != nil {
		return nil, err
	}

	// parse lines
	var EventsToReturn []common.Event
	for _, line := range lines[1:] {
		category := getCategory(line)
		EventsToReturn = append(EventsToReturn, common.Event{
			Time:   getTimestamp(line),
			Type:   category,
			Debit:  getDebit(line, category),
			Credit: getCredit(line, category),
		})
	}
	return EventsToReturn, nil
}

func getCategory(row []string) string {
	cat := row[0]
	if cat == "crypto cashout" {
		return common.TRANSFER
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

func getTimestamp(row []string) time.Time {
	timestamp := row[1]
	//TODO: handle this error
	ret, _ := time.Parse(SHAKEPAY_TIME_FORMAT, timestamp)
	return ret
}

//TODO: Missing functionality:
//	- common.TRANSFER crypto into shakepay wallet
//	- exchange to fiat
//	- fiat withdrawal
func getDebit(row []string, category string) common.Action {
	//TODO: handle this error
	amount, _ := strconv.ParseFloat(row[2], 64)
	currency := strings.ToUpper(row[3])
	var address string
	if category == common.TRANSFER {
		if currency == common.BTC {
			address = SHAKEPAY_BTC_ADDR
		} else if currency == common.ETH {
			address = SHAKEPAY_ETH_ADDR
		}
	}
	rate := getRate(row, currency)
	return common.Action{
		address,
		currency,
		amount,
		amount * rate,
	}
}

//TODO: Missing functionality:
//	- common.TRANSFER crypto into shakepay wallet
//	- exchange to fiat
//	- fiat withdrawal
func getCredit(row []string, category string) common.Action {
	//TODO: handle this error
	amount, _ := strconv.ParseFloat(row[4], 64)
	currency := strings.ToUpper(row[5])
	rate := getRate(row, currency)
	var address string
	if category == common.TRANSFER {
		action := getDebit(row, category)
		action.Address = row[9]
		return action
	} else if category == common.BUY {
		if currency == common.BTC {
			address = SHAKEPAY_BTC_ADDR
		} else if currency == common.ETH {
			address = SHAKEPAY_ETH_ADDR
		}
	}
	return common.Action{
		address,
		currency,
		amount,
		amount * rate,
	}
}

func getRate(row []string, currency string) float64 {
	//TODO: handle this error
	if currency == common.CAD {
		return 1.
	}
	rate, _ := strconv.ParseFloat("0"+row[6]+row[8], 64)
	return rate
}
