package test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
	httpserver "github.com/poportss/go-challenge-flight-price/internal/http"
)

func TestSearchEndpoint_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := flights.NewService(nil, 2*time.Second, flights.NewInMemoryTTL())
	s := httpserver.New(svc, "secret")
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/flights/search?origin=GRU&destination=JFK&date=2025-12-01", nil)
	s.Engine().ServeHTTP(w, req)
	if w.Code != 401 {
		t.Fatalf("expected 401")
	}
}
