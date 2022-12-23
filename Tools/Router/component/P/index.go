package _component_P

import (
	"github.com/gofiber/fiber/v2"
)

var Get = func(c *fiber.Ctx) error {
	return c.SendString(c.Params("params"))
}
