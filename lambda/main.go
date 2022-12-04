package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	util "lambda.GoWeatherLinebot/util"
)

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// debug log
	eventJson, _ := json.Marshal(req)
	log.Print("update!!")
	log.Printf("EVENT: %s", eventJson)

	ctxJson, _ := json.Marshal(ctx)
	log.Printf("context: %s", ctxJson)

	lineEvents, err := util.ParseRequest("", req)
	lineEventsJson, _ := json.Marshal(req)
	log.Printf("lineEventsJson: %s", lineEventsJson)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("message: %s", message.Text)
				return events.LambdaFunctionURLResponse{Body: message.Text, StatusCode: 200}, nil
			}
		}
	}

	return events.LambdaFunctionURLResponse{
		Body: fmt.Sprintf("Hey %s!", eventJson)}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
