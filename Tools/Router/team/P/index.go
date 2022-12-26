package _team_P

import (
	"github.com/gofiber/fiber/v2"
)

var Get = func(c *fiber.Ctx) error {

	return c.SendString(c.Params("params"))
}

func getData(memberId string) ([]map[string]interface{}, error) {
        client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	// 만약 $or가 없다 나오면 {member:{"$or":[...} 
        // 형식으로 변경함
        filter := bson.M{"$or": [
                bson.M{"member.master": memberId}, 
                bson.M{"member.manager": memberId}, 
                bson.M{"member.maker": memberId}, 
                bson.M{"member.reader": memberId}
	]}
	opts := options.Find()
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
