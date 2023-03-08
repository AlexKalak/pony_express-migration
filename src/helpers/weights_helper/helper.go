package weights_helper

import (
	"errors"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetWeightFromDB(weight float64) (*models.Weight, error) {
	database := db.GetDB()

	var weightFromDB models.Weight
	res := database.First(&weightFromDB, "weight = ?", weight)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected < 1 {
		return nil, errors.New("weight not found")
	}

	return &weightFromDB, nil
}
