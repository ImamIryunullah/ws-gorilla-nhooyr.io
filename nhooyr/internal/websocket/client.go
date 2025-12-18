package websocket

import (
	"context"
	"log"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// ✅ after optimization
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "client disconnected")
	}()

	for {
		readCtx, cancel := context.WithTimeout(ctx, 60*time.Second) // ✅ read timeout
		_, msg, err := c.conn.Read(readCtx)
		cancel()
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			log.Println("read error:", err)
			return
		}
		select {
		case c.hub.broadcast <- msg: // ✅ non-blocking broadcast
		default:
			// drop message jika hub penuh, jangan close client
		}
	}
}

// ✅ after optimization
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Second) // ✅ heartbeat ping
	defer func() {
		ticker.Stop()
		c.conn.Close(websocket.StatusNormalClosure, "write pump closed")
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.Close(websocket.StatusNormalClosure, "send channel closed")
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				log.Println("write error:", err)
				return
			}
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			c.conn.Ping(pingCtx) // ✅ heartbeat
			cancel()
		case <-ctx.Done():
			return
		}
	}
}
