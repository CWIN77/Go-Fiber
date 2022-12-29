package _component

import (
	"context"
	"fiber/Tools/mongodb"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type ComponentData struct {
// 	name      string `json:"name" xml:"name" form:"name"`
// 	html      string `json:"html" xml:"HTML" form:"HTML"`
// 	style     string `json:"STYLE" xml:"STYLE" form:"STYLE"`
// 	maker     string `json:"email" xml:"email" form:"email"`
// 	like      string `json:"email" xml:"email" form:"email"`
// 	createdAt string `json:"email" xml:"email" form:"email"`
// 	updatedAt string `json:"email" xml:"email" form:"email"`
// }

// "name": ,
// "search":,
// "html": ,
// "style": ,
// "maker": ,
// "like": ,
// "createdAt": ,
// "updatedAt":,

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

// var Post = func(c *fiber.Ctx) error {

// }

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

// func postData()  {
// 	client := mongodb.GetMongoClient()
// 	coll := client.Database("hvData").Collection("user")
// 	var findData bson.M
// 	insertData := bson.M{
// 		"name": ,
// 		"search":,
// 		"html": ,
// 		"style": ,
// 		"maker": ,
// 		"like": ,
// 		"createdAt": ,
// 		"updatedAt":,
// 	}
// 	if err := coll.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&findData); err == mongo.ErrNoDocuments {
// 		result, err := coll.InsertOne(context.TODO(), insertData)
// 		return result, err
// 	} else if err != nil {
// 		return nil, err
// 	}
// 	for key, value := range insertData {
// 		if findData[key] != value {
// 			filter := bson.D{{Key: "_id", Value: userId}}
// 			update := bson.M{"$set": insertData}
// 			updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
// 			return updateResult, err
// 		}
// 	}
// 	return "Data already exist and fresh state", nil
// }
