package queue

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

//SendConversionData sends conversion request data to the queue.
func (qm *QueueManager) SendConversionData(fileID, filename, sourceFormat, targetFormat, requestID string) error {
	convData := conversionData{
		FileID:       fileID,
		Filename:     filename,
		SourceFormat: sourceFormat,
		TargetFormat: targetFormat,
		RequestID:    requestID,
	}

	body, err := json.Marshal(convData)
	if err != nil {
		return fmt.Errorf("can't marshal the given payload: %w", err)
	}

	err = qm.ch.Publish("", qm.name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
