package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
	"github.com/twistingmercury/api/internal/server"
	"github.com/twistingmercury/heartbeat"
)

var health = flag.Bool("health", false, "Get the current health of the service")

func main() {
	flag.Parse()

	if *health {
		result := server.CheckSelf()
		exCode := 0
		if result.Status != heartbeat.StatusOK {
			exCode = 1
		}
		fmt.Printf("healthcheck: %v\n", result.Status)
		os.Exit(exCode)
	}

	uuidServer := server.NewServer(gin.Default())
	if err := uuidServer.Start(); err != nil {
		fmt.Println(err)
	}
}
