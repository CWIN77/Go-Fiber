package _team

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Get = func(c *fiber.Ctx) error {
	data, err := getData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
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
								// fmt.Println(memberValue.Map())
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
