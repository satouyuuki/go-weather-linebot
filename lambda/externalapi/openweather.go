package externalapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lambda.GoWeatherLinebot/constant"
)

const OPENWEATHER_ORIGIN = "https://api.openweathermap.org"

type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WeatherSummary struct {
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string
}

type WeatherDetail struct {
	Dt        int              `json:"dt"`
	Temp      float32          `json:"temp"`
	FeelsLike float32          `json:"feels_like"`
	Weather   []WeatherSummary `json:"weather"`
}

type OneCall struct {
	Timezone string
	Current  WeatherDetail `json:"current"`
	Hourly   []WeatherDetail
}

func GetGeoLocation(apiKey string, keyword string, geo interface{}) error {
	requestURL := fmt.Sprintf("%s/geo/1.0/direct?q=%s&limit=1&appid=%s", OPENWEATHER_ORIGIN, keyword, apiKey)
	resp, err := http.Get(requestURL)
	if err == nil {
		defer resp.Body.Close()
	} else {
		panic(err)
	}
	return json.NewDecoder(resp.Body).Decode(geo)
}

func GetWeather(apiKey string, lat float64, lon float64, weather interface{}) error {
	requestURL := fmt.Sprintf("%s/data/3.0/onecall?lang=ja&exclude=minutely,daily&units=metric&lat=%f&lon=%f&appid=%s", OPENWEATHER_ORIGIN, lat, lon, apiKey)
	resp, err := http.Get(requestURL)
	if err == nil {
		if resp.StatusCode == http.StatusTooManyRequests {
			return errors.New(constant.WEATHER_API_LIMIT_OVER)
		}
		defer resp.Body.Close()
	} else {
		panic(err)
	}
	return json.NewDecoder(resp.Body).Decode(weather)
}
