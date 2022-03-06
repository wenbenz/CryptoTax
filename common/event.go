package common

import "time"

const (
	// Transaction Types: buy, sell, deposit, withdrawal, mining, interest, fee
	BUY      = "buy"
	SELL     = "sell"
	DEPOSIT  = "deposit"
	WITHDRAW = "withdraw"
	MINING   = "mining"
	INTEREST = "interest"
	FEE      = "fee"
	UNKNOWN  = "unknown"

	// Currencies
	CAD  = "CAD"
	BTC  = "BTC"
	ETH  = "ETH"
	NEXO = "NEXO"
	USDC = "USDC"
	BUSD = "BUSD"
)

// Event describes a change to the balance
type Event struct {
	Time            time.Time
	Currency        string
	Amount          float64
	CadValue        float64
	Wallet          string
	TransactionType string
	Comments        string
}
