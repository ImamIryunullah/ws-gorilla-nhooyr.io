package websocket

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 1024), // sesudah optimasi: buffered channel untuk throughput tinggi
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// sebelum optimasi
				// select {
				// case client.send <- message:
				// default:
				// 	close(client.send)
				// 	delete(h.clients, client)
				// }

				// sesudah optimasi: drop message jika client overload
				select {
				case client.send <- message:
				default:
					// drop message untuk mencegah block & CPU spike
				}
			}
		}
	}
}
