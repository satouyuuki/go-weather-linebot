package externalapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const OPENWEATHER_ORIGIN = "https://api.openweathermap.org"

type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetGeoLocation(apiKey string, keyword string, geo interface{}) error {
	requestURL := fmt.Sprintf("%s/geo/1.0/direct?q=%s&limit=1&appid=%s", OPENWEATHER_ORIGIN, keyword, apiKey)
	resp, err := http.Get(requestURL)
	if err == nil {
		georesp, _ := json.Marshal(resp.Body)
		log.Printf("georesp: %s", georesp)
		defer resp.Body.Close()
	} else {
		panic(err)
	}
	return json.NewDecoder(resp.Body).Decode(geo)
}
