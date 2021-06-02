package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/katiasuya/audio-conversion-service/internal/repository"
)

func main() {
	lambda.Start(repository.NewPostgresDB)
}
