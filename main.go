package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getData() []byte {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("simpleMongo").Collection("Post")
	title := "Join the MongoDB Community"
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "title", Value: title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return []byte("NoData")
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	return jsonData
}

type Request struct {
	Id        string `json:"_id"`
	AuthorId  string `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updatedAt"`
}

func main() {
	app := fiber.New()
	envLoad()
	middleware(app)

	fmt.Println(os.Getenv("TEST"))

	app.Get("/", func(c *fiber.Ctx) error {
		data := Request{}
		json.Unmarshal(getData(), &data)
		return c.Status(200).JSON(data)
	})

	app.Listen(":3000")
}

func envLoad() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func middleware(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
}
