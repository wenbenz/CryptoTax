package nicehash

import (
	"math"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/CryptoTax/tools"
)

const (
	NICEHASH_TIME_FORMAT = "2006-01-02 15:04:05 MST"
)

type nicehashEvent struct {
	time time.Time
	purpose string
	coinType string
	coins float64
	exchangeRate float64
	cad float64
}

func (nhEvent nicehashEvent) toCommonEvent() *common.Event {
	return &common.Event{
		Time: nhEvent.time,
		Currency: nhEvent.coinType,
		Amount: math.Abs(nhEvent.coins),
		CadValue: math.Abs(nhEvent.cad),
		// Wallet          string
		TransactionType: nhEvent.getType(),
		// Comments        string
	}
}

func parseRow(row []string, coinType string) (*nicehashEvent, error) {
	nhEvent := nicehashEvent{
		purpose: row[1],
		coinType: coinType,
	}
	var err error
	if nhEvent.time, err = time.Parse(NICEHASH_TIME_FORMAT, row[0]); err != nil {
		return nil, err
	}
	if nhEvent.coins, err = tools.ParseFloat(row[2]); err != nil {
		return nil, err
	}
	if nhEvent.exchangeRate, err = tools.ParseFloat(row[3]); err != nil {
		return nil, err
	}
	if nhEvent.cad, err = tools.ParseFloat(row[4]); err != nil {
		return nil, err
	}
	return &nhEvent, nil
}

func (event nicehashEvent) getType() string {
	switch event.purpose {
	case "Hashpower mining":
		return common.MINING
	case "Deposit complete":
		return common.DEPOSIT
	case "Withdrawal complete":
		return common.WITHDRAW
	case "Hashpower mining fee", "Withdrawal fee", "Exchange fee":
		return common.FEE
	case "Exchange trade":
		if event.coins < 0 {
			return common.SELL
		} else {
			return common.BUY
		}
	default:
		return common.UNKNOWN
	}
}