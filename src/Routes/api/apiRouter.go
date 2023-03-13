package api

import (
	calculatorController "github.com/alexkalak/pony_express-calculator/src/Routes/api/calculator/controllers"
	"github.com/gofiber/fiber/v2"
)

func ApiRouter(router fiber.Router) {
	router.Route("/calculator", calculatorController.CalculatorController)
}
