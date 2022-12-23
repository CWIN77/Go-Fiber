package _component_P

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
	data, err := getData(mongodb.GetMongoClient())
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

func getData(client *mongo.Client) ([]map[string]interface{}, error) {
	coll := client.Database("hvData").Collection("component")
	opts := options.Find().SetSort(bson.D{{Key: "lije", Value: -1}}).SetLimit(10).SetSkip(0)
	var err error = nil
	var cursor *mongo.Cursor
	cursor, err = coll.Find(context.TODO(), bson.D{}, opts)
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
