package main

import (
	"os"

	apiRouter "github.com/alexkalak/pony_express-calculator/src/Routes/api"
	"github.com/alexkalak/pony_express-calculator/src/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Init()

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
