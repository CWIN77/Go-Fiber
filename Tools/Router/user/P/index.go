package _user_P

import (
	"github.com/gofiber/fiber/v2"
)

var Get = func(c *fiber.Ctx) error {
	return c.SendString(c.Params("params"))
}

func getData(id string) (primitive.M, error) {
 client := mongodb.GetMongoClient()
 coll := client.Database("hvData").Collection("user")
 var result bson.M
 userid, _ := primitive.ObjectIDFromHex(id)
 filter := bson.M{"_id": userid}
 err := coll.FindOne(context.TODO(), filter).Decode(&result)
 return result, err
}
