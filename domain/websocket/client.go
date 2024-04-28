package ws

import (
	"Omnichannel-CRM/package/enum"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Client represents the websocket client at the server
type Client struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	ID       string `json:"id"`
	room     *Room
	platform string
	mu       sync.Mutex
}

func NewWsClient(conn *websocket.Conn, wsServer *WsServer, user_id string, room_id string) *Client {
	room := wsServer.findRoomByID(room_id)
	platform := enum.OMNICHANNEL
	if room == nil {
		room = wsServer.createRoom(room_id)
		platform = enum.WEBHOOK
	}
	client := &Client{
		ID:       user_id,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		room:     room,
		platform: platform,
	}
	room.register <- client
	return client
}

func (room *Client) GetId() string {
	return room.ID
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			client.mu.Lock()
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.mu.Unlock()
				return
			}
			client.mu.Unlock()
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregisterClient <- client
	client.room.unregister <- client

	_, ok := <-client.send
	if ok {
		close(client.send)
	}
	client.conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	log.Printf(" %s", string(jsonMessage))
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	message.Sender = client
	message.Room = client.room

	switch message.Action {
	case SendMessageAction:
		log.Print("INI BROADCAST", message)
		client.room.broadcast <- &message

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	}
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByID(message.Room.ID)
	if room == nil {
		return
	}

	room.unregister <- client
}

// func (client *Client) handleJoinRoomMessage(message Message) {

// 	target := client.wsServer.findClientByID(message.Message)

// 	if target == nil {
// 		return
// 	}

// 	client.joinRoom(target)
// 	target.joinRoom(client)

// }

// func (client *Client) joinRoom(sender *Client) {

// 	client.room.register <- client

// 	client.notifyRoomJoined(client.room, sender)
// }

// func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
// 	message := Message{
// 		Action: RoomJoinedAction,
// 		Target: room,
// 		Sender: sender,
// 	}

// 	client.send <- message.encode()
// }
