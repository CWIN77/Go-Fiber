package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	middleware "fiber/Tools"
	_index "fiber/Tools/Router"
	_teamId_projectId "fiber/Tools/Router/pTeamId/pProjectId"
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

	app.Get("/:teamId/:projectId", _teamId_projectId.Get())
	// app.Get("/api/:id", _api_id.Get())

	// Elastic Beanstalk Deploy Port : 5000
	// Elastic Beanstalk Main Name : application
	app.Listen(":5000")
}
