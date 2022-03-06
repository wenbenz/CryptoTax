package tools

import (
	"encoding/csv"
	"os"
	"strconv"
)

func ReadCsv(path string) ([][]string, error) {
	// read CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func ParseFloat(s string) (float64, error) {
	if len(s) == 0 {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}