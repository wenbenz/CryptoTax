package nicehash

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wenbenz/CryptoTax/common"
)

func TestNewCsvStreamReader(t *testing.T) {
	f, _ := os.Open("testdata/in/btc_short.csv")
	reader := NewCsvStreamReader(f, "NICEHASH_BTC_ADDR")
	assert.Equal(t, "BTC", reader.Currency)
}

func TestNext(t *testing.T) {
	f, _ := os.Open("testdata/in/btc_short.csv")
	reader := NewCsvStreamReader(f, "NICEHASH_BTC_ADDR")

	for i := 2; i <= 13; i++ {
		event, err := reader.Next()
		assert.Nil(t, err)
		assert.Equal(t, 2021, event.Time.Year())
		assert.Less(t, 0, len(event.Type))
	}

	event, err := reader.Next()
	assert.NotNil(t, err)
	assert.Nil(t, event)
}

func TestMining(t *testing.T) {
	reader := NewCsvStreamReader(
		stubReader("BTC",
			"2021-06-05 08:05:18 GMT,Hashpower mining,0.00001393,45421.21,0.63",
			"2021-06-05 08:05:18 GMT,Hashpower mining fee,-0.00000028,45421.21,-0.01"),
		"ADDR")

	expectedEvents := []*common.Event{
		{
			Time: parseTime("2021-06-05 08:05:18 GMT"),
			Type: common.DEPOSIT,
			Credit: common.Action{
				Address:  reader.Address,
				Currency: "BTC",
				Amount:   0.00001393,
				CadValue: 0.63,
			},
			Metadata: map[string]interface{}{
				SOURCE:  "Nicehash report",
				PURPOSE: "Hashpower mining",
			},
		},
		{
			Time: parseTime("2021-06-05 08:05:18 GMT"),
			Type: common.FEE,
			Debit: common.Action{
				Address:  reader.Address,
				Currency: "BTC",
				Amount:   0.00000028,
				CadValue: 0.01,
			},
			Metadata: map[string]interface{}{
				SOURCE:  "Nicehash report",
				PURPOSE: "Hashpower mining fee",
			},
		},
	}

	for _, expected := range expectedEvents {
		event, err := reader.Next()
		assert.Nil(t, err)
		assert.Equal(t, expected, event)
	}
}

func TestDeposit(t *testing.T) {
	reader := NewCsvStreamReader(
		stubReader("BTC",
			"2021-06-09 04:02:29 GMT,Deposit complete,0.00700000,39659.93,277.62"),
		"ADDR")
	event, err := reader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-06-09 04:02:29 GMT"),
		Type: common.DEPOSIT,
		Credit: common.Action{
			Address:  reader.Address,
			Currency: "BTC",
			Amount:   0.00700000,
			CadValue: 277.62,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Deposit complete",
		},
	}, event)
}

func TestWithdrawal(t *testing.T) {
	reader := NewCsvStreamReader(
		stubReader("BTC",
			"2021-07-29 17:15:47 GMT,Withdrawal complete,-0.00100014,49548.48,-49.56",
			"2021-07-29 17:15:47 GMT,Withdrawal fee,-0.00000100,49548.48,-0.05"),
		"ADDR")
	event, err := reader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-07-29 17:15:47 GMT"),
		Type: common.WITHDRAW,
		Debit: common.Action{
			Address:  reader.Address,
			Currency: "BTC",
			Amount:   0.00100014,
			CadValue: 49.56,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Withdrawal complete",
		},
	}, event)
	event, err = reader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-07-29 17:15:47 GMT"),
		Type: common.FEE,
		Debit: common.Action{
			Address:  reader.Address,
			Currency: "BTC",
			Amount:   0.00000100,
			CadValue: 0.05,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Withdrawal fee",
		},
	}, event)
}

func TestExchange(t *testing.T) {
	sellReader := NewCsvStreamReader(
		stubReader("BTC",
			"2021-06-09 16:07:36 GMT,Exchange trade,-0.00358344,39843.51,-142.78"),
		"ADDR")
	sellEvent, err := sellReader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-06-09 16:07:36 GMT"),
		Type: common.EXCHANGE_SELL,
		Debit: common.Action{
			Address:  sellReader.Address,
			Currency: "BTC",
			Amount:   0.00358344,
			CadValue: 142.78,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Exchange trade",
		},
	}, sellEvent)

	buyReader := NewCsvStreamReader(
		stubReader("NEXO",
			"2021-06-09 16:07:36 GMT,Exchange fee,-0.281849560000000000,2.53,-0.71",
			"2021-06-09 16:07:36 GMT,Exchange trade,56.369911100000000000,2.53,142.62"),
		"ADDR")
	buyFeeEvent, err := buyReader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-06-09 16:07:36 GMT"),
		Type: common.FEE,
		Debit: common.Action{
			Address:  sellReader.Address,
			Currency: "NEXO",
			Amount:   0.281849560000000000,
			CadValue: 0.71,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Exchange fee",
		},
	}, buyFeeEvent)
	buyEvent, err := buyReader.Next()
	assert.Nil(t, err)
	assert.Equal(t, &common.Event{
		Time: parseTime("2021-06-09 16:07:36 GMT"),
		Type: common.EXCHANGE_BUY,
		Credit: common.Action{
			Address:  sellReader.Address,
			Currency: "NEXO",
			Amount:   56.369911100000000000,
			CadValue: 142.62,
		},
		Metadata: map[string]interface{}{
			SOURCE:  "Nicehash report",
			PURPOSE: "Exchange trade",
		},
	}, buyEvent)
}

func stubReader(currency string, rows ...string) *strings.Reader {
	headerRow := fmt.Sprintf("Date time,Purpose,Amount (%s),* Exchange rate,Amount (CAD)", currency)
	data := append([]string{headerRow}, rows...)
	return strings.NewReader(strings.Join(data, "\n"))
}

func parseTime(t string) time.Time {
	ret, _ := time.Parse(NICEHASH_TIME_FORMAT, t)
	return ret
}
