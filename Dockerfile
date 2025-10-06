# =========================================
# 1️⃣ Build Stage - Compile the Go binary
# =========================================
FROM golang:1.25-alpine AS build

WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o flight-service ./cmd

# =========================================
# 2️⃣ Runtime Stage - Distroless container
# =========================================
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /app/flight-service /app/flight-service

# Just declare that these variables will exist
# (they will be injected by docker-compose via .env)
ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/flight-service"]
