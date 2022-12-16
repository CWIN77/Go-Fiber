package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetData(client *mongo.Client) (map[string][]primitive.M, error) {
	coll := client.Database("simpleMongo").Collection("Post")
	var result bson.M
	filter := bson.D{{Key: "title", Value: "Join the MongoDB Community"}}
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title")
		return map[string][]primitive.M{"data": {}}, nil
	}
	dataArray := []primitive.M{result}
	return map[string][]primitive.M{"data": dataArray}, err
}

func Connect(uri string) (client *mongo.Client) {
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func CreateTeam(client *mongo.Client)  {
	// db.runCommand( { create: "collection", capped: true, size: 64 * 1024 } )

	db := client.Database("dbName")
	command := bson.D{{Key: "create", Value: "collectionName"}}
	var result bson.M
	if err := db.RunCommand(context.TODO(), command).Decode(&result); err != nil {
		log.Fatal(err)
	}
}