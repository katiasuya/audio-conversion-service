package repository

import (
	"database/sql"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// //NewPostgresDB creates new database connection.
// func NewPostgresDB(c *config.Config) (*sql.DB, error) {
// 	pqInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
// 	db, err := sql.Open("postgres", pqInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, db.Ping()
// }

//NewPostgresDB creates new database connection.
func NewPostgresDB(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pqInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		``, `5432`, `postgres`, ``, `audioconverter`, `disable`)
	db, err := sql.Open("postgres", pqInfo)
	if err != nil {
		fmt.Println("error1", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("ping", err)
	}

	return events.APIGatewayProxyResponse{}, nil
}
