package _component

import (
	"context"
	"fiber/Tools/mongodb"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ComponentData struct {
	NAME  string
	HTML  string
	STYLE string
	MAKER string
}

var Get = func(c *fiber.Ctx) error {
	search := strings.ToLower(c.Query("search"))
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	skip, _ := strconv.ParseInt(c.Query("skip"), 10, 32)
	data, err := getData(search, limit, skip)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := ComponentData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	values := reflect.ValueOf(p)
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			return c.Status(400).JSON("Please send all user data.")
		}
	}
	result, err := postData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

func getData(search string, limit int64, skip int64) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	filter := bson.M{"search": bson.M{"$regex": search}}
	opts := options.Find().SetSort(bson.M{"like": -1}).SetLimit(limit).SetSkip(skip)

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

func postData(componentData ComponentData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	insertData := bson.M{
		"name":      componentData.NAME,
		"keyword":   strings.ToLower(componentData.NAME),
		"html":      componentData.HTML,
		"style":     componentData.STYLE,
		"maker":     componentData.MAKER,
		"like":      0,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}
