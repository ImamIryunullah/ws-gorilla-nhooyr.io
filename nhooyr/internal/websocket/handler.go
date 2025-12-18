package websocket

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionContextTakeover, // ✅ compress messages
	})
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 512), // ✅ buffer lebih besar
	}

	hub.register <- client

	ctx := context.Background()
	go client.WritePump(ctx)
	go client.ReadPump(ctx)
}
