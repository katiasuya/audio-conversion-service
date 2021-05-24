package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/server"
)

var s *server.Server

func Init() {
	log.Fatalln(app.Run())
}

func main() {
	lambda.Start(s.ShowDoc)
}

// }
// func super(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		Body:       "Showing documentation",
// 		StatusCode: 200,
// 	}, nil
// }
