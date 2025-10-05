package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
)

type FlightsController struct {
	service *flights.Service
}

func NewFlightsController(service *flights.Service) *FlightsController {
	return &FlightsController{service: service}
}

func (f *FlightsController) Search(c *gin.Context) {
	var req domain.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := f.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (f *FlightsController) History(c *gin.Context) {
	origin := c.Query("origin")
	dest := c.Query("destination")
	if origin == "" || dest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination required"})
		return
	}

	now := time.Now().UTC()
	history := make([]gin.H, 0, 24)
	for i := 0; i < 24; i++ {
		month := now.AddDate(0, -i, 0)
		history = append(history, gin.H{
			"month":    month.Format("2006-01"),
			"avgPrice": 700 + float64(i)*10 + float64(i%3)*15,
			"currency": "USD",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"origin":      origin,
		"destination": dest,
		"history":     history,
	})
}
