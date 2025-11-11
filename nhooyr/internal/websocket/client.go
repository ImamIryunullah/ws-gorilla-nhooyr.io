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

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "client disconnected")
		log.Println("🔴 Client disconnected, total:", len(c.hub.clients))
	}()

	for {
		_, msg, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				// normal closure
				return
			}
			log.Println("read error:", err)
			return
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) WritePump(ctx context.Context) {
	defer func() {
		c.conn.Close(websocket.StatusNormalClosure, "write pump closed")
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				// hub closed the channel
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

		case <-ctx.Done():
			return
		}
	}
}
