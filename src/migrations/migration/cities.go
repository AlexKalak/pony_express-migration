package migration

import (
	"fmt"
	"strconv"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func MigrateCities() {
	database := db.GetDB()

	var RussiaFromDB models.Country
	var MoldovaFromDB models.Country
	var UkraineFromDB models.Country
	database.Find(&RussiaFromDB, "name = ?", "Rusya")
	database.Find(&MoldovaFromDB, "name = ?", "Moldova")
	database.Find(&UkraineFromDB, "name = ?", "Ukrayna")

	MigrateSenderCities()
	MigrateSenderRegions()

	MigrateRussiaAreas(&RussiaFromDB)
	MigrateMoldovaAreas(&MoldovaFromDB)
	MigrateDistricts(&RussiaFromDB)
	MigrateRussianCities(&RussiaFromDB)
	MigrateMoldovaCities(&MoldovaFromDB)
	MigrateUkraineCities(&UkraineFromDB)
}

func MigrateSenderRegions() {
	database := db.GetDB()
	var Istanbul models.SenderCity
	var Antalya models.SenderCity

	res := database.Find(&Istanbul, "name = ?", "Istanbul")
	if res.Error != nil {
		panic(res.Error)
	}
	res = database.Find(&Antalya, "name = ?", "Antalya")
	if res.Error != nil {
		panic(res.Error)
	}

	AsianIstanbulReg1 := models.SenderRegion{
		Name:         "Азиатская часть Стамбула",
		PriceForDoor: 2700,
		SenderCity:   Istanbul,
	}
	EuropeIstanbulReg1 := models.SenderRegion{
		Name:         "Европейская часть Стамбула",
		PriceForDoor: 1600,
		SenderCity:   Istanbul,
	}
	AsianIstanbulReg2 := models.SenderRegion{
		Name:         "Asian part of Istanbul",
		PriceForDoor: 2700,
		SenderCity:   Istanbul,
	}
	EuropeIstanbulReg2 := models.SenderRegion{
		Name:         "Europe part of Istanbul",
		PriceForDoor: 1600,
		SenderCity:   Istanbul,
	}
	AntalyaReg1 := models.SenderRegion{
		Name:         "Анталия",
		PriceForDoor: 1100,
		SenderCity:   Antalya,
	}
	AntalyaReg2 := models.SenderRegion{
		Name:         "Antalya",
		PriceForDoor: 1100,
		SenderCity:   Antalya,
	}

	database.Create(&AsianIstanbulReg1)
	database.Create(&EuropeIstanbulReg1)
	database.Create(&AntalyaReg1)
	database.Create(&AsianIstanbulReg2)
	database.Create(&EuropeIstanbulReg2)
	database.Create(&AntalyaReg2)
}

func MigrateSenderCities() {
	database := db.GetDB()
	Istanbul := models.SenderCity{
		Name: "Istanbul",
	}
	Antalya := models.SenderCity{
		Name: "Antalya",
	}

	database.Create(&Istanbul)
	database.Create(&Antalya)
}

// ///////////////////////////// CountryCities ///////////////////////////////////////////////////
func MigrateRussiaAreas(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/country_cities.csv")

	for _, entity := range arr {
		id, err := strconv.Atoi(entity[0])
		if err != nil {
			panic(err)
		}
		areaName := entity[2]

		SaveArea(id, areaName, countryFromDB)
	}
}

func MigrateMoldovaAreas(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/moldova/country_cities.csv")

	for _, entity := range arr {
		id, err := strconv.Atoi(entity[0])
		if err != nil {
			panic(err)
		}
		areaName := entity[2]

		SaveArea(id, areaName, countryFromDB)
	}
}

func SaveArea(id int, areaName string, countryFromDB *models.Country) {
	database := db.GetDB()
	area := models.Area{
		ID:      id,
		Name:    areaName,
		Country: *countryFromDB,
	}

	database.Save(&area)
}

// //////////////////////////// Cities and big and small ///////////////////////////////////////////
func MigrateDistricts(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/districts.csv")

	for _, entity := range arr {
		id, err := strconv.Atoi(entity[0])
		if err != nil {
			panic(err)
		}
		areaID, err := strconv.Atoi(entity[1])
		if err != nil {
			panic(err)
		}
		districtName := entity[2]

		SaveDistrict(id, districtName, areaID)
		// SaveDistrict(districtName, areaID)
	}
}

func SaveDistrict(id int, districtName string, areaID int) {
	// func SaveDistrict(districtName string, areaID int) {
	database := db.GetDB()
	area := models.District{
		ID:     id,
		Name:   districtName,
		AreaID: areaID,
	}

	database.Save(&area)
}

func MigrateRussianCities(countryFromDB *models.Country) {
	arrAllCities := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/city_places.csv")
	arrBigCities := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/russia-big-cities.csv")

	SaveIfNotExistCity("Москва", 15, countryFromDB, nil, nil)
	SaveIfNotExistCity("Санкт-Петербург", 15, countryFromDB, nil, nil)

	for _, entity := range arrAllCities {
		if entity[0] != "180" {
			continue
		}

		cityName := entity[4]

		regionID := 17
		if BigCityContains(arrBigCities, cityName) {
			regionID = 16
		}
		if cityName == "Москва" || cityName == "Санкт-Петербург" {
			regionID = 15
		}

		if entity[2] == "" {
			continue
		}

		areaID, err := strconv.Atoi(entity[2])
		if err != nil {
			panic(err)
		}

		if entity[3] == "" {
			SaveIfNotExistCity(cityName, regionID, countryFromDB, nil, &areaID)
			continue
		}

		districtID, err := strconv.Atoi(entity[3])
		if err != nil {
			panic(err)
		}

		SaveIfNotExistCity(cityName, regionID, countryFromDB, &districtID, &areaID)
	}
}
func MigrateMoldovaCities(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/moldova/cities.csv")

	for _, entity := range arr {
		l_reg_id, _ := strconv.Atoi(entity[1])
		regionID := l_reg_id + 17
		area_id, err := strconv.Atoi(entity[2])
		if err != nil {
			panic(err)
		}

		if area_id != 0 {
			SaveIfNotExistCity(entity[0], regionID, countryFromDB, nil, &area_id)
		} else {
			SaveIfNotExistCity(entity[0], regionID, countryFromDB, nil, nil)
		}
	}
}

func MigrateUkraineCities(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/ukraine/cities.csv")

	for _, entity := range arr {
		l_reg_id, _ := strconv.Atoi(entity[1])
		regionID := l_reg_id + 21
		SaveIfNotExistCity(entity[0], regionID, countryFromDB, nil, nil)
	}
}

func SaveIfNotExistCity(cityName string, regionID int, countryFromDB *models.Country, DistrictID *int, areaID *int) {
	fmt.Println(cityName)
	database := db.GetDB()
	var city models.City
	database.Find(&city, "name = ? and district_id = ?", cityName, DistrictID)
	if city.ID != 0 {
		return
	}

	city.Name = cityName
	city.Country = *countryFromDB
	city.Region.ID = regionID
	city.DistrictID = DistrictID
	city.AreaID = areaID
	database.Create(&city)
}

func BigCityContains(bigCityArr [][]string, cityName string) bool {
	for _, city := range bigCityArr {
		if city[0] == cityName {
			return true
		}
	}
	return false
}
