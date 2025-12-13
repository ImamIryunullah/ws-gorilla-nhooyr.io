package main

import (
	"log"
	"net/http"

	ws "websocket-nhooyr/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	router := gin.Default()

	// Route biasa pakai Gin
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "✅ WebSocket server is running 🚀")
	})

	// Buat HTTP multiplexer (gabungkan Gin dan WS)
	mux := http.NewServeMux()
	mux.Handle("/", router) // semua route Gin
	mux.HandleFunc("/nhooyr", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("🚀 Server running on :8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
