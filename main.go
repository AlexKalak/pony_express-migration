package main

import (
	"os"
	"time"

	apiRouter "github.com/alexkalak/pony_express/src/Routes/api"
	currencyhelper "github.com/alexkalak/pony_express/src/currencyHelper"
	"github.com/alexkalak/pony_express/src/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Init()
	go currencyhelper.StartGettingCurrencies(time.Hour * 24)

	app := fiber.New()

	app.Use(cors.New())
	app.Static("/static", "dist/static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "9999"
	}
	app.Route("/api", apiRouter.ApiRouter)

	app.Listen("0.0.0.0:" + port)
}
