package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/katiasuya/audio-conversion-service/internal/app"
	"github.com/katiasuya/audio-conversion-service/internal/server"
)

// func main() {
// 	http.HandleFunc("/", hello)
// 	log.Fatal(gateway.ListenAndServe(":3000", nil))
// }

// func hello(w http.ResponseWriter, r *http.Request) {
// 	// example retrieving values from the api gateway proxy request context.
// 	requestContext, ok := gateway.RequestContext(r.Context())
// 	if !ok || requestContext.Authorizer["sub"] == nil {
// 		fmt.Fprint(w, "Hello World from Go")
// 		return
// 	}

// 	userID := requestContext.Authorizer["sub"].(string)
// 	fmt.Fprintf(w, "Hello %s from Go", userID)

// }

var err error

func init() {
	server.S, err = app.Run()
	log.Fatalf(err.Error())
}

func main() {
	lambda.Start(server.SignUp)
}
