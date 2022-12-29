package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_component "fiber/Tools/Router/component"
	_project "fiber/Tools/Router/project"
	_style "fiber/Tools/Router/style"
	_team "fiber/Tools/Router/team"
	_user "fiber/Tools/Router/user"
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

	if err := mongodb.ConnectDB(os.Getenv("MONGODB_URI")); err != nil {
		log.Fatal("Error mongoDB connect")
	}

	app.Get("/component", _component.Get)
	app.Get("/project/:params", _project.Get)
	app.Get("/team/:params", _team.Get)

	app.Get("/user/:params", _user.Get)
	app.Post("/user", _user.Post)

	app.Get("/style/:params", _style.Get)

	// app.Get("/test", func(c *fiber.Ctx) error {
	// 	return c.Status(200).JSON("test")
	// })

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
