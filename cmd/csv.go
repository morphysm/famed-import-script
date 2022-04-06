package cmd

import (
	"encoding/csv"
	"os"
)

// readCSV return data read from a csv file.
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
