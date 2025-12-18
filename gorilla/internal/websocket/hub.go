package websocket

import (
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 1024), // sesudah optimasi lanjut: buffered channel
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	const workerCount = 4 // sesudah optimasi lanjut: worker pool
	workers := make([]chan []byte, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = make(chan []byte, 256)
		go h.broadcastWorker(workers[i])
	}

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			// sesudah optimasi lanjut: kirim ke worker pool
			h.mu.RLock()
			i := 0
			for client := range h.clients {
				select {
				case workers[i%workerCount] <- msg:
				default:
					client.dropCount++
					if client.dropCount > 10 {
						// auto-disconnect client overload
						h.mu.RUnlock()
						h.unregister <- client
						h.mu.RLock()
					}
				}
				i++
			}
			h.mu.RUnlock()
		}
	}
}

// sesudah optimasi lanjut: worker untuk broadcast
func (h *Hub) broadcastWorker(msgCh chan []byte) {
	for msg := range msgCh {
		h.mu.RLock()
		for client := range h.clients {
			select {
			case client.send <- msg:
			default:
				client.dropCount++
				if client.dropCount > 10 {
					h.mu.RUnlock()
					h.unregister <- client
					h.mu.RLock()
				}
			}
		}
		h.mu.RUnlock()
	}
}
