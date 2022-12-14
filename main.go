package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/joho/godotenv"

	mongodb "src/fiber/Tools"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: "secret-thirty-2-character-string",
	}))

	// env 연결
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	client := mongodb.Connect(os.Getenv("MONGODB_URI"))

	app.Get("/", func(c *fiber.Ctx) error {
		data := mongodb.GetData(client)
		return c.Status(200).JSON(data)
	})

	app.Listen(":5000") // Elastic Beanstalk에 배포시 5000포트이고 파일이름이 application.go여야한다
}
