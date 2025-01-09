package chat

import (
	"context"
	"time"
)

type Hub struct {
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	Clients    map[*Client]bool
	Uuid       string
}

func NewHub(uuid string) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Uuid:       uuid,
	}
}

func (h *Hub) Run() {
	// Create a context with a 30-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel() // Ensure resources are cleaned up

	// A helper function to reset the context
	resetContext := func() {
		cancel() // Cancel the current context
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	}

	for {
		select {
		case <-ctx.Done():
			// Clear clients in case of any unclean left client
			for client := range h.Clients {
				delete(h.Clients, client)
				close(client.Send)
			}

			// Context timeout expired, stop the goroutine
			return
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				if len(h.Clients) == 0 {
					// Stop the context and exit if no clients remain
					cancel()
					return
				}
			}
		case message := <-h.broadcast:
			// Reset the context on receiving a broadcast
			resetContext()

			for client := range h.Clients {
				select {
				case client.Send <- message:
					// Send message to the client
				default:
					// Delete the client if sending fails
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
