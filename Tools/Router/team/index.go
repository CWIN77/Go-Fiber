package _team

import (
	"context"
	"fiber/Tools/mongodb"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

func getData(memberId string) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	filter := bson.M{
		"$or": [4]interface{}{
			bson.M{"member.master": memberId},
			bson.M{"member.manager": memberId},
			bson.M{"member.maker": memberId},
			bson.M{"member.reader": memberId},
		},
	}
	opts := options.Find()
	var err error = nil
	var cursor *mongo.Cursor
	cursor, err = coll.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	dataArray := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		dataMap := map[string]interface{}{}
		for _, k := range result {
			dataMap[k.Key] = k.Value
		}
		dataArray = append(dataArray, dataMap)
	}
	return dataArray, err
}
