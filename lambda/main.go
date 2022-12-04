package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"lambda.GoWeatherLinebot/externalapi"
	util "lambda.GoWeatherLinebot/util"
)

func HandleRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
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
				if message.Text == "" {
					return events.LambdaFunctionURLResponse{Body: "メッセージが入力されていません", StatusCode: http.StatusBadRequest}, nil
				}
				geo := new([]externalapi.GeoLocation)
				if err = externalapi.GetGeoLocation(
					os.Getenv("OPENWEATHER_API_TOKEN"),
					message.Text,
					geo,
				); err != nil {
					log.Println(err)
					return events.LambdaFunctionURLResponse{Body: "位置情報の検索に失敗しました", StatusCode: http.StatusBadRequest}, nil
				} else if len(*geo) == 0 {
					return events.LambdaFunctionURLResponse{Body: "正しい都市名を入力してください。例:「新宿区」", StatusCode: http.StatusBadRequest}, nil
				}
				weather := new(externalapi.OneCall)
				if err = externalapi.GetWeather(
					os.Getenv("OPENWEATHER_API_TOKEN"),
					(*geo)[0].Lat,
					(*geo)[0].Lon,
					weather,
				); err != nil {
					log.Println(err)
					return events.LambdaFunctionURLResponse{Body: fmt.Sprint(err), StatusCode: http.StatusInternalServerError}, nil
				}

				var weatherForecast string
				for i, hour := range weather.Hourly {
					if i >= 15 {
						break
					}
					log.Printf("time: %s, main: %s", util.ToJstFromTimestamp(hour.Dt).Format(time.RFC3339), hour.Weather[0].Main)
					if externalapi.NeedUmbrella(hour.Weather[0].Main) {
						weatherForecast += fmt.Sprintf("%d時ごろに雨\n", util.ToJstFromTimestamp(hour.Dt).Hour())
					}
				}
				if weatherForecast == "" {
					weatherForecast = "傘を持っていく必要はない"
				}
				sendmessage := fmt.Sprintf(
					"今日の%sは%sでしょう。\n現在の天気は%s, 体感気温は%.1f°です。",
					message.Text,
					weatherForecast,
					weather.Current.Weather[0].Description,
					weather.Current.FeelsLike,
				)
				// レスポンスメッセージ
				return events.LambdaFunctionURLResponse{Body: sendmessage, StatusCode: 200}, nil
			default:
				return events.LambdaFunctionURLResponse{Body: "テキスト形式で入力してください", StatusCode: http.StatusBadRequest}, nil
			}
		}
	}

	return events.LambdaFunctionURLResponse{
		Body: "lambda end"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
