package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_component "fiber/Tools/Router/component"
	_component_like "fiber/Tools/Router/component/like"
	_project "fiber/Tools/Router/project"
	_style "fiber/Tools/Router/style"
	_team "fiber/Tools/Router/team"
	_team_member "fiber/Tools/Router/team/member"
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

	go app.Get("/component", _component.Get)
	go app.Post("/component", _component.Post)
	go app.Put("/component", _component.Put)
	go app.Delete("/component", _component.Delete)
	go app.Put("/component/like", _component_like.Put)

	go app.Post("/project", _project.Post)
	go app.Put("/project", _project.Put)
	go app.Delete("/project", _project.Delete)
	go app.Get("/project/:params", _project.Get)

	go app.Get("/team/:params", _team.Get)
	go app.Post("/team", _team.Post)
	go app.Delete("/team", _team.Delete)
	go app.Put("/team", _team.Put)
	go app.Put("/team/member", _team_member.Put)

	go app.Get("/user/:params", _user.Get)
	go app.Post("/user", _user.Post)

	go app.Get("/style/:params", _style.Get)

	go app.Get("/testapi", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(_component.Test())
	})

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
