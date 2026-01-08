package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FerrarioDev/weathercache/internal/cache"
	"github.com/gorilla/mux"
)

type AppServer struct {
	weather *Weather
	cache   *cache.RedisCache
	server  *http.Server
}

func NewServer(addr string) *AppServer {
	cache := cache.NewCache()
	a := &AppServer{
		weather: NewWeather(),
		cache:   cache,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", a.HandleCurrent).Methods("GET")

	a.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return a
}

// Start starts the HTTP server (blocking)
func (a *AppServer) Start() error {
	fmt.Printf("Starting server on %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server and closes cache connection
func (a *AppServer) Shutdown(ctx interface{}) error {
	if err := a.cache.Shutdown(); err != nil {
		fmt.Printf("warning: failed to close cache: %v\n", err)
	}
	return nil
}

func (s *AppServer) HandleCurrent(w http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "you need to insert a location", http.StatusBadRequest)
		return
	}

	// Normalize location once for consistent cache key
	cacheKey := normalizeLocation(location)

	// Try to get from cache first
	if cachedResponse, err := s.cache.GetValue(req.Context(), cacheKey); err == nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Cache miss - fetch from API
	weather, err := s.weather.GetWeather(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store in cache with normalized key
	if err := s.cache.SetValue(req.Context(), cacheKey, weather); err != nil {
		// Log the error but don't fail the request - cache is not critical
		fmt.Printf("warning: failed to cache weather data: %v\n", err)
	}

	// Return the weather data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func normalizeLocation(location string) string {
	normalized := strings.TrimSpace(location)
	normalized = strings.ReplaceAll(normalized, ",", "")
	normalized = strings.ToLower(normalized)
	return normalized
}
