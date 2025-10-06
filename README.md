# ✈️ Go Flight Price Challenge

A concurrent, secure, and extensible **Go microservice** that aggregates **flight offers** from multiple providers — **Amadeus**, **Google Flights (
SerpAPI)** and a **Mock provider** — returning the **cheapest**, **fastest**, and **comparable offers** for a given route.

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
│   ├── flights # Business logic, cache, and aggregator service
│   ├── http  # Gin-based HTTP server and controllers
│   │   ├── controllers # Auth, Flights, SSE controllers
│   │   ├── middleware # JWT authentication 
│   │   └── server.go
│   ├── providers  # Integrations with external APIs (Amadeus, Google, Mock)
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
    "Google Flights",
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
    "provider": "Google Flights",
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
      "provider": "Google Flights",
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

## 🐳 Run with Docker

### Build the container:
```bash
docker build -t flight-service .
```

### Run with environment variables:
```bash
docker run -p 8080:8080 \
  -e JWT_SECRET=devsecret \
  -e AMADEUS_CLIENT_ID=your_client_id \
  -e AMADEUS_CLIENT_SECRET=your_client_secret \
  -e SERP_API_GOOGLEFLIGHTS_API_KEY=your_serpapi_key \
  flight-service
```

### Or via Docker Compose:
```bash
docker-compose up --build
```

---

## 🔒 HTTPS / TLS Configuration for Production

To enable **secure HTTPS communication** in production, the Flight Price Aggregator can be deployed either behind a **reverse proxy** (recommended) or configured to handle TLS directly.

---

### 🧩 Option 1: Reverse Proxy (Recommended)

Run the Go service internally over HTTP (port `8080`) and use a reverse proxy (Nginx, Caddy, or Traefik) to terminate HTTPS.

**Example Nginx configuration:**
```nginx
server {
    listen 443 ssl;
    server_name api.example.com;

    ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    location / {
        proxy_pass http://flight-service:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

Use [Certbot](https://certbot.eff.org/) to automatically issue and renew SSL certificates.

---

### 🔧 Option 2: Direct TLS in Go (Alternative)
If you prefer to run HTTPS directly from Go:

```bash
export TLS_CERT_FILE=/etc/ssl/certs/server.crt
export TLS_KEY_FILE=/etc/ssl/private/server.key
```

Modify `main.go`:
```go
log.Fatal(server.RunTLS(":443", os.Getenv("TLS_CERT_FILE"), os.Getenv("TLS_KEY_FILE")))
```

---

### 🧪 Option 3: Local Development (Self-Signed)

For local testing:
```bash
openssl req -x509 -newkey rsa:4096 -nodes -keyout server.key -out server.crt -days 365
go run ./cmd/api --tls --cert=server.crt --key=server.key
```

Or with Docker Compose:
```yaml
services:
  flight-service:
    build: .
    environment:
      - TLS_CERT_FILE=/certs/server.crt
      - TLS_KEY_FILE=/certs/server.key
    volumes:
      - ./certs:/certs:ro
```

---

### ✅ Recommendation

- Use **Reverse Proxy (Option 1)** for production deployments.
- Automate renewal with **Certbot** or **Traefik’s built-in Let’s Encrypt integration**.
- Never embed or copy TLS secrets inside Docker images — mount them securely at runtime.

---

## 👨‍💻 Author

**Rafael Portela**  
Go Developer | Cloud & API Engineering  
📧 [LinkedIn Profile](https://www.linkedin.com/in/rafaelportela-dev)

---

## 📄 License

**MIT License** — free to use, modify, and share.
