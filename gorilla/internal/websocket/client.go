package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	dropCount int // sesudah optimasi lanjut: track drop message
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

	// sesudah optimasi lanjut
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

	batch := make([][]byte, 0, 10) // sesudah optimasi lanjut: batch message

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// batch message
			batch = append(batch, message)
		case <-ticker.C:
			if len(batch) == 0 {
				c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
				continue
			}

			for _, msg := range batch {
				c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
			batch = batch[:0] // reset batch
		}
	}
}
