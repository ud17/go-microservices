package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const(
	WEB_PORT = "80"
	RPC_PORT = "5001"
	MONGO_URL = "mongodb://mongo:27017"
	GRPC_PORT = "50001"
)

var client *mongo.Client

type Config struct {

}

func main() {

	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func ()  {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	} ()
}

func connectToMongo() (*mongo.Client, error) {
	// connection options
	clientOption := options.Client().ApplyURI(MONGO_URL)
	clientOption.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return client, nil
}