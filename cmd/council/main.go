package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/pkg/config"
)

func main() {
	fmt.Println("The Council Backend is starting...")

	cfg := config.Load()

	r := gin.Default()

	// Basic health check handler
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	fmt.Printf("Server listening on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
