package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// debug log
	eventJson, _ := json.Marshal(req)
	log.Printf("EVENT: %s", eventJson)

	ctxJson, _ := json.Marshal(ctx)
	log.Printf("context: %s", ctxJson)

	return events.LambdaFunctionURLResponse{
		Body: fmt.Sprintf("Hey %s!", eventJson)}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
