package _user

import (
	"context"
	"fiber/Tools/mongodb"
	"fmt"

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

type UserData struct {
	UID  string `json:"uid" xml:"uid" form:"uid"`
	NAME string `json:"name" xml:"name" form:"name"`
	IMG  string `json:"img" xml:"img" form:"img"`
	MAIL string `json:"mail" xml:"mail" form:"mail"`
}

var Post = func(c *fiber.Ctx) error {
	p := new(UserData)
	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	createData(p)
	return c.Status(200).JSON([4]string{p.UID, p.NAME, p.IMG, p.MAIL})
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

func createData(userData *UserData) {
	if userData.IMG != "" && userData.MAIL != "" && userData.NAME != "" && userData.UID != "" {
		fmt.Println(userData)
	}
}
