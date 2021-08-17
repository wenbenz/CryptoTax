package common

import "time"

const (
	// Transaction Types
	BUY           = "buy"
	SELL          = "sell"
	DEPOSIT       = "deposit"
	WITHDRAW      = "withdraw"
	EXCHANGE_FULL = "exchange"
	EXCHANGE_SELL = "exchange_sell"
	EXCHANGE_BUY  = "exchange_buy"
	FEE           = "fee"
	UNKNOWN       = "unknown"

	//Currencies
	CAD  = "CAD"
	BTC  = "BTC"
	ETH  = "ETH"
	NEXO = "NEXO"
	USDC = "USDC"
	BUSD = "BUSD"
)

//An action is an operation on an address and consists of
type Action struct {
	Address  string
	Currency string
	Amount   float64
	CadValue float64
}

//An event represents any crypto event, containing:
//	- a timestamp
//	- event type
// 	- debit action
// 	- credit action
type Event struct {
	Time     time.Time
	Type     string
	Debit    Action
	Credit   Action
	Metadata map[string]interface{}
}

//EventStreamReader defines a way to retrieve events from a stream such as a response body.
//Next returns the next event in the stream. Implementations should return events in chronological order.
type EventIterator interface {
	Next() (*Event, error)
}
