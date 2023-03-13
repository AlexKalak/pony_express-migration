package migration

import (
	"fmt"
	"strconv"

	"github.com/alexkalak/pony_express-calculator/src/db"
	"github.com/alexkalak/pony_express-calculator/src/models"
)

func MigrateCities() {
	database := db.GetDB()

	var RussiaFromDB models.Country
	var MoldovaFromDB models.Country
	var UkraineFromDB models.Country
	database.Find(&RussiaFromDB, "name = ?", "Россия")
	database.Find(&MoldovaFromDB, "name = ?", "Молдова")
	database.Find(&UkraineFromDB, "name = ?", "Украина")

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

	res := database.Find(&Istanbul, "tr_name = ?", "Istanbul")
	if res.Error != nil {
		panic(res.Error)
	}
	res = database.Find(&Antalya, "tr_name = ?", "Antalya")
	if res.Error != nil {
		panic(res.Error)
	}

	AsianIstanbulReg := models.SenderRegion{
		Name:         "Азиатская часть Стамбула",
		TrName:       "Istanbul'un Asya yakası",
		PriceForDoor: 2700,
		SenderCity:   Istanbul,
	}
	EuropeIstanbulReg := models.SenderRegion{
		Name:         "Европейская часть Стамбула",
		TrName:       "İstanbul'un Avrupa yakası",
		PriceForDoor: 1600,
		SenderCity:   Istanbul,
	}
	AntalyaReg := models.SenderRegion{
		Name:         "Анталия",
		TrName:       "Antalya",
		PriceForDoor: 1100,
		SenderCity:   Antalya,
	}

	database.Create(&AsianIstanbulReg)
	database.Create(&EuropeIstanbulReg)
	database.Create(&AntalyaReg)
}

func MigrateSenderCities() {
	database := db.GetDB()
	Istanbul := models.SenderCity{
		Name:   "Стамбул",
		TrName: "Istanbul",
	}
	Antalya := models.SenderCity{
		Name:   "Анталия",
		TrName: "Antalya",
	}

	database.Create(&Istanbul)
	database.Create(&Antalya)
}

// ///////////////////////////// Areas ///////////////////////////////////////////////////
func MigrateRussiaAreas(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/country_cities.csv")

	for _, entity := range arr {
		id, err := strconv.Atoi(entity[0])
		if err != nil {
			panic(err)
		}
		areaName := entity[2]
		areaTrName := entity[3]

		SaveArea(id, areaName, areaTrName, countryFromDB)
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
		areaTrName := entity[3]

		SaveArea(id, areaName, areaTrName, countryFromDB)
	}
}

func SaveArea(id int, areaName string, trName string, countryFromDB *models.Country) {
	database := db.GetDB()
	area := models.Area{
		ID:     id,
		Name:   areaName,
		TrName: trName,

		Country: *countryFromDB,
	}

	database.Save(&area)
}

// //////////////////////////// Districts ///////////////////////////////////////////
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
		districtTrName := entity[3]

		SaveDistrict(id, districtName, districtTrName, areaID)
		// SaveDistrict(districtName, areaID)
	}
}

func SaveDistrict(id int, districtName string, districtTrName string, areaID int) {
	// func SaveDistrict(districtName string, areaID int) {
	database := db.GetDB()
	area := models.District{
		ID:     id,
		Name:   districtName,
		TrName: districtTrName,
		AreaID: areaID,
	}

	database.Save(&area)
}

// //////////////////////////// Cities and big and small ///////////////////////////////////////////
func MigrateRussianCities(countryFromDB *models.Country) {
	arrAllCities := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/city_places.csv")
	arrBigCities := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/russia/cities/russia-big-cities.csv")

	SaveIfNotExistCity("Москва", "Moskva", 15, countryFromDB, nil, nil)
	SaveIfNotExistCity("Санкт-Петербург", "Sankt-Peterburg", 15, countryFromDB, nil, nil)

	for _, entity := range arrAllCities {
		if entity[0] != "180" {
			continue
		}

		cityName := entity[4]
		cityTrName := entity[5]

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
			SaveIfNotExistCity(cityName, cityTrName, regionID, countryFromDB, nil, &areaID)
			continue
		}

		districtID, err := strconv.Atoi(entity[3])
		if err != nil {
			panic(err)
		}

		SaveIfNotExistCity(cityName, cityTrName, regionID, countryFromDB, &districtID, &areaID)
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

		cityName := entity[0]
		cityTrName := entity[3]

		if area_id != 0 {
			SaveIfNotExistCity(cityName, cityTrName, regionID, countryFromDB, nil, &area_id)
		} else {
			SaveIfNotExistCity(cityName, cityTrName, regionID, countryFromDB, nil, nil)
		}
	}
}

func MigrateUkraineCities(countryFromDB *models.Country) {
	arr := ReadCSV("/home/alexkalak/Desktop/pony_express/csvtables/ukraine/cities.csv")

	for _, entity := range arr {
		l_reg_id, _ := strconv.Atoi(entity[1])
		regionID := l_reg_id + 21

		cityName := entity[0]
		cityTrName := entity[2]

		SaveIfNotExistCity(cityName, cityTrName, regionID, countryFromDB, nil, nil)
	}
}

func SaveIfNotExistCity(cityName string, cityTrName string, regionID int, countryFromDB *models.Country, DistrictID *int, areaID *int) {
	fmt.Println(cityName)
	database := db.GetDB()
	var city models.City
	database.Find(&city, "name = ? and district_id = ?", cityName, DistrictID)
	if city.ID != 0 {
		return
	}

	city.Name = cityName
	city.TrName = cityTrName
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
