// Package queue provides methods to work with request queuing.
package queue

import (
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/config"
	"github.com/streadway/amqp"
)

// NewRabbitMQClient creates new rabbitmq connection.
func NewRabbitMQClient(conf *config.RabbitMQData) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(conf.URI)
	if err != nil {
		return nil, nil, fmt.Errorf("can't connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("can't open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(conf.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("can't declare a queue: %w", err)
	}

	return conn, ch, nil
}
