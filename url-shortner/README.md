# URL Shortener with Analytics Backend 
## Overview 
This is a backend service for a URL shortener that also tracks analytics data. It uses Go, Gin, GORM, PostgreSQL, and Redis to provide the following functionalities: 

    URL Shortening : Convert long URLs into short, unique slugs.
    URL Retrieval : Redirect users from short slugs back to the original URLs.
    Analytics Logging : Track and log each URL visit, including IP address, device type, and geolocation.

## Features

- URL Shortening
    - POST `/shorten` with a JSON body containing the original URL.
    - Returns a JSON response with a unique slug.
- URL Retrieval
    - GET `/r/:slug` redirects to the original URL.
- Analytics Logging
    - Logs visit details, including IP address, device type, and geolocation.
- Rate Limiting
    - Limits the number of requests per minute per IP.
- Background Processing
    - Uses Redis to queue and process analytics data asynchronously.
- Real-Time Analytics
    - Uses Redis Pub/Sub to broadcast analytics updates to clients.

## Folder Structure
```bash
url-shortener/
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── handlers/
│   │   ├── url_handler.go
│   │   └── analytics_handler.go
│   ├── models/
│   │   ├── models.go
│   ├── services/
│   │   ├── url_service.go
│   │   └── analytics_service.go
│   └── middleware/
│       └── rate_limiter.go
├── pkg/
│   ├── redis/
│   │   └── redis_client.go
│   └── postgres/
│       └── postgres_client.go
├── .env
└── go.mod
```