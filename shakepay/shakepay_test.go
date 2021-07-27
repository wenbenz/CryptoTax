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

func TestGetCategory(t *testing.T) {
	assert.Equal(t, common.DEPOSIT, getCategory(rows[1]))
	assert.Equal(t, common.BUY, getCategory(rows[2]))
	assert.Equal(t, common.TRANSFER, getCategory(rows[3]))
	//TODO: implement and test other Shakepay types.
	//	currently lack sample from own records.
}

func TestGetTime(t *testing.T) {
	testFormat := "15:04:05 2006-01-02"
	assert.Equal(t, "04:37:12 2021-06-01", getTimestamp(rows[1]).Format(testFormat))
	assert.Equal(t, "04:48:51 2021-06-01", getTimestamp(rows[2]).Format(testFormat))
	assert.Equal(t, "04:52:39 2021-06-01", getTimestamp(rows[3]).Format(testFormat))
}

func TestGetDebit(t *testing.T) {
	assert.Equal(t, common.Action{}, getDebit(rows[1], common.DEPOSIT))
	assert.Equal(t, common.Action{"", common.CAD, 800., 800.}, getDebit(rows[2], common.BUY))
	assert.Equal(t, common.Action{SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 789.4578429348551}, getDebit(rows[3], common.TRANSFER))
}

func TestGetCredit(t *testing.T) {
	assert.Equal(t, common.Action{"", common.CAD, 800., 800.}, getCredit(rows[1], common.DEPOSIT))
	assert.Equal(t, common.Action{SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 799.999740804494}, getCredit(rows[2], common.BUY))
	assert.Equal(t, common.Action{"address1", common.BTC, 0.01777697, 789.4578429348551}, getCredit(rows[3], common.TRANSFER))
}

// "purchase/sale","",800,"CAD",0.01777697,"common.BTC","45002.0302","purchase",,,
// "crypto cashout","",0.01777697,"common.BTC",,,,"debit","44409.0215","address1","blockchainid1"
func TestGetEventsFromCSV(t *testing.T) {
	expected := []common.Event{
		{
			parseTime("2021-06-01T04:37:12+00"),
			common.DEPOSIT,
			common.Action{},
			common.Action{"", common.CAD, 800., 800.},
		},
		{
			parseTime("2021-06-01T04:48:51+00"),
			common.BUY,
			common.Action{"", common.CAD, 800., 800.},
			common.Action{SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 799.999740804494},
		},
		{
			parseTime("2021-06-01T04:52:39+00"),
			common.TRANSFER,
			common.Action{SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 789.4578429348551},
			common.Action{"address1", common.BTC, 0.01777697, 789.4578429348551},
		},
	}
	actual, err := GetEventsFromCSV("testdata/in/shakepay.csv")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func parseTime(t string) time.Time {
	ret, _ := time.Parse(SHAKEPAY_TIME_FORMAT, t)
	return ret
}
