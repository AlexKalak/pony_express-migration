package migration

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("a")
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println("a")

	fmt.Println(records)
	return records
}
