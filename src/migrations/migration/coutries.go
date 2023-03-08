package migration

import (
	"strconv"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func MigrateCountries() {
	// database := db.GetDB()
	array := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/coutries.csv")

	for _, entity := range array {
		SaveIfNotExistsCountryCode(entity[1])
		SaveIfNotExistsCountry(entity[0], entity[1], entity[2])
	}
}

func SaveIfNotExistsCountryCode(code string) {
	database := db.GetDB()

	countryCode := models.CountryCode{}

	database.Model(&models.CountryCode{}).Where("code = ?", code).Find(&countryCode)
	if countryCode.ID == 0 {
		database.Create(&models.CountryCode{Code: code})
	}
}

func SaveIfNotExistsCountry(countryName string, code string, regionID string) {
	database := db.GetDB()

	countryCode := models.CountryCode{}
	database.Model(&models.CountryCode{}).Where("code = ?", code).Find(&countryCode)

	country := models.Country{}
	database.Model(&models.Country{}).Where("name = ?", countryName).Find(&country)
	if country.ID != 0 {
		return
	}

	country.CountryCode = countryCode
	country.Name = countryName
	if regionID != "" {
		intRegID, err := strconv.Atoi(regionID)
		if err != nil {
			panic(err)
		}
		country.RegionID = intRegID
	}
	database.Create(&country)

}
