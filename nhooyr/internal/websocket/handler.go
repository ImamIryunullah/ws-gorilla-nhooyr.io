package websocket

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // biar Origin null tidak masalah
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

	// ❗ gunakan context.Background() agar koneksi tetap hidup
	ctx := context.Background()

	go client.WritePump(ctx)
	go client.ReadPump(ctx)
}
