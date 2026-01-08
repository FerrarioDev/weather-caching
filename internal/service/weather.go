// Package service handles all the api requests and caching
package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/FerrarioDev/weathercache/config"
	"github.com/FerrarioDev/weathercache/internal/domain"
)

const visualCrossingAPI = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

type Weather struct {
	key    string
	client *http.Client
}

func NewWeather() *Weather {
	if err := config.Load(); err != nil {
		fmt.Println(err)
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &Weather{config.WeatherAPI.GetValue(), &client}
}

func (w *Weather) GetWeather(location string) (*domain.TransformedWeatherResponse, error) {
	u := fmt.Sprintf("%s%s", visualCrossingAPI, url.QueryEscape(location))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", w.key)
	q.Add("unitGroup", "metric")
	q.Add("include", "days,current")

	req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api error: status code %d", resp.StatusCode)
	}

	var raw domain.RawWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	weather, err := w.ParseWeather(location, &raw)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (s *Weather) ParseWeather(location string, raw *domain.RawWeatherResponse) (*domain.TransformedWeatherResponse, error) {
	// Validate that we have data
	if len(raw.Days) == 0 {
		return nil, fmt.Errorf("api returned no weather data for the requested location")
	}

	var current domain.Current
	day := raw.Days[0]
	current = domain.Current{
		Datetime:       day.Datetime,
		Timestamp:      day.DatetimeEpoch,
		Temperature:    day.Temp,
		FeelsLike:      day.FeelsLike,
		Humidity:       day.Humidity,
		DewPoint:       day.Dew,
		Pressure:       day.Pressure,
		WindSpeed:      day.WindSpeed,
		WindGust:       day.WindGust,
		WindDirection:  day.WindDir,
		CloudCover:     day.CloudCover,
		Visibility:     day.Visibility,
		UVIndex:        day.UVIndex,
		SolarRadiation: day.SolarRadiation,
		Precipitation:  day.Precip,
		PrecipProb:     day.PrecipProb,
		PrecipType:     day.PrecipType,
		Conditions:     day.Conditions,
		Icon:           day.Icon,
		Description:    raw.Description,
		Sunrise:        day.Sunrise,
		Sunset:         day.Sunset,
		MoonPhase:      day.MoonPhase,
	}

	result := domain.TransformedWeatherResponse{
		Location: domain.Location{
			Name:    raw.ResolvedAddress,
			Address: raw.Address,
			Coordinates: domain.Coordinates{
				Latitude:  raw.Latitude,
				Longitude: raw.Longitude,
			},
			Timezone: raw.Timezone,
		},
		Current:     current,
		QueryCost:   raw.QueryCost,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	return &result, nil
}
