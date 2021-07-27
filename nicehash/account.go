package nicehash

type Account struct {
	Active         bool                  `json:"active"`
	Currency       string                `json:"currency"`
	TotalBalance   string                `json:"totalBalance"`
	Available      string                `json:"available"`
	Pending        string                `json:"pending"`
	PendingDetails AccountExtendedDetail `json:"pendingDetails"`
	BtcRate        float64               `json:"btcRate"`
}

type AccountExtendedDetail struct {
	Deposit         string `json:"deposit"`
	Withdrawal      string `json:"withdrawal"`
	Exchange        string `json:"exchange"`
	HashpowerOrders string `json:"hashpowerOrders"`
	UnpaidMining    string `json:"unpaidMining"`
}

func (c *Client) GetAccount(currency string, extended bool) (Account, error) {
	var account Account
	err := c.Do("GET", "/main/api/v2/accounting/account2/BTC", map[string]string{"extendedResponse": "true"}, &account)
	return account, err
}
