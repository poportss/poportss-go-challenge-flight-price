# ✈️ Go Flight Price Challenge

A concurrent, secure, and extensible **Go microservice** that aggregates **flight offers** from multiple providers — **Amadeus**, **Google Flights (SerpAPI)**, **AirScraper (RapidAPI)**, and a **Mock provider** — returning the **cheapest**, **fastest**, and **comparable offers** for a given route.

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
cmd/api/main.go
├── internal/
│   ├── domain/           # Core models and DTOs
│   ├── flights/          # Business logic, cache, and aggregator service
│   ├── providers/        # Integrations with external APIs (Amadeus, Google, AirScraper, Mock)
│   ├── http/             # Gin-based HTTP server and controllers
│   │   ├── controllers/  # Auth, Flights, SSE controllers
│   │   ├── middleware/   # JWT authentication and provider injection
│   ├── util/             # HTTP utilities and env helpers
│
└── test/                 # E2E and integration tests
```

---

## ⚙️ Features

✅ **Parallel API requests** – all providers are queried concurrently using `errgroup`.  
✅ **JWT Authentication** – required for accessing `/flights/*` routes.  
✅ **Dynamic Provider Registration** – providers are registered after successful login.  
✅ **Amadeus OAuth2 Integration** – retrieves and uses real access tokens.  
✅ **Google Flights via SerpAPI** – fetches flight data through the SerpAPI integration.  
✅ **AirScraper (RapidAPI)** – uses the RapidAPI flight search endpoint.  
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
  "providers": ["GoogleFlights", "Amadeus", "Mock"]
}
```

---

### ✈️ `GET /flights/search`
Searches for flight offers from all registered providers concurrently.

#### Query Parameters:
| Parameter | Type | Required | Example |
|------------|------|-----------|----------|
| origin | string | ✅ | GRU |
| destination | string | ✅ | JFK |
| startDate | date | ✅ | 2025-12-01 |
| endDate | date | ❌ | 2025-12-10 |

#### Response:
```json
{
  "cheapest": {
    "provider": "Amadeus",
    "airline": "LATAM",
    "price": 355.34,
    "currency": "EUR"
  },
  "fastest": {
    "provider": "GoogleFlights",
    "airline": "United",
    "price": 400.00,
    "currency": "USD"
  },
  "offers": [
    { "provider": "Amadeus", ... },
    { "provider": "GoogleFlights", ... },
    { "provider": "Mock", ... }
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
    {"month": "2025-10", "avgPrice": 700, "currency": "USD"},
    {"month": "2025-09", "avgPrice": 720, "currency": "USD"}
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

| Layer | Technology |
|--------|-------------|
| Language | Go 1.22+ |
| Framework | Gin Gonic |
| HTTP Client | Native `net/http` |
| Auth | JWT Middleware |
| Concurrency | `errgroup` |
| Caching | In-memory TTL store |
| APIs | Amadeus, SerpAPI (Google Flights), AirScraper, Mock |
| Streaming | SSE |
| Testing | `testing` + `httptest` |

---

## ⚙️ Environment Variables

| Variable | Description | Example |
|-----------|--------------|----------|
| `PORT` | HTTP port | `8080` |
| `JWT_SECRET` | JWT signing secret | `devsecret` |
| `AMADEUS_CLIENT_ID` | Amadeus API client ID | `abc123` |
| `AMADEUS_CLIENT_SECRET` | Amadeus API client secret | `xyz456` |
| `SERP_API_GOOGLEFLIGHTS_API_KEY` | SerpAPI key for Google Flights | `your_serpapi_key` |
| `RAPIDAPI_AIRSCRAPER_KEY` | RapidAPI key for AirScraper | `your_rapidapi_key` |

---

## 🧠 Caching Behavior
Each `SearchRequest` result is cached for **30 seconds** using an in-memory TTL cache. Expired entries are automatically cleaned up by a background goroutine.

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
