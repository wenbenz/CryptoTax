package common

import "time"

//Categories
const (
	BUY           = "buy"
	SELL          = "sell"
	DEPOSIT       = "deposit"
	EXCHANGE_FULL = "exchange"
	EXCHANGE_SELL = "exchange_sell"
	EXCHANGE_BUY  = "exchange_buy"
	TRANSFER      = "transfer"
	MINING_FEE    = "mining fee"
	UNKNOWN       = "unknown"
)

//An event represents any crypto event, containing:
//	- a timestamp
//	- event type
// 	- debit action
// 	- credit action
type Event struct {
	Time   time.Time
	Type   string
	Debit  Action
	Credit Action
}
