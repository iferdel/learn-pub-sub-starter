package main

import ( 
    "fmt"
	"log"

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

    newGameState := gamelogic.NewGameState(username)

    for {

        userInput := gamelogic.GetInput()
        if len(userInput) == 0 {
            continue
        }
            
        switch word := userInput[0]; word {
        case "spawn":
            err = newGameState.CommandSpawn(userInput)
            if err != nil {
                fmt.Println(err)
                continue
            }
        case "move":
            _, err := newGameState.CommandMove(userInput)
            if err != nil {
                fmt.Println(err)
                continue
            }
        case "status":
            newGameState.CommandStatus()
        case "help":
            gamelogic.PrintClientHelp()
        case "spam":
            fmt.Println("Spamming not allowed yet!")
        case "quit":
            gamelogic.PrintQuit()
            return
        default:
            fmt.Println("could not find valid command as first argument")
        }
    }
}
