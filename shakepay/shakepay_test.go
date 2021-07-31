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
	assert.Equal(t, common.NewAction("", common.CAD, 800., 800.), getDebit(rows[2], common.BUY))
	assert.Equal(t, common.NewAction(SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 789.4578429348551), getDebit(rows[3], common.TRANSFER))
}

func TestGetCredit(t *testing.T) {
	assert.Equal(t, common.NewAction("", common.CAD, 800., 800.), getCredit(rows[1], common.DEPOSIT))
	assert.Equal(t, common.NewAction(SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 799.999740804494), getCredit(rows[2], common.BUY))
	assert.Equal(t, common.NewAction("address1", common.BTC, 0.01777697, 789.4578429348551), getCredit(rows[3], common.TRANSFER))
}
func TestProcessLine(t *testing.T) {
	testFilePath := "testdata/in/shakepay.csv"
	expected := []*common.Event{
		{
			Time:     parseTime("2021-06-01T04:37:12+00"),
			Type:     common.DEPOSIT,
			Debit:    common.Action{},
			Credit:   common.NewAction("", common.CAD, 800., 800.),
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
		},
		{
			Time:     parseTime("2021-06-01T04:48:51+00"),
			Type:     common.BUY,
			Debit:    common.NewAction("", common.CAD, 800., 800.),
			Credit:   common.NewAction(SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 799.999740804494),
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
		},
		{
			Time:     parseTime("2021-06-01T04:52:39+00"),
			Type:     common.TRANSFER,
			Debit:    common.NewAction(SHAKEPAY_BTC_ADDR, common.BTC, 0.01777697, 789.4578429348551),
			Credit:   common.NewAction("address1", common.BTC, 0.01777697, 789.4578429348551),
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
		},
	}

	iterator, _ := NewCsvStreamReaderFromFile(testFilePath)
	for _, event := range expected {
		next, err := iterator.Next()
		assert.Nil(t, err)
		assert.Equal(t, event, next)
	}
}

func parseTime(t string) time.Time {
	ret, _ := time.Parse(SHAKEPAY_TIME_FORMAT, t)
	return ret
}
