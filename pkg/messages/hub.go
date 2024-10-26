package messages

import "fmt"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
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
			fmt.Printf("New Message to broadcast: %s\n", message)
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

func parseHTML(message, time []byte) (HTMLmessage []byte) {
	HTMLmessage = ([]byte) fmt.Sprintf(
		`<div id="message" hx-swap-oob="beforeend"><div class="flex justify-end"><div class="bg-gray-300 text-black p-2 rounded-lg max-w-xs">%s: %s%s</div></div></div>`,
		"Test",
		message,
    time,
	)
	return HTMLmessage
}
