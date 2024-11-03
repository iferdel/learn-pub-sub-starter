package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)


type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

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

    err = newCh.QueueBind(queue.Name, key, exchange, false, nil)
    if err != nil {
        return nil, amqp.Queue{}, fmt.Errorf("error while binding queue: %v", err)
    }

    return newCh, queue, nil
}

