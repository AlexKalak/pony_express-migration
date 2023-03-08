package region_helper

import (
	"fmt"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetRegionByID(id int) (*models.Region, error) {
	database := db.GetDB()
	var region models.Region
	res := database.First(&region, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected < 1 {
		return nil, fmt.Errorf("err region not found")
	}

	return &region, nil
}
