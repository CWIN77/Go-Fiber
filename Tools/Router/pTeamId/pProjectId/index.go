package _teamId_projectId

import (
	"github.com/gofiber/fiber/v2"
)

func Get() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString(c.Params("teamId") + ", " + c.Params("projectId"))
	}
}
