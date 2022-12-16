package _index

import (
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(client *mongo.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data, err := mongodb.GetData(client)
		if err != nil {
			return c.Status(400).JSON(err.Error())
		}
		return c.Status(200).JSON(data)
	}
}
