package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/katiasuya/audio-conversion-service/internal/converter"
	"github.com/katiasuya/audio-conversion-service/internal/logger"
	"github.com/streadway/amqp"
)

// QueueManager provides methods to use queuing fot conversion requests.
type QueueManager struct {
	name      string
	ch        *amqp.Channel
	converter *converter.Converter
}

// New creates a new queue manager.
func New(name string, ch *amqp.Channel, converter *converter.Converter) *QueueManager {
	return &QueueManager{
		name:      name,
		ch:        ch,
		converter: converter,
	}
}

//ProcessMsgs processes messages coming from the queue, i.e conversion requests.
func (qm *QueueManager) ProcessMsgs() error {
	err := qm.ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("can't set QoS: %w", err)
	}

	msgs, err := qm.ch.Consume(qm.name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("can't register a consumer: %w", err)
	}

	consumeMsgs := make(chan bool)

	type conversionData struct {
		FileID       string
		Filename     string
		SourceFormat string
		TargetFormat string
		RequestID    string
	}

	go func() {
		for d := range msgs {
			var data conversionData
			if err := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&data); err != nil {
				logger.Error(context.Background(), fmt.Errorf("can't decode message: %w", err))
			}
			err := qm.converter.Process(data.FileID, data.Filename, data.SourceFormat, data.TargetFormat, data.RequestID)
			if err != nil {
				logger.Error(context.Background(), err)
			}

			d.Ack(false)
		}
	}()

	<-consumeMsgs

	return nil
}
