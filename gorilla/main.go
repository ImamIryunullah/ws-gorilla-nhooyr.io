package main

import (
	"log"
	"net/http"

	"websocket-go/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("WS Gorilla Optimized | Hub & Pump tuning + Worker Pool")

	router := gin.Default()

	hub := websocket.NewHub()
	go hub.Run()

	router.GET("/gorilla", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WebSocket server is running 🚀")
	})

	log.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server error:", err)
	}
}
