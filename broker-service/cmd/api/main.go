package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const WEB_PORT = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// connect to rabbitmq
	rabbitMQConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1) // exits program immediately without running `defer` functions
	}
	defer rabbitMQConn.Close()
	
	app := Config{
		Rabbit: rabbitMQConn,
	}

	log.Printf("Starting broker service on port %s\n", WEB_PORT)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.Routes(),
	}

	// start server
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
	}

	return connection, nil
}