package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// ======================
// ReadPump
// ======================
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// sebelum optimasi
	// for {
	// 	_, message, err := c.conn.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	c.hub.broadcast <- message
	// }

	// sesudah optimasi
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

// ======================
// WritePump
// ======================
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	// sebelum optimasi
	// defer c.conn.Close()
	// for {
	// 	message, ok := <-c.send
	// 	if !ok {
	// 		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	// 		return
	// 	}
	// 	c.conn.WriteMessage(websocket.TextMessage, message)
	// }

	// sesudah optimasi
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
