package _style_P

import (
	"context"
	"fiber/Tools/mongodb"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Get = func(c *fiber.Ctx) error {
	data, err := getData(mongodb.GetMongoClient())
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

func getData(client *mongo.Client) (map[string][]primitive.M, error) {
	coll := client.Database("hvData").Collection("user")
	var result bson.M
	filter := bson.D{{Key: "name", Value: "최우승"}}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title")
		return map[string][]primitive.M{"data": {}}, nil
	}
	dataArray := []primitive.M{result}
	return map[string][]primitive.M{"data": dataArray}, err
}
