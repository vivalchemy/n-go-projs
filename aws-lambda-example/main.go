package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type EventType struct {
	Age  int    `json:"What is your age?"`
	Name string `json:"What is your name?"`
}

type ResponseType struct {
	StatusCode int
	Body       string `json:"body"`
}

func handleLambdaEvent(event EventType) (ResponseType, error) {
	return ResponseType{
		StatusCode: 200,
		Body:       fmt.Sprintf("Hello %v you are %v years old", event.Name, event.Age),
	}, nil
}

func main() {
	lambda.Start(handleLambdaEvent)
}
