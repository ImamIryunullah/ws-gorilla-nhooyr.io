package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	hub.register <- client

	// Context tidak dibatasi 1 jam, supaya broadcast tetap aktif
	ctx := context.Background()

	// Jalankan dua goroutine: write & read
	go client.WritePump(ctx)
	go client.ReadPump(ctx)
}
