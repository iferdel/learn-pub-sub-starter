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

    // create channel for publish
    publishCh, err := conn.Channel()
    if err != nil {
        log.Fatalf("could not create publish channel: %v", err)
    }

    username, err := gamelogic.ClientWelcome()
    if err != nil {
        log.Fatalf("error on specifying username: %v", err)
    }

    newGameState := gamelogic.NewGameState(username)

    // subscribe to army move queue
    err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+newGameState.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(newGameState),
    )
    if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

    // subscribe to pause queue
    err = pubsub.SubscribeJSON(
        conn, 
        routing.ExchangePerilDirect, 
        routing.PauseKey+"."+newGameState.GetUsername(),
        routing.PauseKey,
        pubsub.SimpleQueueTransient,
        handlerPause(newGameState),
    )
    if err != nil {
        log.Fatalf("could not subscribe to pause: %v", err)
    }

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
            mv, err := newGameState.CommandMove(userInput)
            if err != nil {
                fmt.Println(err)
                continue
            }
            err = pubsub.PublishJSON(
                publishCh,
                routing.ExchangePerilTopic, 
                routing.ArmyMovesPrefix+"."+mv.Player.Username, 
                mv,
            )
            if err != nil {
                fmt.Printf("error: %s\n", err)
                continue
            }
            fmt.Printf("Moved %v units to %s\n", len(mv.Units), mv.ToLocation)
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

