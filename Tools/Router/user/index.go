package _user

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserData struct {
	UID  string `json:"uid" xml:"uid" form:"uid"`
	NAME string `json:"name" xml:"name" form:"name"`
	IMG  string `json:"img" xml:"img" form:"img"`
	MAIL string `json:"mail" xml:"mail" form:"mail"`
}

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := new(UserData)
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if p.IMG == "" || p.MAIL == "" || p.NAME == "" || p.UID == "" {
		return c.Status(400).JSON("Please send all user data.")
	}
	result, err := createData(p)
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

func createData(userData *UserData) (interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	userId, _ := primitive.ObjectIDFromHex(userData.UID)
	var findData bson.M
	insertData := bson.M{
		"name": userData.NAME,
		"img":  userData.IMG,
		"mail": userData.MAIL,
		"_id":  userId,
	}
	if err := coll.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&findData); err == mongo.ErrNoDocuments {
		result, err := coll.InsertOne(context.TODO(), insertData)
		return result, err
	} else if err != nil {
		return nil, err
	}

	if (findData["name"] != insertData["name"]) || (findData["img"] != insertData["img"]) || (findData["mail"] != insertData["mail"]) || (findData["_id"] != insertData["_id"]) {
		filter := bson.D{{Key: "_id", Value: userId}}
		update := bson.M{"$set": insertData}
		
		updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
		return updateResult, err
	}
	return "Data already exist and fresh state", nil
}
