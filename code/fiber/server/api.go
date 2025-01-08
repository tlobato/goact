package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApiRoutes(app *fiber.App) {
	api := app.Group("/api", cors.New(cors.Config{
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
	}))

	api.Get("/hello", helloHandler)
}

func helloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World!"})
}
