package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MyDate struct {
	Hello string `bson:"hello",json:"hello"`
}

func main() {
	getDate()
}

func getDate() (result []MyDate) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://one:123456@10.28.109.5:27017"))
	collection := client.Database("one").Collection("one_coll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, _ := collection.Find(ctx, bson.M{})
	err := cursor.All(context.TODO(), &result)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success get data")
	return
}
