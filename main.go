package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_index "fiber/Tools/Router"
	__component_P "fiber/Tools/Router/_api/_component/P"
	__project_P "fiber/Tools/Router/_api/_project/P"
	__team_P "fiber/Tools/Router/_api/_team/P"
	__user_P "fiber/Tools/Router/_api/_user/P"
	middleware "fiber/Tools/middleware"
	"fiber/Tools/mongodb"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	middleware.Setting(app)

	mongoClient := mongodb.Connect(os.Getenv("MONGODB_URI"))

	app.Get("/", _index.Get(mongoClient))

	app.Get("/api/component/:makerId", __component_P.Get())
	app.Get("/api/project/:ownerId", __project_P.Get())
	app.Get("/api/team/:memberId", __team_P.Get())
	app.Get("/api/user/:userId", __user_P.Get())

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
