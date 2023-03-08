package price_helper

import (
	"errors"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetPriceFromDB(regionID int, packageTypeID int, weightID int, senderCityID int) (*models.Price, error) {
	database := db.GetDB()

	var priceFromDB models.Price
	res := database.First(&priceFromDB, "region_id = ? AND package_type_id = ? AND weight_id = ? AND sender_city_id = ?", regionID, packageTypeID, weightID, senderCityID)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected < 1 {
		return nil, errors.New("price not found")
	}

	return &priceFromDB, nil
}
