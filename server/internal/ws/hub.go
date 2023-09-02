package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"Clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// Check if room exist
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				// Check if client already in the room or not
				if _, ok := r.Clients[cl.ID]; !ok {
					// If client didnt exist in the room, then add client
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			// Check if client exist in the room or not
			if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
				// Broadcast a message to inform client has left the room
				if len(h.Rooms[cl.RoomID].Clients) != 0 {
					h.Broadcast <- &Message{
						Content:  "user left the chat",
						RoomID:   cl.RoomID,
						Username: cl.Username,
					}
				}

				delete(h.Rooms[cl.RoomID].Clients, cl.ID)
				close(cl.Message)
			}
		case m := <-h.Broadcast:
			// Check if room exist
			if _, ok := h.Rooms[m.RoomID]; ok {
				// Send message to each client
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
