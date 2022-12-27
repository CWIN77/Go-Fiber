package _user

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error  {
  return c.Send(c.Body())
}

func getData(id string) (primitive.M, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	var result bson.M
	userId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userId}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

// func createData() {
  
// }