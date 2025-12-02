package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twistingmercury/heartbeat"
)

func main() {
	start()
}

func start() {
	engine := gin.Default()

	dep := heartbeat.DependencyDescriptor{
		Name:        "My custom dependency",
		Type:        "My dependency",
		HandlerFunc: checkSelf,
	}

	engine.GET("/health", heartbeat.Handler("time-api", dep))

	engine.GET("/api/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time": time.Now().String(),
		})
	})

	if err := engine.Run(); err != nil {
		log.Fatal("failed to run time server: v%", err)
	}
}

func checkSelf() heartbeat.StatusResult {
	return heartbeat.StatusResult{
		Status:     heartbeat.StatusOK,
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
}
