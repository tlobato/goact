package main

import (
		"log"

		"github.com/tlobato/goact/server"

		"github.com/joho/godotenv"
		"github.com/gofiber/fiber/v2"
		"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
		godotenv.Load(".env")

		app := fiber.New()

		//handle frontend routes
		app.Static("/", "./client/dist")

		api := app.Group("/api", cors.New(cors.Config{
			AllowOrigins: "http://localhost:3000", //change this later
			AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		}))
		api.Get("/hello", server.HandleHello)

		log.Fatal(app.Listen(":3000"))
}
