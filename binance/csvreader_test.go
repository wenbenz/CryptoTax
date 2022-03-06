package binance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenbenz/CryptoTax/common"
)

func TestReader(t *testing.T) {
	f, _ := os.Open("testdata/in/binance.csv")
	stream := NewCsvStreamReader(f)

	// 163059421,2021-09-10 05:31:08,Spot,Deposit,BTC,0.04195844,""
	event, err := stream.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{}, event)
	// 163059421,2021-09-10 14:56:31,Spot,Large OTC trading,USDC,3794.50675184,""
	// 163059421,2021-09-10 14:56:31,Spot,Large OTC trading,BTC,-0.08392160,""
	// 163059421,2021-09-10 15:12:43,Spot,Withdraw,USDC,-3794.50675100,Withdraw fee is included
	// 163059421,2021-09-27 02:59:07,Spot,Deposit,ETH,0.62316269,""
	// 163059421,2021-09-27 03:02:20,Spot,Fee,BUSD,-0.38564198,""
	// 163059421,2021-09-27 03:02:20,Spot,Sell,ETH,-0.01000000,""
	// 163059421,2021-09-27 03:02:20,Spot,Fee,BUSD,-1.32115200,""
	// 163059421,2021-09-27 04:46:41,Spot,Withdraw,BUSD,-1958.05405892,Withdraw fee is included
	// 163059421,2021-11-05 03:14:22,Spot,Savings Interest,BUSD,0.02192000,""

}
