package shakepay

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wenbenz/CryptoTax/common"
	"github.com/wenbenz/CryptoTax/tools"
)

var rows [][]string

func TestMain(m *testing.M) {
	rows, _ = tools.ReadCsv("testdata/in/shakepay.csv")
	m.Run()
}

func TestParser(t *testing.T) {
	// "fiat funding","2021-06-01T04:37:12+00",,,800,"CAD",,"credit",,"abcd@email.com",
	spEvent, err := parseRow(rows[1])
	assert.Nil(t, err)
	assert.Equal(t, shakepayEvent{"fiat funding", parseTime("2021-06-01T04:37:12+00"), 0., "", 800., "CAD", 0., "credit", 0., "abcd@email.com", ""}, spEvent)
	
	// "purchase/sale","2021-06-01T04:48:51+00",800,"CAD",0.01777697,"BTC","45002.0302","purchase",,,
	spEvent, err = parseRow(rows[2])
	assert.Nil(t, err)
	assert.Equal(t, shakepayEvent{"purchase/sale", parseTime("2021-06-01T04:48:51+00"), 800., "CAD", 0.01777697, "BTC", 45002.0302, "purchase", 0., "", ""}, spEvent)
	
	// "crypto cashout","2021-06-01T04:52:39+00",0.01777697,"BTC",,,,"debit","44409.0215","address1","blockchainid1"
	spEvent, err = parseRow(rows[3])
	assert.Nil(t, err)
	assert.Equal(t, shakepayEvent{"crypto cashout",parseTime("2021-06-01T04:52:39+00"),0.01777697,"BTC",0.,"",0.,"debit",44409.0215,"address1","blockchainid1"}, spEvent)
}

func TestCsvReader(t *testing.T) {
	testFilePath := "testdata/in/shakepay.csv"
	expected := []*common.Event{
		{
			Time:     parseTime("2021-06-01T04:37:12+00"),
			TransactionType: common.DEPOSIT,
			Currency: common.CAD,
			Amount: 800.,
			CadValue: 800.,
			Wallet: "abcd@email.com",
		},
		{
			Time:     parseTime("2021-06-01T04:48:51+00"),
			TransactionType: common.SELL,
			Currency: common.CAD,
			Amount: 800.,
			CadValue: 800.,
		},
		{
			Time:     parseTime("2021-06-01T04:48:51+00"),
			TransactionType: common.BUY,
			Currency: common.BTC,
			Amount: 0.01777697,
			CadValue: 800.,
		},
		{
			Time:     parseTime("2021-06-01T04:52:39+00"),
			TransactionType: common.WITHDRAW,
			Wallet: "address1",
			Currency: common.BTC,
			Amount: 0.01777697,
			CadValue: 789.4578429348551,
			Comments: "blockchainid1",
		},
	}

	iterator, err := NewCsvStreamReaderFromFile(testFilePath)
	assert.Nil(t, err)
	for _, event := range expected {
		next, err := iterator.Next()
		assert.Nil(t, err)
		assert.Equal(t, event, next)
	}
	next, err := iterator.Next()
	assert.Nil(t, next)
	assert.Nil(t, err)
}

func parseTime(t string) time.Time {
	ret, _ := time.Parse(SHAKEPAY_TIME_FORMAT, t)
	return ret
}
