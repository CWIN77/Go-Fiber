package _style_P

import (
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
)

var Get = func(c *fiber.Ctx) error {
	data, err := mongodb.GetData(mongodb.GetMongoClient())
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}
