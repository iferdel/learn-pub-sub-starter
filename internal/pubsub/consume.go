package pubsub

import (
	"encoding/json"
	"fmt"
    "log"

	amqp "github.com/rabbitmq/amqp091-go"
)


type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

type AckType int

const (
    Ack AckType = iota
    NackRequeue 
    NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {

    ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
    if err != nil {
        return fmt.Errorf("could not subscribe to pause: %v", err)
    }
    fmt.Printf("Log queue %v declared and bound!\n", queue.Name)

    msgs, err := ch.Consume(
        queue.Name, //queue 
        "",         //consumer
        false,      // auto-ack
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        //args
    )
    if err != nil {
        return fmt.Errorf("could not consume channel: %v", err)
    }

    unmarshaller := func(data []byte) (T, error) {
            var target T
            err = json.Unmarshal(data, &target)
            return target, err
    }

    go func() {
        defer ch.Close()
        for msg := range msgs {
            dat, err := unmarshaller(msg.Body)
            if err != nil {
                fmt.Printf("could not unmarshal message: %v\n", err)
				continue
            }
            ack := handler(dat)
            switch ack {
            case Ack:
                log.Println("Acknowledging message")
                msg.Ack(false)
            case NackRequeue:
                log.Println("Not Acknowledged message, requeuing...")
                msg.Nack(false, true)
            case NackDiscard:
                log.Println("Not Acknowledged message, discarding...")
                msg.Nack(false, false)
            }
        }
    }()

    return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

    newCh, err := conn.Channel()
    if err != nil {
        return nil, amqp.Queue{}, fmt.Errorf("could not create new channel: %v", err)
    }

    queue, err := newCh.QueueDeclare(
        queueName,                             // name
        simpleQueueType == SimpleQueueDurable, // durable
        simpleQueueType != SimpleQueueDurable, // delete when unused
        simpleQueueType != SimpleQueueDurable, // exclusive
        false,                                 // no-wait
        nil,                                   // arguments
    )
    if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %v", err)
	}

    err = newCh.QueueBind(
        queue.Name, // queue name
        key,        // routing key
        exchange,   // exchange
        false,      // no-wait
        nil,        // args
    )
    if err != nil {
        return nil, amqp.Queue{}, fmt.Errorf("error while binding queue: %v", err)
    }

    return newCh, queue, nil
}

