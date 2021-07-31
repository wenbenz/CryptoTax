package nicehash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wenbenz/go-nicehash/client"
)

func TestNicehashReader(t *testing.T) {
	nh, _ := client.NewClientReadFrom("../data/nicehash.credentials")
	reports, _ := nh.GetReportsList()
	report, _ := nh.GetReport(reports[0].Id)
	reader := NewCsvStreamReader(report)
	read, _ := reader.Next()
	assert.Empty(t, read)
}
