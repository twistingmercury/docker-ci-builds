package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/twistingmercury/api/internal/handlers"
	"github.com/twistingmercury/heartbeat"
)

type Server interface {
	Start() error
}

type server struct {
	engine *gin.Engine
}

func NewServer(engine *gin.Engine) Server {
	return &server{engine: engine}
}

func (s *server) Start() error {
	s.setupHealthCheck()
	s.engine.GET("/api/uuid", handlers.GetUUID)

	if err := s.engine.Run(); err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}

func (s *server) setupHealthCheck() {
	dep := heartbeat.DependencyDescriptor{
		Name:        "My custom dependency",
		Type:        "My dependency",
		HandlerFunc: CheckSelf,
	}

	s.engine.GET("/health", heartbeat.Handler("uuid-api", dep))
}

func CheckSelf() heartbeat.StatusResult {
	return heartbeat.StatusResult{
		Status:     heartbeat.StatusOK,
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
}
