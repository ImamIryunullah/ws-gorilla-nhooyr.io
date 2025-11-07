package websocket

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "closed")
	}()

	for {
		_, message, err := c.conn.Read(ctx)
		if err != nil {
			log.Println("read:", err)
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump(ctx context.Context) {
	defer c.conn.Close(websocket.StatusNormalClosure, "closed")

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			err := c.conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
