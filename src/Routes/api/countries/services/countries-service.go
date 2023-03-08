package services

import (
	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/models"
	"github.com/gofiber/fiber/v2"
)

type CountriesService interface {
	GetAllCountries(c *fiber.Ctx) (*[]models.Country, error)
}

type countriesService struct {
}

func New() CountriesService {
	return &countriesService{}
}

func (cs *countriesService) GetAllCountries(c *fiber.Ctx) (*[]models.Country, error) {
	database := db.GetDB()

	p_countries := new([]models.Country)
	res := database.Preload("Cities.District").Preload("Cities.Area").Find(p_countries)
	if res.Error != nil {
		return nil, res.Error
	}

	return p_countries, nil
}
