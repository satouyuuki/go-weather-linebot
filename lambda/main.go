package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"lambda.GoWeatherLinebot/constant"
	"lambda.GoWeatherLinebot/externalapi"
	util "lambda.GoWeatherLinebot/util"
)

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// debug
	lineEventsJson, _ := json.Marshal(req)
	log.Printf("lineEventsJson: %s", lineEventsJson)

	// set env
	wApiToken := os.Getenv("OPENWEATHER_API_TOKEN")
	cSecret := os.Getenv("CHANNNE_SECRET")
	cToken := os.Getenv("CHANNNE_TOKEN")

	// create linebot
	bot, err := linebot.New(cSecret, cToken)
	if err != nil {
		log.Fatal(err)
	}

	// bloa
	if req.Headers == nil {
		return externalapi.Broadcast("渋谷区", wApiToken, bot)
	}

	lineEvents, err := util.ParseRequest(cSecret, req)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "" {
					if _, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(constant.MESSAGE_NOT_FOUND)).Do(); err != nil {
						log.Print(err)
					}
					return events.LambdaFunctionURLResponse{Body: constant.MESSAGE_NOT_FOUND, StatusCode: http.StatusBadRequest}, nil
				}
				return externalapi.Reply(message.Text, wApiToken, event.ReplyToken, bot)
			default:
				if _, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(constant.INVALID_TYPE_MESSAGE)).Do(); err != nil {
					log.Print(err)
				}
				return events.LambdaFunctionURLResponse{Body: constant.INVALID_TYPE_MESSAGE, StatusCode: http.StatusBadRequest}, nil
			}
		}
	}
	return events.LambdaFunctionURLResponse{
		Body: constant.FAILED_READ_MESSAGE, StatusCode: http.StatusInternalServerError}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
