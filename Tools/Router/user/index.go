package _user

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserData struct {
	UID   string
	NAME  string
	IMG   string
	EMAIL string
}

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := UserData{}
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

func getData(id string) (primitive.M, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	var result bson.M
	userId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userId}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func postData(userData UserData) (interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	userId, _ := primitive.ObjectIDFromHex(userData.UID)
	var findData bson.M
	insertData := bson.M{
		"name":  userData.NAME,
		"img":   userData.IMG,
		"email": userData.EMAIL,
		"_id":   userId,
	}
	if err := coll.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&findData); err == mongo.ErrNoDocuments {
		result, err := coll.InsertOne(context.TODO(), insertData)
		return result, err
	} else if err != nil {
		return nil, err
	}
	for key, value := range insertData {
		if findData[key] != value {
			filter := bson.D{{Key: "_id", Value: userId}}
			update := bson.M{"$set": insertData}
			updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
			return updateResult, err
		}
	}
	return "Data already exist and fresh state", nil
}
