package domain

// RawWeatherResponse represents the original API response structure
type RawWeatherResponse struct {
	QueryCost       int     `json:"queryCost"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	ResolvedAddress string  `json:"resolvedAddress"`
	Address         string  `json:"address"`
	Timezone        string  `json:"timezone"`
	// TzOffset        int     `json:"tzoffset"`
	Description string `json:"description"`
	Days        []Day  `json:"days"`
}

type Day struct {
	Datetime       string   `json:"datetime"`
	DatetimeEpoch  int64    `json:"datetimeEpoch"`
	TempMax        float64  `json:"tempmax"`
	TempMin        float64  `json:"tempmin"`
	Temp           float64  `json:"temp"`
	FeelsLikeMax   float64  `json:"feelslikemax"`
	FeelsLikeMin   float64  `json:"feelslikemin"`
	FeelsLike      float64  `json:"feelslike"`
	Dew            float64  `json:"dew"`
	Humidity       float64  `json:"humidity"`
	Precip         float64  `json:"precip"`
	PrecipProb     float64  `json:"precipprob"`
	PrecipCover    float64  `json:"precipcover"`
	PrecipType     []string `json:"preciptype"`
	Snow           float64  `json:"snow"`
	SnowDepth      float64  `json:"snowdepth"`
	WindGust       float64  `json:"windgust"`
	WindSpeed      float64  `json:"windspeed"`
	WindDir        float64  `json:"winddir"`
	Pressure       float64  `json:"pressure"`
	CloudCover     float64  `json:"cloudcover"`
	Visibility     float64  `json:"visibility"`
	SolarRadiation float64  `json:"solarradiation"`
	SolarEnergy    float64  `json:"solarenergy"`
	UVIndex        float64  `json:"uvindex"`
	SevereRisk     float64  `json:"severerisk"`
	Sunrise        string   `json:"sunrise"`
	SunriseEpoch   int64    `json:"sunriseEpoch"`
	Sunset         string   `json:"sunset"`
	SunsetEpoch    int64    `json:"sunsetEpoch"`
	MoonPhase      float64  `json:"moonphase"`
	Conditions     string   `json:"conditions"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	Hours          []Hour   `json:"hours"`
}

type Hour struct {
	Datetime       string   `json:"datetime"`
	DatetimeEpoch  int64    `json:"datetimeEpoch"`
	Temp           float64  `json:"temp"`
	FeelsLike      float64  `json:"feelslike"`
	Humidity       float64  `json:"humidity"`
	Dew            float64  `json:"dew"`
	Precip         float64  `json:"precip"`
	PrecipProb     float64  `json:"precipprob"`
	Snow           float64  `json:"snow"`
	SnowDepth      float64  `json:"snowdepth"`
	PrecipType     []string `json:"preciptype"`
	WindGust       float64  `json:"windgust"`
	WindSpeed      float64  `json:"windspeed"`
	WindDir        float64  `json:"winddir"`
	Pressure       float64  `json:"pressure"`
	Visibility     float64  `json:"visibility"`
	CloudCover     float64  `json:"cloudcover"`
	SolarRadiation float64  `json:"solarradiation"`
	SolarEnergy    float64  `json:"solarenergy"`
	UVIndex        float64  `json:"uvindex"`
	SevereRisk     float64  `json:"severerisk"`
	Conditions     string   `json:"conditions"`
	Icon           string   `json:"icon"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Location struct {
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Country     string      `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
	Timezone    string      `json:"timezone"`
	TzOffset    int         `json:"tzOffset"`
}

type Current struct {
	Datetime       string   `json:"datetime"`
	Timestamp      int64    `json:"timestamp"`
	Temperature    float64  `json:"temperature"`
	FeelsLike      float64  `json:"feelsLike"`
	Humidity       float64  `json:"humidity"`
	DewPoint       float64  `json:"dewPoint"`
	Pressure       float64  `json:"pressure"`
	WindSpeed      float64  `json:"windSpeed"`
	WindGust       float64  `json:"windGust"`
	WindDirection  float64  `json:"windDirection"`
	CloudCover     float64  `json:"cloudCover"`
	Visibility     float64  `json:"visibility"`
	UVIndex        float64  `json:"uvIndex"`
	SolarRadiation float64  `json:"solarRadiation"`
	Precipitation  float64  `json:"precipitation"`
	PrecipProb     float64  `json:"precipProb"`
	PrecipType     []string `json:"precipType"`
	Conditions     string   `json:"conditions"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	Sunrise        string   `json:"sunrise"`
	Sunset         string   `json:"sunset"`
	MoonPhase      float64  `json:"moonPhase"`
}

type TransformedWeatherResponse struct {
	Location    Location `json:"location"`
	Current     Current  `json:"current"`
	Units       string   `json:"units"`
	Cached      bool     `json:"cached"`
	QueryCost   int      `json:"queryCost"`
	LastUpdated string   `json:"lastUpdated"`
}

// {
//   "location": {
//     "name": "London",
//     "country": "UK",
//     "coordinates": {
//       "latitude": 51.5074,
//       "longitude": -0.1278
//     }
//   },
//   "current": {
//     "temperature": 15.5,
//     "feelsLike": 13.2,
//     "humidity": 72,
//     "pressure": 1013,
//     "windSpeed": 12.5,
//     "windDirection": 180,
//     "cloudCover": 45,
//     "visibility": 10,
//     "uvIndex": 3,
//     "description": "Partly cloudy",
//     "icon": "partly-cloudy-day"
//   },
//   "units": "metric",
//   "timestamp": "2026-01-02T14:30:00Z",
//   "cached": false
// }
