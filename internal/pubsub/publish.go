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

