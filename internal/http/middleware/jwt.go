package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ProviderTokens struct {
	AmadeusToken     string `json:"amadeus_token"`
	GoogleFlightsKey string `json:"google_flights_key"`
}

type CustomClaims struct {
	User      string         `json:"user"`
	Providers ProviderTokens `json:"providers"`
	jwt.RegisteredClaims
}

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer"})
			return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		_, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) { return []byte(secret), nil })
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		c.Next()
	}
}

func GenerateJWT(secret, user string, ttl time.Duration, tokens ProviderTokens) (string, error) {
	claims := CustomClaims{
		User:      user,
		Providers: tokens,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
