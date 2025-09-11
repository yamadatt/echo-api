package main

import (
	"echo-api/internal/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Create a new Lambda handler
	h := handler.NewLambdaHandler()
	
	// Start the Lambda function
	lambda.Start(h.HandleRequest)
}