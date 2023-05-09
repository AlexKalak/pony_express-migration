package migration

import (
	"github.com/alexkalak/migration/src/db"
	"github.com/alexkalak/migration/src/models"
)

func MigrateDeliveryTypes() {
	database := db.GetDB()
	database.Create(&models.DeliveryType{
		Name: "documents",
	})
	database.Create(&models.DeliveryType{
		Name: "standart",
	})
	database.Create(&models.DeliveryType{
		Name: "B2B",
	})
}
