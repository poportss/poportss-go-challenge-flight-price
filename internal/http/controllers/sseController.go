package controllers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
)

type SSEController struct {
	service *flights.Service
}

func NewSSEController(service *flights.Service) *SSEController {
	return &SSEController{service: service}
}

func (s *SSEController) Stream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	route := c.Param("route")
	parts := strings.Split(route, "|")
	if len(parts) < 3 {
		c.SSEvent("error", "invalid route")
		return
	}

	req := domain.SearchRequest{Origin: parts[0], Destination: parts[1]}
	req.StartDate, _ = time.Parse("2006-01-02", parts[2])
	if len(parts) > 3 {
		req.EndDate, _ = time.Parse("2006-01-02", parts[3])
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		resp, err := s.service.Search(c, req)
		if err != nil {
			c.SSEvent("error", err.Error())
		} else {
			c.SSEvent("update", resp)
		}
		c.Writer.Flush()
		<-ticker.C
	}
}
