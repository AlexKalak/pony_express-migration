package controllers

import (
	"net/http"

	"github.com/alexkalak/pony_express/src/Routes/api/calculator/services"
	"github.com/gofiber/fiber/v2"
)

var CalculatorService = services.New()

func CalculatorController(router fiber.Router) {
	router.Post("/", calculate)
}

func calculate(c *fiber.Ctx) error {
	prices, usedVolumeWeights, validationErrors, err := CalculatorService.Calculate(c)
	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"ok":  false,
			"err": err.Error(),
		})
	}

	if len(validationErrors) > 0 {
		c.SendStatus(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"ok":               false,
			"validationErrors": validationErrors,
		})
	}

	return c.JSON(fiber.Map{
		"ok":                 true,
		"used-volume-weight": usedVolumeWeights,
		"prices":             prices,
	})
}
