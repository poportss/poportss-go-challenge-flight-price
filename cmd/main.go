package main

import (
	"log"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/flights"
	httpserver "github.com/poportss/go-challenge-flight-price/internal/http"
	"github.com/poportss/go-challenge-flight-price/internal/util"
)

func main() {
	log.Println("🚀 Starting Flight Price Aggregator...")

	jwtSecret := util.EnvOr("JWT_SECRET", "devsecret")
	port := util.EnvOr("PORT", "8080")

	// Cache with automatic cleanup every 1 minute
	cache := flights.NewInMemoryTTL()
	cache.StartCleanup(1 * time.Minute)
	log.Println("✓ Cache initialized with automatic cleanup")

	log.Printf("No providers configured yet. They will be set up during login.")

	// Create service with 1-minute timeout and cache
	svc := flights.NewService(nil, 1*time.Minute, cache)
	log.Printf("✓ Service initialized with %d provider(s)", 0)

	// Create and start HTTP server
	server := httpserver.New(svc, jwtSecret)

	log.Printf("🌐 Server running at http://localhost:%s", port)
	log.Printf("📖 Available endpoints:")
	log.Printf("   POST /login - Authentication")
	log.Printf("   GET  /flights/search - Search flights")
	log.Printf("   GET  /flights/history - Flight price history")
	log.Printf("   GET  /sse/:route - Server-Sent Events stream")

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
