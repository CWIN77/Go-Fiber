package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	mongodb "src/fiber/Tools"
)

func main() {
	app := fiber.New()

	// env 연결
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	client := mongodb.Connect(os.Getenv("MONGODB_URI"))

	app.Use(cors.New())
	app.Use(compress.New())

	app.Get("/", func(c *fiber.Ctx) error {
		data := mongodb.GetData(client)
		return c.Status(200).JSON(data)
	})

	app.Listen(":5000") // Elastic Beanstalk에 배포시 5000포트이고 파일이름이 application.go여야한다
}
