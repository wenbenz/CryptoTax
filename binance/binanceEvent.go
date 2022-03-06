package binance

import (
	"math"
	"strconv"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/CryptoTax/tools"
)

const BINANCE_TIME_FORMAT = "2006-01-02 15:04:05"

type binanceEvent struct {
	userId int
	time time.Time
	account string
	operation string
	coin string
	change float64
	remark string
}

func parseRow(row []string) binanceEvent {
	id, _ := strconv.Atoi(row[0])
	time, _ := time.ParseInLocation(BINANCE_TIME_FORMAT, row[1], time.UTC)
	change, _ := tools.ParseFloat(row[5])
	return binanceEvent{
		userId: id,
		time: time,
		account: row[2],
		operation: row[3],
		coin: row[4],
		change:  change,
		remark: row[6],
	}
}

func (be binanceEvent) toCommonEvent() (*common.Event, error) {
	value, err := tools.MarketDataClient{}.GetValueAtTime(be.time, be.coin)
	if err != nil {
		return nil, err
	}
	return &common.Event{
		Time: be.time,
		TransactionType: getType(be),
		Currency: be.coin,
		Amount: math.Abs(be.change),
		CadValue: value,
		Wallet: "BINANCE",
		// Comments        string
	}, nil
}

func getType(be binanceEvent) string {
	switch be.operation {
	case "Deposit": return common.DEPOSIT
	case "Large OTC trading":
		if be.change >= 0. {
			return common.BUY
		} else {
			return common.SELL
		}
	case "Withdraw": return common.WITHDRAW
	case "Fee": return common.FEE
	case "Sell":
		if be.change >= 0 {
			return common.BUY
		} else {
			return common.SELL
		}
	case "Savings Interest": return common.INTEREST
	default: return common.UNKNOWN
	}
}
