package constant

type WeatherType int

const (
	Clear = iota + 1
	Clouds
	Rain
	Snow
	Thunderstorm
	Drizzle
	Other
)

var weatherName = [...]string{
	1: "晴れ",
	2: "曇り",
	3: "雨",
	4: "雪",
	5: "雷雨",
	6: "霧雨",
	7: "竜巻など",
}

var weatherValue = map[string]int{
	"Clear":        1,
	"Clouds":       2,
	"Rain":         3,
	"Snow":         4,
	"Thunderstorm": 5,
	"Drizzle":      6,
	"Other":        7,
}

func (w WeatherType) String() string { return weatherName[w] }

func ParseWeatherType(name string) WeatherType {
	for display, num := range weatherValue {
		if name == display {
			return WeatherType(num)
		}
	}
	return Other
}

func NeedUmbrella(w WeatherType) bool {
	return !(w == Clear || w == Clouds)
}
