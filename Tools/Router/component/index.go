package _component

import (
	"bytes"
	"context"
	"encoding/json"
	"fiber/Tools/mongodb"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TPostData struct { // * Create
	NAME  string
	HTML  string
	STYLE string
	MAKER string
}

type TPutData struct { // * Update
	ID    string
	NAME  string
	HTML  string
	STYLE string
}

type TDeleteData struct { // * Delete
	ID    string
	MAKER string
}

var Get = func(c *fiber.Ctx) error {
	keyword := strings.ToLower(c.Query("keyword"))
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	skip, _ := strconv.ParseInt(c.Query("skip"), 10, 32)
	data, err := getData(keyword, limit, skip)

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

var Test = func() string {
	const URL = "http://127.0.0.1:5000/component"
	postData := map[string]interface{}{
		"name":  strconv.Itoa(rand.Int()),
		"html":  strconv.Itoa(rand.Int()),
		"style": strconv.Itoa(rand.Int()),
		"maker": strconv.Itoa(rand.Int()),
	}
	pbytes, _ := json.Marshal(postData)
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", URL, buff)
	if err != nil {
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	var postResult map[string]interface{}
	if err := json.Unmarshal(respBody, &postResult); err != nil {
		return err.Error()
	}

	insertedID := postResult["InsertedID"]
	if insertedID == "" || insertedID == nil {
		return "/component POST error"
	}

	putData := map[string]interface{}{
		"name":  strconv.Itoa(rand.Int()),
		"html":  strconv.Itoa(rand.Int()),
		"style": strconv.Itoa(rand.Int()),
		"id":    insertedID,
	}
	pbytes, _ = json.Marshal(putData)
	buff = bytes.NewBuffer(pbytes)

	req, err = http.NewRequest("PUT", URL, buff)
	if err != nil {
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err.Error()
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	var putResult map[string]interface{}
	if err := json.Unmarshal(respBody, &putResult); err != nil {
		return err.Error()
	}
	if putResult["ModifiedCount"] == nil || putResult["ModifiedCount"].(float64) != 1 {
		return "/component PUT error"
	}
	resp, err = http.Get(URL + "?keyword=" + putData["name"].(string) + "&limit=0&skip=0")
	fmt.Println(URL + "?keyword=" + putData["name"].(string) + "&limit=0&skip=0")
	if err != nil {
		return err.Error()
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	var getResult []map[string]interface{}
	if err := json.Unmarshal(respBody, &getResult); err != nil {
		return err.Error()
	}
	if getResult[0]["_id"].(string) != insertedID {
		return "/component GET error"
	}

	deleteData := map[string]interface{}{
		"id":    getResult[0]["_id"],
		"maker": getResult[0]["maker"],
	}
	pbytes, _ = json.Marshal(deleteData)
	buff = bytes.NewBuffer(pbytes)

	req, err = http.NewRequest("DELETE", URL, buff)
	if err != nil {
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		return err.Error()
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	var deleteResult map[string]interface{}
	if err := json.Unmarshal(respBody, &deleteResult); err != nil {
		return err.Error()
	}
	if deleteResult["DeletedCount"] == nil || deleteResult["DeletedCount"].(float64) != 1 {
		return "/component DELETE error"
	}
	return "OK: /component test"
}

func getData(keyword string, limit int64, skip int64) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	filter := bson.M{"keyword": bson.M{"$regex": keyword}}
	opts := options.Find().SetSort(bson.M{"likeCount": -1}).SetLimit(limit).SetSkip(skip)

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	dataArray := []map[string]interface{}{}
	for _, result := range results {
		dataMap := map[string]interface{}{}
		for _, k := range result {
			dataMap[k.Key] = k.Value
		}
		dataArray = append(dataArray, dataMap)
	}
	return dataArray, err
}

func postData(data TPostData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	insertData := bson.M{
		"name":      data.NAME,
		"keyword":   strings.ToLower(data.NAME),
		"html":      data.HTML,
		"style":     data.STYLE,
		"maker":     data.MAKER,
		"like":      [0]string{},
		"likeCount": 0,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}

func putData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	updateData := bson.M{"updatedAt": time.Now()}
	values := reflect.ValueOf(data)
	for i := 0; i < values.NumField(); i++ {
		dataName := strings.ToLower(values.Type().Field(i).Name)
		if values.Field(i).Interface() != "" && dataName != "id" {
			if dataName == "name" {
				updateData["keyword"] = strings.ToLower(values.Field(i).String())
			}
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

func deleteData(data TDeleteData) (*mongo.DeleteResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": [2]primitive.M{bson.M{"_id": id}, bson.M{"maker": data.MAKER}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	return result, err
}
