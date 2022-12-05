package externalapi

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"lambda.GoWeatherLinebot/constant"
	util "lambda.GoWeatherLinebot/util"
)

func Reply(msg, wApiToken, replyToken string, bot *linebot.Client) (events.LambdaFunctionURLResponse, error) {
	geo := new([]GeoLocation)
	if err := GetGeoLocation(
		wApiToken,
		msg,
		geo,
	); err != nil {
		log.Println(err)
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(constant.GEOLOCATION_API_EXEC_FAIL)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
	} else if len(*geo) == 0 {
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(constant.GEOLOCATION_API_NOT_FOUND)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_NOT_FOUND, StatusCode: http.StatusBadRequest}, nil
	}
	weather := new(OneCall)
	if err := GetWeather(
		wApiToken,
		(*geo)[0].Lat,
		(*geo)[0].Lon,
		weather,
	); err != nil {
		log.Println(err)
		if _, err = bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(constant.WEATHER_API_EXEC_FAIL)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.WEATHER_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
	}

	var hourRains []string
	for i, hour := range weather.Hourly {
		if i >= 15 {
			break
		}
		log.Printf("time: %s, main: %s", util.ToJstFromTimestamp(hour.Dt).Format(time.RFC3339), hour.Weather[0].Main)
		wType := constant.ParseWeatherType(hour.Weather[0].Main)
		if constant.NeedUmbrella(wType) {
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
		msg,
		strings.Join(hourRains, "\n"),
		weather.Current.Weather[0].Description,
		weather.Current.FeelsLike,
	)
	// レスポンスメッセージ
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(sendmessage)).Do(); err != nil {
		log.Print(err)
	}
	return events.LambdaFunctionURLResponse{Body: sendmessage, StatusCode: 200}, nil
}

func Broadcast(msg, wApiToken string, bot *linebot.Client) (events.LambdaFunctionURLResponse, error) {
	geo := new([]GeoLocation)
	if err := GetGeoLocation(
		wApiToken,
		msg,
		geo,
	); err != nil {
		log.Println(err)
		if _, err = bot.BroadcastMessage(
			linebot.NewTextMessage(constant.GEOLOCATION_API_EXEC_FAIL)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
	} else if len(*geo) == 0 {
		if _, err = bot.BroadcastMessage(
			linebot.NewTextMessage(constant.GEOLOCATION_API_NOT_FOUND)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.GEOLOCATION_API_NOT_FOUND, StatusCode: http.StatusBadRequest}, nil
	}
	weather := new(OneCall)
	if err := GetWeather(
		wApiToken,
		(*geo)[0].Lat,
		(*geo)[0].Lon,
		weather,
	); err != nil {
		log.Println(err)
		if _, err = bot.BroadcastMessage(
			linebot.NewTextMessage(constant.WEATHER_API_EXEC_FAIL)).Do(); err != nil {
			log.Print(err)
		}
		return events.LambdaFunctionURLResponse{Body: constant.WEATHER_API_EXEC_FAIL, StatusCode: http.StatusInternalServerError}, nil
	}

	var hourRains []string
	for i, hour := range weather.Hourly {
		if i >= 15 {
			break
		}
		log.Printf("time: %s, main: %s", util.ToJstFromTimestamp(hour.Dt).Format(time.RFC3339), hour.Weather[0].Main)
		wType := constant.ParseWeatherType(hour.Weather[0].Main)
		if constant.NeedUmbrella(wType) {
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
		msg,
		strings.Join(hourRains, "\n"),
		weather.Current.Weather[0].Description,
		weather.Current.FeelsLike,
	)
	// レスポンスメッセージ
	if _, err := bot.BroadcastMessage(
		linebot.NewTextMessage(sendmessage)).Do(); err != nil {
		log.Print(err)
	}
	return events.LambdaFunctionURLResponse{Body: sendmessage, StatusCode: 200}, nil
}
