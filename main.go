package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	_component "fiber/Tools/Router/component"
	_component_like "fiber/Tools/Router/component/like"
	_project "fiber/Tools/Router/project"
	_project_component "fiber/Tools/Router/project/component"
	_style "fiber/Tools/Router/style"
	_team "fiber/Tools/Router/team"
	_team_member "fiber/Tools/Router/team/member"
	_test "fiber/Tools/Router/test"
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
	fmt.Println(primitive.NewObjectID().Hex())

	app.Get("/component", _component.Get)
	app.Post("/component", _component.Post)
	app.Put("/component", _component.Put)
	app.Delete("/component", _component.Delete)
	app.Put("/component/like", _component_like.Put)

	app.Post("/project", _project.Post)
	app.Put("/project", _project.Put)
	app.Delete("/project", _project.Delete)
	app.Get("/project/:params", _project.Get)
	app.Put("/project/component", _project_component.Put)
	app.Delete("/project/component", _project_component.Delete)

	app.Get("/team/:params", _team.Get)
	app.Post("/team", _team.Post)
	app.Delete("/team", _team.Delete)
	app.Put("/team", _team.Put)
	app.Put("/team/member", _team_member.Put)

	app.Get("/user/:params", _user.Get)
	app.Post("/user", _user.Post)

	app.Get("/style/:params", _style.Get)

	if os.Getenv("ENV_MODE") == "development" {
		app.Get("/test", _test.Test)
	}

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
