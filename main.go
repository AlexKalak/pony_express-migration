package main

import (
	"fmt"
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
	fmt.Println("port: ", port)
	if port == "" {
		port = "80"
	}
	app.Route("/api", apiRouter.ApiRouter)

	app.Listen(":" + port)
	fmt.Println(app)
	fmt.Println("end")
}
