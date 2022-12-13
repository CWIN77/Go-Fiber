package main

import (
	"log"

	"github.com/gofiber/fiber"
	// "github.com/gofiber/fiber/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Send("11")
	})

	log.Fatal(app.Listen(":3000"))
}
