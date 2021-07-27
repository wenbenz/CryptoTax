package common

//Currencies
const (
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
