package main

import (
	"log"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/flights"
	httpserver "github.com/poportss/go-challenge-flight-price/internal/http"
	"github.com/poportss/go-challenge-flight-price/internal/util"
)

func main() {
	log.Println("🚀 Iniciando Flight Price Aggregator...")

	jwtSecret := util.EnvOr("JWT_SECRET", "devsecret")
	port := util.EnvOr("PORT", "8080")

	// Cache com limpeza automática a cada 1 minuto
	cache := flights.NewInMemoryTTL()
	cache.StartCleanup(1 * time.Minute)
	log.Println("✓ Cache inicializado com limpeza automática")

	log.Printf("Nenhum provider configurado. Eles seram configurados no login.")

	// Criar service com timeout de 1 minuto e cache
	svc := flights.NewService(nil, 1*time.Minute, cache)
	log.Printf("✓ Service inicializado com %d provider(s)", 0)

	// Criar e iniciar servidor HTTP
	server := httpserver.New(svc, jwtSecret)

	log.Printf("🌐 Servidor rodando em http://localhost:%s", port)
	log.Printf("📖 Endpoints disponíveis:")
	log.Printf("   POST /login - Autenticação")
	log.Printf("   GET  /flights/search - Buscar voos")
	log.Printf("   GET  /flights/history - Histórico de preços")
	log.Printf("   GET  /sse/:route - Server-Sent Events")

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
