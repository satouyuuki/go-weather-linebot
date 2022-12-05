package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	lineEvents, err := util.ParseRequest("", req)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "" {
					return events.LambdaFunctionURLResponse{Body: constant.MESSAGE_NOT_FOUND, StatusCode: http.StatusBadRequest}, nil
				}
				geo := new([]externalapi.GeoLocation)
				if err = externalapi.GetGeoLocation(
					os.Getenv("OPENWEATHER_API_TOKEN"),
					message.Text,
					geo,
				); err != nil {
					log.Println(err)
					return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
				} else if len(*geo) == 0 {
					return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_NOT_FOUND, StatusCode: http.StatusBadRequest}, nil
				}
				weather := new(externalapi.OneCall)
				if err = externalapi.GetWeather(
					os.Getenv("OPENWEATHER_API_TOKEN"),
					(*geo)[0].Lat,
					(*geo)[0].Lon,
					weather,
				); err != nil {
					log.Println(err)
					return events.LambdaFunctionURLResponse{Body: constant.WEATHER_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
				}

				var hourRains []string
				for i, hour := range weather.Hourly {
					if i >= 15 {
						break
					}
					log.Printf("time: %s, main: %s", util.ToJstFromTimestamp(hour.Dt).Format(time.RFC3339), hour.Weather[0].Main)
					wType := constant.ParseWeatherType(hour.Weather[0].Main)
					if externalapi.NeedUmbrella(wType) {
						hourRains = append(hourRains, fmt.Sprintf("%d時ごろに%s",
							util.ToJstFromTimestamp(hour.Dt).Hour(),
							wType.String(),
						))
					}
				}
				if len(hourRains) > 0 {
					hourRains = append([]string{constant.NEED_UMBRELLA}, hourRains...)
				} else {
					hourRains = append([]string{constant.NO_NEED_UMBRELLA}, hourRains...)
				}

				sendmessage := fmt.Sprintf(
					"今日の%sは%s\n現在の天気は%s\n体感気温は%.1f°",
					message.Text,
					strings.Join(hourRains, "\n"),
					weather.Current.Weather[0].Description,
					weather.Current.FeelsLike,
				)
				// レスポンスメッセージ
				return events.LambdaFunctionURLResponse{Body: sendmessage, StatusCode: 200}, nil
			default:
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
