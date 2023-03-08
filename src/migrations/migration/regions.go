package migration

import (
	"fmt"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func MigrateRegions() {
	array := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/regions.csv")

	for _, entity := range array {
		fmt.Println(entity)
		SaveIfNotExists(entity[0])
	}
}

func SaveIfNotExists(regionName string) {
	database := db.GetDB()
	var region models.Region
	fmt.Println(regionName)
	database.Find(&region, "name = ?", regionName)
	if region.ID != 0 {
		return
	}

	region.Name = regionName
	database.Create(&region)
}
