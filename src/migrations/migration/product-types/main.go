package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/alexkalak/migration/src/db"
	"github.com/alexkalak/migration/src/models"
)

func main() {
	database := db.GetDB()

	database.Migrator().DropTable("product_types")
	database.AutoMigrate(&models.ProductType{})

	records := ReadCSV("/home/alexkalak/Desktop/pony_express-migration/csvtables/product-types.csv")
	for _, record := range records {
		gtipCode, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}

		ruName := record[1]

		trName := record[2]

		enName := record[3]

		itemCode, err := strconv.Atoi(record[4])
		if err != nil {
			panic(err)
		}

		weight, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			panic(err)
		}
		weight = roundFloat(weight, 2)

		warning, err := strconv.ParseBool(record[6])
		if err != nil {
			panic(err)
		}

		productType := models.ProductType{
			GtipCode: gtipCode,
			RuName:   ruName,
			TrName:   trName,
			EnName:   enName,
			ItemCode: itemCode,
			Weight:   weight,
			Warning:  warning,
		}

		database.Create(&productType)
	}

}

func ReadCSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	// length += len(records)*len(records[0]) - len(records)
	fmt.Println(len(records))
	return records
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
