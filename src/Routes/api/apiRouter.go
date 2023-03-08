package api

import (
	calculatorController "github.com/alexkalak/pony_express/src/Routes/api/calculator/controllers"
	coutriesController "github.com/alexkalak/pony_express/src/Routes/api/countries/controllers"
	"github.com/gofiber/fiber/v2"
)

func ApiRouter(router fiber.Router) {
	router.Route("/countries", coutriesController.CountriesController)
	router.Route("/calculator", calculatorController.CalculatorController)
}
