package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
	"github.com/katiasuya/audio-conversion-service/internal/storage"
	"github.com/streadway/amqp"
)

var status = []string{"processing", "done", "failed"}

// Converter converts audio files to other formats.
type Converter struct {
	repo    *repository.Repository
	storage *storage.Storage
}

// New creates a new Converter with given fields.
func New(repo *repository.Repository, storage *storage.Storage) *Converter {
	return &Converter{
		repo:    repo,
		storage: storage,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"conversion_requests", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	type conversionData struct {
		FileID       string
		Filename     string
		SourceFormat string
		TargetFormat string
		RequestID    string
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var data conversionData
			if err := json.NewDecoder(bytes.NewReader(d.Body)).Decode(&data); err != nil {
				log.Printf("Error")
			}
			// err := processConversion(data.FileID, data.Filename, data.SourceFormat, data.TargetFormat, data.RequestID)
			// if err != nil {
			// 	log.Printf("Error")
			// }
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c *Converter) processConversion(fileID, filename, sourceFormat, targetFormat, requestID string) error {
	if err := c.repo.UpdateRequest(requestID, status[0], ""); err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	targetFileID, err := uuid.NewRandom()
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't generate target file uuid: %w", err)
	}
	targetFileIDStr := targetFileID.String()

	sourceLocation := fmt.Sprintf(storage.LocationTemplate, fileID, sourceFormat)
	targetLocation := fmt.Sprintf(storage.LocationTemplate, targetFileIDStr, targetFormat)

	cmd := exec.Command("ffmpeg", "-i", sourceLocation, targetLocation)
	if err = cmd.Run(); err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't perform conversion")
	}

	targetFile, err := os.Open(targetLocation)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't generate targetFileID: %w", err)
	}

	err = c.storage.UploadFileToCloud(targetFile, targetFileIDStr, targetFormat)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't upload file to s3: %w", err)
	}

	targetID, err := c.repo.InsertAudio(filename, targetFormat, targetFileIDStr)
	if err != nil {
		if updateErr := c.repo.UpdateRequest(requestID, status[2], ""); updateErr != nil {
			return fmt.Errorf("can't update request: %w", err)
		}
		return fmt.Errorf("can't insert audio: %w", err)
	}

	if err := c.repo.UpdateRequest(requestID, status[1], targetID); err != nil {
		return fmt.Errorf("can't update request: %w", err)
	}

	return nil
}
