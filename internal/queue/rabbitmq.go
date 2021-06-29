// Package queue provides methods to work with request queuing.
package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

// NewRabbitMQClient creates new rabbitmq connection.
func NewRabbitMQClient(url, queue string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("can't connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("can't open a channel: %w", err)
	}

	_, err = ch.QueueDeclare("conversion_requests", true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("can't declare a queue: %w", err)
	}

	return conn, ch, nil
}
