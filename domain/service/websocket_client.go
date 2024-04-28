package service

import (
	"encoding/json"
	"fmt"
	"net/url"
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

type WebSocketClient struct {
	id        string
	configStr string
	sendBuf   chan []byte

	mu     sync.RWMutex
	wsconn *websocket.Conn
}

func NewWebSocketClient(host, channel, userId, query string) (*WebSocketClient, error) {

	if connections == nil {
		connections = make(map[string]*WebSocketClient)
	}

	connection, ok := connections[userId]
	if ok {
		go connection.Connect()
		go connection.listenWrite()
		return connection, nil
	}

	temp := WebSocketClient{
		sendBuf: make(chan []byte, 256),
		id:      userId,
	}
	
	u := url.URL{Scheme: "ws", Host: host, Path: channel, RawQuery: query}
	temp.configStr = u.String()
	
	go temp.Connect()
	go temp.listenWrite()

	time.Sleep(200 * time.Millisecond)
	
	connections[userId] = &temp

	return &temp, nil
}

func (conn *WebSocketClient) Connect() *websocket.Conn {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.wsconn != nil {
		return conn.wsconn
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		ws, _, err := websocket.DefaultDialer.Dial(conn.configStr, nil)
		if err != nil {
			conn.log("connect", err, fmt.Sprintf("Cannot connect to websocket: %s", conn.configStr))
			continue
		}
		conn.log("connect", nil, fmt.Sprintf("connected to websocket to %s", conn.configStr))
		conn.wsconn = ws
		return conn.wsconn
	}
}

// Write data to the websocket server
func (conn *WebSocketClient) Write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()
	ws := conn.wsconn
	if ws == nil {
		err := fmt.Errorf("conn.ws is nil")
		return err
	}
	conn.sendBuf <- data
	return nil

}

func (conn *WebSocketClient) listenWrite() {
	for data := range conn.sendBuf {
		ws := conn.Connect()
		if ws == nil {
			err := fmt.Errorf("conn.ws is nil")
			conn.log("listenWrite", err, "No websocket connection")
			continue
		}

		if err := ws.WriteMessage(
			websocket.TextMessage,
			data,
		); err != nil {
			conn.log("listenWrite", nil, "WebSocket Write Error")
		}
		conn.log("listenWrite", nil, fmt.Sprintf("send: %s", data))
	}
}

// Close will send close message and shutdown websocket connection
func (conn *WebSocketClient) Stop() {
	conn.mu.Lock()
	if conn.wsconn != nil {
		delete(connections, conn.id)
		conn.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.wsconn.Close()
		conn.wsconn = nil
	}
	conn.mu.Unlock()
}

// Log print log statement
// In real word I would recommend to use zerolog or any other solution
func (conn *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}
