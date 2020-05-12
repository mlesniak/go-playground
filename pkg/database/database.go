package database

import (
	ctx "context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Collection returns the default collection for this application.
func Collection() *mongo.Collection {
	api := client.Database("dev").Collection("api")
	return api

	// c, _ = ctx.WithTimeout(ctx.Background(), 5*time.Second)
	// _, err = api.InsertOne(c, bson.M{"name": "pi", "value": 3.14159})
	// if err != nil {
	// 	panic(err)
	// }
}

func init() {
	c, _ := ctx.WithTimeout(ctx.Background(), 10*time.Second)
	options := options.Client().
		SetAuth(options.Credential{
			Username: "admin",
			Password: "admin",
		}).
		ApplyURI("mongodb://localhost:27017")
	cl, err := mongo.Connect(c, options)
	if err != nil {
		panic(err)
	}
	client = cl
}
