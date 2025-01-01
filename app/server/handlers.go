package server

import (
		"github.com/gofiber/fiber/v2"
)

func HandleHello(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello!"})
}