package shakepay

import (
	"strings"
	"time"

	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/CryptoTax/tools"
)

const SHAKEPAY_TIME_FORMAT = "2006-01-02T15:04:05-07"

// Headers:
// "Transaction Type","Date","Amount Debited","Debit Currency","Amount Credited","Credit Currency",
// "Buy / Sell Rate","Direction","Spot Rate","Source / Destination","Blockchain Transaction ID"
type shakepayEvent struct {
	transactionType         string
	date                    time.Time
	amountDebited           float64
	debitCurrency           string
	amountCredited          float64
	creditCurrency          string
	buySellRate             float64
	direction               string
	spotRate                float64
	address                 string
	blockchainTransactionId string
}

func (spEvent shakepayEvent) toCommonEvents() []*common.Event {
	switch spEvent.transactionType {
	case "crypto cashout":
		return []*common.Event{spEvent.toCommonDebitEvent()}
	case "purchase/sale":
		return []*common.Event{
			spEvent.toCommonDebitEvent(),
			spEvent.toCommonCreditEvent(),
		}
	case "fiat funding":
		return []*common.Event{spEvent.toCommonCreditEvent()}
	}
	return []*common.Event{}
}

func (spEvent shakepayEvent) toCommonDebitEvent() *common.Event {
	return &common.Event{
		Time:     spEvent.date,
		Currency: spEvent.debitCurrency,
		Amount:    spEvent.amountDebited,
		CadValue: spEvent.getCadValue(),
		Wallet:   spEvent.address,
		TransactionType: spEvent.getCategory(true),
		Comments: spEvent.blockchainTransactionId,
	}
}

func (spEvent shakepayEvent) toCommonCreditEvent() *common.Event {
	return &common.Event{
		Time:     spEvent.date,
		Currency: spEvent.creditCurrency,
		Amount:    spEvent.amountCredited,
		CadValue: spEvent.getCadValue(),
		Wallet:   spEvent.address,
		TransactionType: spEvent.getCategory(false),
		Comments: spEvent.blockchainTransactionId,
	}
}

func (event shakepayEvent) getCategory(isDebitEvent bool) string {
	switch event.transactionType {
	case "crypto cashout":
		return common.WITHDRAW
	case "purchase/sale":
		// if event.direction == "purchase" {
		// 	return common.BUY
		// } else {
		// 	return common.SELL
		// }
		if isDebitEvent {
			return common.SELL
		} else {
			return common.BUY
		}
	case "fiat funding":
		return common.DEPOSIT
	//TODO: implement and test other Shakepay types.
	//	currently lack sample from own records.
	default:
		return common.UNKNOWN
	}
}

func (event shakepayEvent) getCadValue() float64 {
	if event.debitCurrency == common.CAD {
		return event.amountDebited
	} else if event.creditCurrency == common.CAD {
		return event.amountCredited
	} else {
		return event.amountDebited * event.spotRate
	}
}

func parseRow(row []string) (shakepayEvent, error) {
	var err error
	event := shakepayEvent{
		transactionType:         row[0],
		debitCurrency:           strings.ToUpper(row[3]),
		creditCurrency:          strings.ToUpper(row[5]),
		direction:               row[7],
		address:                 row[9],
		blockchainTransactionId: row[10],
	}
	if event.date, err = time.Parse(SHAKEPAY_TIME_FORMAT, row[1]); err != nil {
		return shakepayEvent{}, err
	}
	if event.amountDebited, err = tools.ParseFloat (row[2]); err != nil {
		return shakepayEvent{}, err
	}
	if event.amountCredited, err = tools.ParseFloat (row[4]); err != nil {
		return shakepayEvent{}, err
	}
	if event.buySellRate, err = tools.ParseFloat (row[6]); err != nil {
		return shakepayEvent{}, err
	}
	if event.spotRate, err = tools.ParseFloat (row[8]); err != nil {
		return shakepayEvent{}, err
	}
	return event, nil
}
