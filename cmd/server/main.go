package main

import (
    "fmt"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"

    "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func main() {
    fmt.Println("Starting Peril server...")

    gamelogic.PrintClientHelp()

    for {
        userInput := gamelogic.GetInput()
        if len(userInput) == 0 {
            continue
        }

        const rabbitConnString = "amqp://guest:guest@localhost:5672/"
        conn, err := amqp.Dial(rabbitConnString)

        if err != nil {
            log.Fatalf("could not connect to RabbitMQ: %v", err)
        }

        defer conn.Close()
        fmt.Println("connection succeeded")

        newCh, err := conn.Channel()

        switch firstWordUserInput := userInput[0]; firstWordUserInput {
        case "pause":
            fmt.Println("sending pause message")
            pubsub.PublishJSON(newCh, routing.ExchangePerilTopic, routing.GameLogSlug, routing.PlayingState{IsPaused: true}) 
        case "resume":
            fmt.Println("sending resume message")
            pubsub.PublishJSON(newCh, routing.ExchangePerilTopic, routing.GameLogSlug, routing.PlayingState{IsPaused: false}) 
        case "quit":
            fmt.Println("exiting...")
            return
        default:
            fmt.Println("could not find valid command as first argument")
        }

        if err != nil {
            log.Fatalf("could not create new channel: %v", err)
        }
    }
}
