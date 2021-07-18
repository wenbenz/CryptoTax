package cryptotax

import (
	"encoding/csv"
	"os"
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
