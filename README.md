# ✈️ Go Flight Price Challenge

A concurrent, secure, and extensible **Go microservice** that aggregates **flight offers** from multiple providers — **Amadeus**, **Google Flights (
SerpAPI)**, and a **Mock provider** — returning the **cheapest**, **fastest**, and **comparable offers** for a given route.

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
Each search result is cached for **30 seconds** in memory.  
If the same query is made within the TTL, the cached response is returned immediately.

---

## 🧪 Running Tests
```bash
go test -v ./...
```

---

## 🐳 Running with Docker

The project supports full Docker-based builds using a multi-stage `Dockerfile`.

### 🔹 Build the image
```bash
docker build -t flight-price-service .
```

### 🔹 Run the container
```bash
docker run -p 8080:8080   -e PORT=8080   -e JWT_SECRET=devsecret   -e AMADEUS_BASE_URL=https://test.api.amadeus.com   -e AMADEUS_CLIENT_ID=your_amadeus_client_id   -e AMADEUS_CLIENT_SECRET=your_amadeus_client_secret   -e SERP_API_GOOGLEFLIGHTS_BASE_URL=https://serpapi.com/search   -e SERP_API_GOOGLEFLIGHTS_API_KEY=your_serpapi_key   flight-price-service
```

Then open: [http://localhost:8080](http://localhost:8080)

---

## 🧩 Run with Docker Compose (recommended)

If you prefer to use Docker Compose:

```bash
docker compose up --build
```

The compose file automatically reads variables from your `.env`.

### Example `.env`
```env
PORT=8080
JWT_SECRET=devsecret

AMADEUS_BASE_URL=https://test.api.amadeus.com
AMADEUS_CLIENT_ID=your_amadeus_client_id
AMADEUS_CLIENT_SECRET=your_amadeus_client_secret

SERP_API_GOOGLEFLIGHTS_BASE_URL=https://serpapi.com/search
SERP_API_GOOGLEFLIGHTS_API_KEY=your_serpapi_key
```

---

## 👨‍💻 Author

**Rafael Portela**  
Go Developer | Cloud & API Engineering  
📧 [LinkedIn](https://www.linkedin.com/in/rafaelportela-dev)

---

## 📄 License
MIT License — free to use, modify, and distribute.
