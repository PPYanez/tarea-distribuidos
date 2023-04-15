package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	ctx    context.Context
)

func init() {
	// initialize the MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	// establish a connection to the MongoDB server
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	// check the connection
	connected := pingDb(client)
	if !connected {
		panic("Failed to connect to database")
	}

}

func pingDb(client *mongo.Client) bool {
	var result bson.M
	if err := client.Database("distribuidos").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return false
	}
	return true
}

func GetClient() *mongo.Client {
	return client
}
