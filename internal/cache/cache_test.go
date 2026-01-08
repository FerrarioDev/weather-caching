package cache

import (
	"context"
	"testing"

	"github.com/FerrarioDev/weathercache/internal/domain"
)

func TestRedisCache(t *testing.T) {
	testWeather := &domain.TransformedWeatherResponse{
		Location: domain.Location{
			Name: "Cordoba",
		},
		Current: domain.Current{
			Temperature: 25.5,
			Humidity:    60,
		},
	}

	r := NewCache()
	defer r.Shutdown()

	if err := r.SetValue(context.Background(), testWeather.Location.Name, testWeather); err != nil {
		t.Errorf("failed to get value from redis")
	}

	got, err := r.GetValue(context.Background(), testWeather.Location.Name)
	if err != nil {
		t.Errorf("failed to get value from redis: %v", err)
	}

	if got.Location.Name != testWeather.Location.Name {
		t.Errorf("location name mismatch: got %s, want %s",
			got.Location.Name, testWeather.Location.Name)
	}

	if got.Current.Temperature != testWeather.Current.Temperature {
		t.Errorf("temperature mismatch: got %f, want %f",
			got.Current.Temperature, testWeather.Current.Temperature)
	}
}
