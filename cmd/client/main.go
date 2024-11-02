package main

import ( 
    "fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
    "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")

    const rabbitConnString = "amqp://guest:guest@localhost:5672/"
    conn, err := amqp.Dial(rabbitConnString)
    
    if err != nil {
        log.Fatalf("could not connect to RabbitMQ: %v", err)
    }
    
    defer conn.Close()
    fmt.Println("connection succeeded")

    username, err := gamelogic.ClientWelcome()
    if err != nil {
        log.Fatalf("error on specifying username: %v", err)
    }

    _, queue, err := pubsub.DeclareAndBind(
        conn, 
        routing.ExchangePerilDirect, 
        routing.PauseKey+"."+username, 
        routing.PauseKey, 
        pubsub.SimpleQueueTransient,
    ) 
    if err != nil {
        log.Fatalf("could not subscribe to pause: %v", err)
    }
    fmt.Printf("Queue %v declared and bound!\n", queue.Name)

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    <- signalChan
    fmt.Println("Program shutting down, closing connection")
}
