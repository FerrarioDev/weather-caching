# WeatherCache 🌤️

A Go application that fetches real-time weather data and caches it in Redis for better performance.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API](#api)
- [Development](#development)
- [Troubleshooting](#troubleshooting)

## Features

- **Fast Caching**: Redis-backed caching with 1-hour expiration
- **Real-time Weather**: Integrates with Visual Crossing Weather API
- **Smart Cache Keys**: Automatic location normalization (London = london = LONDON)
- **Error Handling**: Comprehensive error messages and validation
- **Graceful Shutdown**: Proper resource cleanup
- **Environment Config**: Fully configurable via .env

## Prerequisites

- **Go** 1.25.2 or higher
- **Redis** 6.0+ (running locally or via Docker)
- **Weather API Key** from [Visual Crossing](https://www.visualcrossing.com/)

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/FerrarioDev/weathercache.git
cd weathercache
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment

```bash
# Copy template
cp .env.example .env

# Edit with your API key
nano .env
```

**Required variables in `.env`:**
```env
WEATHER_API=your_api_key_here
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
SERVER_PORT=:8084
```

### 4. Start Redis

```bash
# Option A: If Redis installed
redis-server

# Option B: Using Docker
docker run -d -p 6379:6379 redis:latest
```

### 5. Run Application

```bash
go run ./src/weather/main.go
```

You should see:
```
Starting server on :8084
```

## Configuration

Edit `.env` file to customize:

| Variable | Description | Default |
|----------|-------------|---------|
| `WEATHER_API` | Visual Crossing API key | Required |
| `REDIS_ADDR` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | (empty) |
| `SERVER_PORT` | HTTP server port | `:8084` |

### Getting a Weather API Key

1. Visit [Visual Crossing](https://www.visualcrossing.com/)
2. Sign up for free account
3. Copy API key from dashboard
4. Add to `.env`

## Usage

### Start Server

```bash
go run ./src/weather/main.go
```

### Make Requests

```bash
# Get weather for a location
curl "http://localhost:8084/?location=London"

# Pretty print JSON
curl -s "http://localhost:8084/?location=Paris" | jq

# With spaces (URL encoded)
curl "http://localhost:8084/?location=New%20York"
```

### Stop Server

Press `Ctrl+C` to gracefully shutdown.

## API

### Endpoint: `GET /`

Returns current weather for specified location.

**Query Parameters:**
- `location` (required): City name or location (e.g., "London", "Tokyo")

**Success Response (200):**
```json
{
  "location": {
    "name": "London, England, United Kingdom",
    "address": "london",
    "coordinates": {
      "latitude": 51.5074,
      "longitude": -0.1278
    },
    "timezone": "Europe/London"
  },
  "current": {
    "temperature": 8.5,
    "humidity": 72,
    "windSpeed": 12.5,
    "conditions": "Partly cloudy",
    "icon": "partly-cloudy-day"
  },
  "cached": false,
  "lastUpdated": "2026-01-08T14:30:00Z"
}
```

**Error: Missing Location (400)**
```
you need to insert a location
```

**Error: Invalid API Key (401)**
```
api error: status code 401
```

### Cache Behavior

- **First request**: Calls API (~500-1500ms), stores in cache
- **Same location within 1 hour**: Returns cached data (~10-50ms)
- **Location normalization**: "London", "london", "LONDON" all use same cache

## Development

### Project Structure

```
weathercache/
├── config/              # Configuration management
│   └── config.go
├── internal/
│   ├── cache/          # Redis caching
│   │   ├── redis.go
│   │   └── cache_test.go
│   ├── domain/         # Data models
│   │   └── model.go
│   └── service/        # Business logic
│       ├── weather.go
│       └── handler.go
├── src/weather/        # Application entry
│   └── main.go
└── tests/              # Integration tests
    └── weather_test.go
```

### Key Components

**`config/config.go`**
- Manages environment variables
- Validates required configuration
- Provides type-safe access

**`internal/service/weather.go`**
- Calls Visual Crossing API
- Transforms API response to internal format
- Handles API errors

**`internal/service/handler.go`**
- HTTP request handling
- Location normalization
- Cache-first strategy

**`internal/cache/redis.go`**
- Redis connection management
- JSON serialization for cache
- 1-hour TTL on cached entries

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test -v ./internal/cache

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestRedisCache -v ./internal/cache
```

### Code Style

```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...
```

## Troubleshooting

### Redis Connection Error

```
Error: redis: connection refused
```

**Solution:**
```bash
# Start Redis
redis-server

# Or verify it's running
redis-cli ping  # Should return PONG
```

### Missing API Key Error

```
Error: required configuration missing: WEATHER_API
```

**Solution:**
```bash
# Check .env file
cat .env

# Add your key
echo 'WEATHER_API=your_key_here' >> .env
```

### Port Already in Use

```
Error: listen tcp :8084: bind: address already in use
```

**Solution:**
```bash
# Use different port
SERVER_PORT=:8085 go run ./src/weather/main.go

# Or kill process using port
lsof -i :8084
kill -9 <PID>
```

### API Rate Limit

```
Error: api error: status code 429
```

**Solution:**
- Cache is working! Most requests should hit cache
- Wait before making more API requests
- Upgrade Visual Crossing plan for higher limits

### First Request Slow

This is normal! First request calls the API:
1. Makes HTTP request to Visual Crossing
2. Parses response
3. Stores in Redis
4. Returns to client (~1 second total)

Subsequent requests within 1 hour are instant (from cache).

## Performance Tips

- **Reuse locations**: Same location twice benefits from cache
- **Check cache hit rate**: Monitor if caching is working
- **Monitor Redis**: Use `redis-cli INFO memory` to check memory
- **Batch requests**: Make many requests at once to maximize cache

## Security

- ✅ `.env` is in `.gitignore` - never committed
- ✅ API key not logged or in error messages
- ✅ Use `.env.example` as template
- ✅ Rotate API keys periodically

## Testing with cURL

```bash
# First request (hits API)
time curl "http://localhost:8084/?location=London"

# Second request (from cache) - much faster!
time curl "http://localhost:8084/?location=London"

# Different location normalization (same cache)
curl "http://localhost:8084/?location=london"
curl "http://localhost:8084/?location=LONDON"
curl "http://localhost:8084/?location=london,"
```

## Testing with Python

```python
import requests

response = requests.get(
    "http://localhost:8084/",
    params={"location": "London"}
)

weather = response.json()
print(f"Temperature: {weather['current']['temperature']}°C")
print(f"Conditions: {weather['current']['conditions']}")
print(f"From cache: {weather['cached']}")
```

## Testing with JavaScript

```javascript
const location = "London";

fetch(`http://localhost:8084/?location=${location}`)
  .then(res => res.json())
  .then(data => {
    console.log(`Temperature: ${data.current.temperature}°C`);
    console.log(`Conditions: ${data.current.conditions}`);
  });
```

## Monitoring Cache

```bash
# Connect to Redis CLI
redis-cli

# List all cache keys
KEYS *

# Get specific cache entry
GET london

# Check memory usage
INFO memory

# Clear all cache
FLUSHALL  # ⚠️  WARNING: Deletes all data!
```

## Environment-Specific Setup

### Local Development

```env
WEATHER_API=your_dev_key
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
SERVER_PORT=:8084
```

### Docker Development

```env
WEATHER_API=your_key
REDIS_ADDR=redis:6379
REDIS_PASSWORD=your_password
SERVER_PORT=:8084
```

### Production

```env
WEATHER_API=your_prod_key
REDIS_ADDR=redis.prod.example.com:6379
REDIS_PASSWORD=secure_password
SERVER_PORT=:8080
```

## Architecture Overview

```
Client Request (GET /?location=London)
  ↓
Handler validates & normalizes location
  ↓
Check Redis cache for "london" key
  ├─ Cache HIT → Return data (~50ms)
  └─ Cache MISS → Continue
    ↓
  Call Visual Crossing API
    ↓
  Parse & validate response
    ↓
  Store in Redis (1 hour TTL)
    ↓
  Return to client (~1000ms for first request)
```

## Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/name`
3. Make changes and test: `go test ./...`
4. Commit: `git commit -m "Description"`
5. Push and create Pull Request

## Support

For issues:
1. Check Troubleshooting section
2. Search existing GitHub issues
3. Create new issue with details:
   - Error message
   - Steps to reproduce
   - Go version
   - OS/Platform

## License

Open source - see LICENSE file

## Author

FerrarioDev

---

**Last Updated**: January 8, 2026
