package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
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
	Models data.Models
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

	app := Config{
		Models: data.New(client),
	}

	// register to rpc server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	// start server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", RPC_PORT)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", GRPC_PORT))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
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