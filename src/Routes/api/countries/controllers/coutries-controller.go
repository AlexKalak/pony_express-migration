package controllers

import (
	"net/http"

	"github.com/alexkalak/pony_express/src/Routes/api/countries/services"
	"github.com/gofiber/fiber/v2"
)

var CountriesService = services.New()

func CountriesController(router fiber.Router) {
	router.Get("/", getCountries)
}

func getCountries(c *fiber.Ctx) error {
	countries, err := CountriesService.GetAllCountries(c)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"ok":        true,
		"countries": countries,
	})
}
