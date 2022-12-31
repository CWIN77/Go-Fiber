package _team

import (
	"context"
	"fiber/Tools/mongodb"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPostData struct {
	NAME   string
	MEMBER []map[string]string
}

type TDeleteData struct {
	ID     string
	MASTER string
}

type TPutData struct {
	ID   string
	NAME string
}

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := TPostData{}
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

var Delete = func(c *fiber.Ctx) error {
	p := TDeleteData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	result, err := deleteData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
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

func getData(memberId string) ([]interface{}, error) {
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
	var cursor *mongo.Cursor
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var teamResults []bson.D
	err = cursor.All(context.TODO(), &teamResults)
	if err != nil {
		return nil, err
	}

	coll = client.Database("hvData").Collection("user")
	var memberList []interface{}
	for _, result := range teamResults {
		values := reflect.ValueOf(result[0].Value)
		for i, v := range result {
			if v.Key == "member" {
				values = reflect.ValueOf(result[i].Value)
				break
			}
		}
		for i := 0; i < values.Len(); i++ {
			value := values.Index(i).Interface().(primitive.E).Value.(string)
			if value != memberId {
				inMember := false
				for _, memberListValue := range memberList {
					if value == memberListValue {
						inMember = true
					}
				}
				if !inMember {
					id, err := primitive.ObjectIDFromHex(value)
					if err != nil {
						return nil, err
					}
					memberList = append(memberList, bson.M{"_id": id})
				}
			}
		}
	}
	fmt.Println(memberList)
	filter = bson.M{"$or": memberList}
	cursor, err = coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var memberResults []bson.D
	err = cursor.All(context.TODO(), &memberResults)

	var newTeamList []interface{}
	for _, result := range teamResults {
		newTeam := make(map[string]interface{})
		newMemberList := []interface{}{}
		for i, v := range result {
			if v.Key == "member" {
				values := reflect.ValueOf(result[i].Value)
				for i := 0; i < values.Len(); i++ {
					value := values.Index(i).Interface().(primitive.E).Value.(string)
					id, err := primitive.ObjectIDFromHex(value)
					if err != nil {
						return nil, err
					}
					if value != memberId {
						for _, memberValue := range memberResults {
							if memberValue.Map()["_id"] == id {
								memberClass := values.Index(i).Interface().(primitive.E).Key
								newMemberList = append(newMemberList, map[string]interface{}{memberClass: memberValue.Map()})
							}
						}
					}
				}
			} else {
				newTeam[v.Key] = v.Value
			}
		}
		newTeam["member"] = newMemberList
		newTeamList = append(newTeamList, newTeam)
	}
	return newTeamList, err
}

func postData(data TPostData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	insertData := bson.M{
		"name":   data.NAME,
		"member": data.MEMBER,
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}

func deleteData(data TDeleteData) (*mongo.DeleteResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": [2]primitive.M{bson.M{"_id": id}, bson.M{"member.master": data.MASTER}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	return result, err
}

func putData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	updateData := bson.M{}
	values := reflect.ValueOf(data)
	for i := 0; i < values.NumField(); i++ {
		dataName := strings.ToLower(values.Type().Field(i).Name)
		if values.Field(i).Interface() != "" && dataName != "id" {
			updateData[dataName] = values.Field(i).Interface()
		}
	}
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": updateData}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}
