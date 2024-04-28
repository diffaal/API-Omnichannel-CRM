package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents the websocket listener at the server
type Listener struct {
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
}

func NewWsListener(conn *websocket.Conn, wsServer *WsServer) *Listener {
	listener := &Listener{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}
	return listener
}

func (listener *Listener) readPump() {
	defer func() {
		listener.disconnect()
	}()

	listener.conn.SetReadLimit(maxMessageSize)
	listener.conn.SetReadDeadline(time.Now().Add(pongWait))
	listener.conn.SetPongHandler(func(string) error { listener.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := listener.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		listener.handleNewMessage(jsonMessage)
	}

}

func (listener *Listener) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		listener.conn.Close()
	}()
	for {
		select {
		case message, ok := <-listener.send:
			listener.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				listener.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := listener.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(listener.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-listener.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			listener.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := listener.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (listener *Listener) disconnect() {
	listener.wsServer.unregisterListener <- listener
	close(listener.send)
	listener.conn.Close()
}

func (listener *Listener) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}
}

// func (listener *Listener) handleJoinRoomMessage(message Message) {

// 	target := listener.wsServer.findClientByID(message.Message)

// 	if target == nil {
// 		return
// 	}

// 	listener.joinRoom(target)
// 	target.joinRoom(listener)

// }

// func (listener *Client) joinRoom(sender *Client) {

// 	listener.room.register <- listener

// 	listener.notifyRoomJoined(listener.room, sender)
// }

// func (listener *Client) notifyRoomJoined(room *Room, sender *Client) {
// 	message := Message{
// 		Action: RoomJoinedAction,
// 		Target: room,
// 		Sender: sender,
// 	}

// 	listener.send <- message.encode()
// }
