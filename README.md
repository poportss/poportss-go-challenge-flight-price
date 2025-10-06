# ✈️ Go Flight Price Challenge

A concurrent, secure, and extensible **Go microservice** that aggregates **flight offers** from multiple providers — **Amadeus**, **Google Flights (
SerpAPI)**, **AirScraper (RapidAPI)**, and a **Mock provider** — returning the **cheapest**, **fastest**, and **comparable offers** for a given route.

This project was built as part of a **technical interview challenge** to demonstrate:

- concurrent API requests,
- proper data structuring,
- JWT authentication,
- server-sent events (SSE),
- caching with TTL, and
- clean modular design.

---

## 🧱 Architecture Overview

```
├── !example.env
├── Dockerfile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── domain  # Core models and DTOs
│   │   ├── amadeus.go
│   │   ├── googleFlights.go
│   │   └── models.go
│   ├── flights # Business logic, cache, and aggregator service
│   │   ├── cache.go
│   │   └── service.go
│   ├── http  # Gin-based HTTP server and controllers
│   │   ├── controllers # Auth, Flights, SSE controllers
│   │   ├── middleware # JWT authentication 
│   │   └── server.go
│   ├── providers  # Integrations with external APIs (Amadeus, Google, AirScraper, Mock)
│   │   ├── amadeus.go
│   │   ├── googleFlights.go
│   │   ├── interface.go
│   │   └── mock.go
│   └── util  # HTTP utilities and env helpers
│       └── httpClient.go
└── test   # E2E and integration tests
    ├── e2eSearchTest.go
    └── serviceTest.go
```

---

## ⚙️ Features

✅ **Parallel API requests** – all providers are queried concurrently using `errgroup`.  
✅ **JWT Authentication** – required for accessing `/flights/*` routes.  
✅ **Dynamic Provider Registration** – providers are registered after successful login.  
✅ **Amadeus OAuth2 Integration** – retrieves and uses real access tokens.  
✅ **Google Flights via SerpAPI** – fetches flight data through the SerpAPI integration.   
✅ **Mock Provider** – simulates data when external APIs are limited.  
✅ **Server-Sent Events (SSE)** – provides periodic flight updates every 30 seconds.  
✅ **In-memory TTL Cache** – results cached for 30s to reduce API usage.  
✅ **Unit and E2E Tests** – validate endpoints, error handling, and aggregation.

---

## 🚀 Endpoints

### 🔐 `POST /login`

Authenticates the user and retrieves tokens from external providers (e.g., Amadeus).  
Returns a JWT that must be used in subsequent requests.

#### Request:

```json
{
  "user": "admin",
  "pass": "secret"
}
```

#### Response:

```json
{
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR...",
  "expires_in": 3600,
  "providers": [
    "GoogleFlights",
    "Amadeus",
    "Mock"
  ]
}
```

---

### ✈️ `GET /flights/search`

Searches for flight offers from all registered providers concurrently.

#### Query Parameters:

| Parameter   | Type   | Required | Example    |
|-------------|--------|----------|------------|
| origin      | string | ✅        | GRU        |
| destination | string | ✅        | JFK        |
| startDate   | date   | ✅        | 2025-12-01 |
| endDate     | date   | ✅        | 2025-12-10 |

#### Response:

```json
{
  "cheapest": {
    "provider": "Ports Airlines",
    "airline": "MockAir",
    "price": 875,
    "currency": "USD",
    "duration": 28800000000000,
    "departure_at": "2025-12-01T10:00:00Z",
    "arrival_at": "2025-12-01T18:00:00Z",
    "origin": "GRU",
    "destination": "JFK"
  },
  "fastest": {
    "provider": "GoogleFlights",
    "airline": "American",
    "price": 907,
    "currency": "USD",
    "duration": 27600000000000,
    "departure_at": "2025-12-01T22:35:00Z",
    "arrival_at": "2025-12-02T06:15:00Z",
    "origin": "GRU",
    "destination": "JFK"
  },
  "offers": [
    {
      "provider": "Amadeus",
      ...
    },
    {
      "provider": "GoogleFlights",
      ...
    },
    {
      "provider": "Mock",
      ...
    }
  ]
}
```

---

### 📈 `GET /flights/history`

Simulated endpoint returning average monthly flight prices for the past 24 months.

> APIs like Amadeus and SerpAPI do not expose historical price data, so this is a **mocked** endpoint that demonstrates the format and concept.

#### Response:

```json
{
  "origin": "GRU",
  "destination": "JFK",
  "history": [
    {
      "month": "2025-10",
      "avgPrice": 700,
      "currency": "USD"
    },
    {
      "month": "2025-09",
      "avgPrice": 720,
      "currency": "USD"
    }
  ]
}
```

---

### 🔁 `GET /sse/:route`

Provides **real-time flight updates** every 30 seconds via **Server-Sent Events (SSE)**.

#### Route format:

```
/sse/{origin}|{destination}|{startDate}|{endDate?}
```

#### Example:

```
/sse/GRU|JFK|2025-12-01|2025-12-10
```

The server pushes updates like:

```
event: update
data: {"cheapest": {...}, "fastest": {...}, "offers": [...]}
```

Test via:

```bash
curl -N -H "Authorization: Bearer <JWT_TOKEN>"   http://localhost:8080/sse/GRU|JFK|2025-12-01
```

---

## 🧰 Tech Stack

| Layer       | Technology                                          |
|-------------|-----------------------------------------------------|
| Language    | Go 1.22+                                            |
| Framework   | Gin Gonic                                           |
| HTTP Client | Native `net/http`                                   |
| Auth        | JWT Middleware                                      |
| Concurrency | `errgroup`                                          |
| Caching     | In-memory TTL store                                 |
| APIs        | Amadeus, SerpAPI (Google Flights), AirScraper, Mock |
| Streaming   | SSE                                                 |
| Testing     | `testing` + `httptest`                              |

---

## ⚙️ Environment Variables

| Variable                         | Description                    | Example             |
|----------------------------------|--------------------------------|---------------------|
| `PORT`                           | HTTP port                      | `8080`              |
| `JWT_SECRET`                     | JWT signing secret             | `devsecret`         |
| `AMADEUS_BASE_URL`               | Amadeus API base_url           | `http`              |
| `AMADEUS_CLIENT_ID`              | Amadeus API client ID          | `abc123`            |
| `AMADEUS_CLIENT_SECRET`          | Amadeus API client secret      | `xyz456`            |
| `SERP_API_GOOGLEFLIGHTS_API_KEY` | SerpAPI key for Google Flights | `your_serpapi_key`  |
| `RAPIDAPI_AIRSCRAPER_KEY`        | RapidAPI key for AirScraper    | `your_rapidapi_key` |

---

## 🧠 Caching Behavior

Each `SearchRequest` result is cached for **30 seconds** using an in-memory TTL cache. Expired entries are automatically cleaned up by a background
goroutine.

---

## 🧪 Running Tests

```bash
go test -v ./...
```

---

## 🚀 Run the Project

```bash
git clone https://github.com/<your-username>/go-challenge-flight-price.git
cd go-challenge-flight-price
export PORT=8080
export JWT_SECRET=devsecret
go run ./cmd/api
```

Then test the endpoints with Postman or curl.

---

## 👨‍💻 Author

**Rafael Portela**  
Go Developer | Cloud & API Engineering  
📧 [Contact via LinkedIn](https://linkedin.com)

---

## 📄 License

MIT License — free to use, modify, and share.
