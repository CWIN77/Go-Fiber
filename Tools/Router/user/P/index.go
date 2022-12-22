package router

import (
	"github.com/gofiber/fiber/v2"
)

func Get() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString(c.Params("params"))
	}
}
