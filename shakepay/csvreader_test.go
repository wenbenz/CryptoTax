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
	assert.Equal(t, common.WITHDRAW, getCategory(rows[3]))
	//TODO: implement and test other Shakepay types.
	//	currently lack sample from own records.
}

func TestGetTime(t *testing.T) {
	testFormat := "15:04:05 2006-01-02"
	expectedValues := []string{"04:37:12 2021-06-01", "04:48:51 2021-06-01", "04:52:39 2021-06-01"}
	for i, expected := range expectedValues {
		actual, err := getTimestamp(rows[i+1])
		assert.Nil(t, err)
		assert.Equal(t, expected, actual.Format(testFormat))
	}
}

func TestGetDebit(t *testing.T) {
	expectedActions := []common.Action{
		{},
		{Currency: common.CAD, Amount: 800., CadValue: 800.},
		{Currency: common.BTC, Amount: 0.01777697, CadValue: 789.4578429348551},
	}
	for i, expected := range expectedActions {
		actual, err := getDebit(rows[i+1], getCategory(rows[i+1]))
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestGetCredit(t *testing.T) {
	expectedActions := []common.Action{
		{Currency: common.CAD, Amount: 800., CadValue: 800.},
		{Currency: common.BTC, Amount: 0.01777697, CadValue: 799.999740804494},
		{Currency: common.BTC, Amount: 0.01777697, CadValue: 789.4578429348551},
	}
	for i, expected := range expectedActions {
		actual, err := getCredit(rows[i+1], getCategory(rows[i+1]))
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestCsvReader(t *testing.T) {
	testFilePath := "testdata/in/shakepay.csv"
	expected := []*common.Event{
		{
			Time:     parseTime("2021-06-01T04:37:12+00"),
			Type:     common.DEPOSIT,
			Debit:    common.Action{Address: "abcd@email.com"},
			Credit:   common.Action{Currency: common.CAD, Amount: 800., CadValue: 800.},
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
		},
		{
			Time:     parseTime("2021-06-01T04:48:51+00"),
			Type:     common.BUY,
			Debit:    common.Action{Currency: common.CAD, Amount: 800., CadValue: 800.},
			Credit:   common.Action{Address: "SHAKEPAY_BTC_ADDR", Currency: common.BTC, Amount: 0.01777697, CadValue: 799.999740804494},
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
		},
		{
			Time:     parseTime("2021-06-01T04:52:39+00"),
			Type:     common.WITHDRAW,
			Debit:    common.Action{Address: "SHAKEPAY_BTC_ADDR", Currency: common.BTC, Amount: 0.01777697, CadValue: 789.4578429348551},
			Credit:   common.Action{Address: "address1", Currency: common.BTC, Amount: 0.01777697, CadValue: 789.4578429348551},
			Metadata: map[string]interface{}{common.SOURCE: testFilePath},
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
