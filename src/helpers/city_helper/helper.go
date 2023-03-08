package city_helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetCityByName(name string) (*models.City, error) {
	database := db.GetDB()
	var city models.City
	res := database.Model(&models.City{}).Preload("District.Area").Where("name = ?", name).First(&city)
	if res.Error != nil {
		return nil, res.Error
	}

	return &city, nil
}

func GetCityByCityNameCountryDistrictAndArea(name string, countryName, districtName string, areaName string) (*models.City, error) {
	database := db.GetDB()

	conditions := make([]string, 0, 4)
	conditions = append(conditions, fmt.Sprintf("name = '%s'", name))

	areaFromDB, err := GetAreaByName(areaName)
	if err != nil {
		return nil, err
	}

	if areaFromDB != nil {
		conditions = append(conditions, "area_id = "+strconv.Itoa(areaFromDB.ID))
	} else {
		conditions = append(conditions, "area_id is NULL")
	}

	districtFromDB, err := GetDistrictByNameAndArea(districtName, areaName)
	if err != nil {
		return nil, err
	}

	if districtFromDB != nil {
		conditions = append(conditions, "district_id = "+strconv.Itoa(districtFromDB.ID))
	} else {
		conditions = append(conditions, "district_id is NULL")
	}

	countryFromDB, err := GetCountryByName(countryName)
	if err != nil {
		return nil, err
	}

	if countryFromDB != nil {
		conditions = append(conditions, "country_id = "+strconv.Itoa(countryFromDB.ID))
	} else {
		conditions = append(conditions, "country_id is NULL")
	}

	var city models.City

	queryConditionsStr := strings.Join(conditions[:], " AND ")

	res := database.
		Model(&models.City{}).
		Preload("District").
		Preload("Area").
		Where(queryConditionsStr).
		First(&city)

	if res.Error != nil {
		return nil, res.Error
	}
	// str, _ := json.MarshalIndent(city, "", "\t")
	// fmt.Println(string(str))

	return &city, nil
}
func GetCityByCityNameCountryAndArea(name string, countryName string, areaName string) (*models.City, error) {
	database := db.GetDB()

	countryFromDB, err := GetCountryByName(countryName)
	if err != nil {
		return nil, err
	}
	var country_id *int
	if countryFromDB != nil {
		country_id = &countryFromDB.ID
	}

	areaFromDB, err := GetAreaByName(areaName)
	if err != nil {
		return nil, err
	}
	var areaID *int
	if areaFromDB != nil {
		areaID = &areaFromDB.ID
	}

	var city models.City
	res := database.
		Model(&models.City{}).
		Preload("District.Area").
		Where("name = ? AND county_id and area_id = ?", name, country_id, areaID).
		First(&city)

	if res.Error != nil {
		return nil, res.Error
	}

	return &city, nil
}

func GetSenderCityByName(name string) (*models.SenderCity, error) {
	database := db.GetDB()
	var senderCity models.SenderCity
	res := database.First(&senderCity, "name = ?", name)
	if res.Error != nil {
		return nil, res.Error
	}

	return &senderCity, nil
}

func GetAreaByName(name string) (*models.Area, error) {
	database := db.GetDB()

	if name == "" {
		return nil, nil
	}

	var area models.Area
	res := database.First(&area, "name = ?", name)
	if res.Error != nil {
		return nil, res.Error
	}

	return &area, nil
}

func GetDistrictByNameAndArea(name string, areaName string) (*models.District, error) {
	database := db.GetDB()
	if name == "" {
		return nil, nil
	}

	var condition string
	areaFromDB, err := GetAreaByName(areaName)
	if err != nil {
		return nil, err
	}
	if areaFromDB == nil {
		condition = "area_id IS NULL"
	} else {
		condition = fmt.Sprintf("area_id = %d", areaFromDB.ID)
	}

	fmt.Println("condition: ", condition)

	var district models.District
	res := database.
		Where(condition).
		First(&district, "name = ?", name)

	if res.Error != nil {
		return nil, res.Error
	}

	return &district, nil
}

func GetCountryByName(name string) (*models.Country, error) {
	database := db.GetDB()
	if name == "" {
		return nil, nil
	}

	var country models.Country
	res := database.First(&country, "name = ?", name)
	if res.Error != nil {
		return nil, res.Error
	}

	return &country, nil
}
