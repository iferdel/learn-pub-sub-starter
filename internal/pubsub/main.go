package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	jsonVal, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("not able to marshal value: %v", err)
	}

    publishMsg := amqp.Publishing{
        ContentType: "application/json",
        Body: jsonVal,
    } 

    err = ch.PublishWithContext(context.Background(), exchange, key, false, false,  publishMsg)
    if err != nil {
        return fmt.Errorf("error while publishing content as json: %v", err)
    }

	return nil
}

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
