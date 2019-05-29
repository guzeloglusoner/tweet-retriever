package main

// Hub is used for maintaining the active clients and broadcasts messages to all clients
type Hub struct {
	clients    map[*Client]bool
	inbound    chan []byte
	register   chan *Client
	unregister chan *Client
}

func initHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		inbound:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// endlessly listens register / unregister / inbound channels
func (h *Hub) start() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.inbound:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
