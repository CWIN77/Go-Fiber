package _component_like

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPutData struct {
	COMP_ID string
	USER_ID string
}

var Put = func(c *fiber.Ctx) error {
	p := TPutData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	result, err := putData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

func putData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()

	coll := client.Database("hvData").Collection("component")
	compId, err := primitive.ObjectIDFromHex(data.COMP_ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": compId}
	update := bson.M{"$pull": bson.M{"like": bson.M{"$eq": data.USER_ID}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if result.ModifiedCount == 0 {
		update := bson.M{"$push": bson.M{"like": data.USER_ID}}
		result, err := coll.UpdateOne(context.TODO(), filter, update)
		return result, err
	}
	return nil, err
}
