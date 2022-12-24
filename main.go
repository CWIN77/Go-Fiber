package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	_index "fiber/Tools/Router"
	_component "fiber/Tools/Router/component"
	_project_P "fiber/Tools/Router/project/P"
	_team_P "fiber/Tools/Router/team/P"
	_user_P "fiber/Tools/Router/user/P"
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

	app.Get("/", _index.Get)

	app.Get("/component", _component.Get)
	app.Get("/project/:params", _project_P.Get)
	app.Get("/team/:params", _team_P.Get)
	app.Get("/user/:params", _user_P.Get)

	app.Get("/test", _project_P.Get)

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
