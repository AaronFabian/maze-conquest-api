package chat

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
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				// fmt.Println("unregister: ", client)
				delete(h.Clients, client)
				close(client.Send)
				if len(h.Clients) == 0 {
					return // Stop goroutine
				}
			}
		case message := <-h.broadcast:
			for client := range h.Clients {
				select {

				// Client send message into all client
				case client.Send <- message:

				// Delete client
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
