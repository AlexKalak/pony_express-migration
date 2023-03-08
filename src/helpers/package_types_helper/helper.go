package package_types_helper

import (
	"errors"

	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
)

func GetPackageTypeFromDB(name string) (*models.PackageType, error) {
	database := db.GetDB()

	var packageTypeFromDB models.PackageType
	res := database.First(&packageTypeFromDB, "name = ?", name)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected < 1 {
		return nil, errors.New("package type not found")
	}

	return &packageTypeFromDB, nil
}
