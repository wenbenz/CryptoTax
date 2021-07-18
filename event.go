package cryptotax

import "time"

//Categories
const (
	BUY            = "buy"
	SELL           = "sell"
	DEPOSIT        = "deposit"
	EXCHANGE_FULL  = "exchange"
	EXCHANGE_LEFT  = "exchange_left"
	EXCHANGE_RIGHT = "exchange_right"
	TRANSFER       = "transfer"
	MINING_FEE     = "mining fee"
	UNKNOWN        = "unknown"
)

//An event represents any crypto event, containing:
//	- a timestamp
//	- event type
// 	- debit action
// 	- credit action
type Event struct {
	Time     time.Time
	Category string
	Debit    Action
	Credit   Action
}
