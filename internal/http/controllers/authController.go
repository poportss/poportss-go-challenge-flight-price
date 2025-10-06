package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
	"github.com/poportss/go-challenge-flight-price/internal/http/middleware"
	"github.com/poportss/go-challenge-flight-price/internal/providers"
	"github.com/poportss/go-challenge-flight-price/internal/util"
)

type AuthController struct {
	service   *flights.Service
	jwtSecret string
}

func NewAuthController(service *flights.Service, jwtSecret string) *AuthController {
	return &AuthController{service: service, jwtSecret: jwtSecret}
}

func (a *AuthController) Login(c *gin.Context) {
	var body struct{ User, Pass string }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "bad body"})
		return
	}

	if body.User != "admin" || body.Pass != "secret" {
		c.JSON(401, gin.H{"error": "invalid creds"})
		return
	}

	ctx := c.Request.Context()
	amadeusToken, err := providers.GetAmadeusAccessToken(ctx)
	if err != nil {
		c.JSON(502, gin.H{"error": "amadeus auth failed", "details": err.Error()})
		return
	}

	client := util.NewHTTPClient(1 * time.Minute)
	a.service.AddProvider(providers.NewAmadeus(client, amadeusToken))
	a.service.AddProvider(providers.NewGoogleFlights(client, os.Getenv("SERP_API_GOOGLEFLIGHTS_API_KEY")))
	a.service.AddProvider(providers.NewMockProvider("Ports Airlines"))

	providerTokens := middleware.ProviderTokens{
		AmadeusToken:     amadeusToken,
		GoogleFlightsKey: os.Getenv("GOOGLE_FLIGHTS_KEY"),
	}

	token, err := middleware.GenerateJWT(a.jwtSecret, body.User, time.Hour, providerTokens)
	if err != nil {
		c.JSON(500, gin.H{"error": "jwt error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt_token":  token,
		"expires_in": 3600,
		"providers":  []string{"Google Flights", "Amadeus", "Ports Airlines"},
	})
}
