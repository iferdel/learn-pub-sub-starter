package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"

    "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")

    const rabbitConnString = "amqp://guest:guest@localhost:5672/"
    conn, err := amqp.Dial(rabbitConnString)
    
    if err != nil {
        log.Fatalf("could not connect to RabbitMQ: %v", err)
    }
    
    defer conn.Close()
    fmt.Println("connection succeeded")

    newCh, err := conn.Channel()

    pubsub.PublishJSON(newCh,routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}) 

    if err != nil {
        log.Fatalf("could not create new channel: %v", err)
    }

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    <- signalChan
    fmt.Println("Program shutting down, closing connection")
}
