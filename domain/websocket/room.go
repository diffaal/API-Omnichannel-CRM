// room.go
package ws

import "log"

type Room struct {
	ID         string `json:"id"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

// NewRoom creates a new Room
func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)
			
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
			
		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
			log.Printf("INI CLIENT DI ROOM %s: %v", room.ID, room.clients)
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	delete(room.clients, client)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		log.Print("ini clients", client)
		client.send <- message
	}
}

func (room *Room) GetId() string {
	return room.ID
}