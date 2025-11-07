package main

import (
	"log"
	"net/http"

	"websocket-nhooyr/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	hub := websocket.NewHub()
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WebSocket server (nhooyr.io) is running 🚀")
	})

	log.Println("Server running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Server error:", err)
	}
}
