package _team_member

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPutData struct {
	ID     string
	MEMBER string
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

	coll := client.Database("hvData").Collection("team")
	teamId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": teamId}
	var result bson.M
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	newMemberList := map[string]string{}
	for i, v := range result {
		if i == "member" {
			values := reflect.ValueOf(v)
			for k, v := range values.Interface().(primitive.M) {
				if data.MEMBER != v {
					newMemberList[k] = v.(string)
				}
			}
		}
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter = bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"member": newMemberList}}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}
