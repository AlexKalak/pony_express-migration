package services

import (
	"fmt"
	"math"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/helpers/price_helper"
	"github.com/alexkalak/pony_express/src/helpers/weights_helper"
	"github.com/alexkalak/pony_express/src/models"
)

func GetPrice(weight float64, regionID int, packageTypeID int, packageType string, senderCityID int) (int, error) {
	allWeightsOfPackage, err := GetAllWeightsOfPackage(regionID, packageTypeID)
	if err != nil {
		return 0, err
	}

	maxWeight := GetMaxWeightFromArray(allWeightsOfPackage)
	RoundUpWeight(&weight, allWeightsOfPackage, maxWeight)

	if weight < maxWeight {
		weightFromDB, err := weights_helper.GetWeightFromDB(weight)
		if err != nil {
			return 0, err
		}

		priceFromDB, err := price_helper.GetPriceFromDB(regionID, packageTypeID, weightFromDB.ID, senderCityID)
		if err != nil {
			return 0, err
		}
		return priceFromDB.Price, nil
	}

	p_maxWeightFromDB, err := weights_helper.GetWeightFromDB(maxWeight)
	if err != nil {
		return 0, err
	}

	maxPriceFromDB, err := price_helper.GetPriceFromDB(regionID, packageTypeID, p_maxWeightFromDB.ID, senderCityID)
	if err != nil {
		return 0, nil
	}

	overPrice, hasOverPrice := GetOverPrice(regionID, packageTypeID, senderCityID)
	if !hasOverPrice {
		return GetPriceUsingMaxWeightAndReminder(weight, maxWeight, maxPriceFromDB, senderCityID)
	} else {
		return GetPriceUsingOverPrice(weight, maxWeight, overPrice, maxPriceFromDB)
	}

}

func GetPriceUsingMaxWeightAndReminder(weight float64, maxWeight float64, maxPriceFromDB *models.Price, senderCityID int) (int, error) {
	numOfMaxWeights := int(weight / maxWeight)
	reminder := weight - float64(numOfMaxWeights)*maxWeight

	fmt.Println("num of max weights: ", numOfMaxWeights)
	fmt.Println("reminder: ", reminder)

	price := numOfMaxWeights * maxPriceFromDB.Price

	err := addReminderToPrice(&price, reminder, maxPriceFromDB.RegionID, maxPriceFromDB.PackageTypeID, senderCityID)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func GetPriceUsingOverPrice(weight float64, maxWeight float64, overPrice *models.PriceOverMaxWeight, maxPriceFromDB *models.Price) (int, error) {
	numOfMaxWeights := int(weight / maxWeight)
	reminder := weight - float64(numOfMaxWeights)*maxWeight
	numOfOverPrices := int(math.Round(reminder / overPrice.Weight.Weight))

	fmt.Println("num of max weights: ", numOfMaxWeights)
	fmt.Println("num of over prices: ", numOfOverPrices)

	//Finding price for max weight in db

	price := numOfMaxWeights * maxPriceFromDB.Price
	price += overPrice.Price * numOfOverPrices

	return price, nil
}

func GetOverPrice(regionID int, packageTypeID int, senderCityID int) (*models.PriceOverMaxWeight, bool) {
	database := db.GetDB()

	var overPrice models.PriceOverMaxWeight
	res := database.Preload("Weight").Find(&overPrice, "region_id = ? AND package_type_id = ? AND sender_city_id = ?", regionID, packageTypeID, senderCityID)
	if res.RowsAffected < 1 {
		return nil, false
	}

	if overPrice.ID == 0 {
		return nil, false
	}

	return &overPrice, true
}

func addReminderToPrice(price *int, reminder float64, regionID int, packageTypeID int, senderCityID int) error {
	if reminder > 0 {
		p_reminderFromDB, err := weights_helper.GetWeightFromDB(reminder)
		if err != nil {
			return err
		}

		//Finding price for reminder weight in db
		reminderPriceFromDB, err := price_helper.GetPriceFromDB(regionID, packageTypeID, p_reminderFromDB.ID, senderCityID)
		if err != nil {
			return err
		}

		*price += reminderPriceFromDB.Price
	}

	return nil
}
