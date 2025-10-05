package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
	"github.com/poportss/go-challenge-flight-price/internal/http/controllers"
	"github.com/poportss/go-challenge-flight-price/internal/http/middleware"
)

type Server struct {
	engine  *gin.Engine
	service *flights.Service
}

func New(service *flights.Service, jwtSecret string) *Server {
	r := gin.Default()
	srv := &Server{engine: r, service: service}

	// Controllers
	authCtrl := controllers.NewAuthController(service, jwtSecret)
	flightsCtrl := controllers.NewFlightsController(service)
	sseCtrl := controllers.NewSSEController(service)

	// Rotas públicas
	r.POST("/login", authCtrl.Login)

	// Rotas autenticadas
	auth := r.Group("/", middleware.JWT(jwtSecret))
	auth.GET("/flights/search", flightsCtrl.Search)
	auth.GET("/flights/history", flightsCtrl.History)
	auth.GET("/sse/:route", sseCtrl.Stream)

	return srv
}

func (s *Server) Run(addr string) error { return s.engine.Run(addr) }
func (s *Server) Engine() *gin.Engine   { return s.engine }
