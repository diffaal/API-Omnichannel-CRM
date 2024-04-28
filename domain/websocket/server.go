package ws

import (
	"Omnichannel-CRM/domain/repository"
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/logger"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WsServer struct {
	clients            map[string]*Client
	listeners          map[*Listener]bool
	registerListener   chan *Listener
	unregisterListener chan *Listener
	registerClient     chan *Client
	unregisterClient   chan *Client
	notification       chan []byte
	rooms              map[string]*Room
}

var lock = &sync.Mutex{}

var singleInstance *WsServer

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		singleInstance = &WsServer{
			listeners:          make(map[*Listener]bool),
			registerListener:   make(chan *Listener),
			unregisterListener: make(chan *Listener),
			clients:            make(map[string]*Client),
			registerClient:     make(chan *Client),
			unregisterClient:   make(chan *Client),
			notification:       make(chan []byte),
			rooms:              make(map[string]*Room),
		}
	}
	return singleInstance
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case listnener := <-server.registerListener:
			server.registerListenerToServer(listnener)

		case listnener := <-server.unregisterListener:
			server.unregisterListenerToServer(listnener)

		case client := <-server.registerClient:
			server.registerClientToServer(client)

		case client := <-server.unregisterClient:
			server.unregisterClientToServer(client)

		case message := <-server.notification:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for _, client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) registerListenerToServer(listener *Listener) {
	server.listeners[listener] = true
	server.listOnlineRooms(JoinRoomAction)
}

func (server *WsServer) unregisterListenerToServer(listener *Listener) {
	delete(server.listeners, listener)
}

func (server *WsServer) registerClientToServer(client *Client) {
	server.clients[client.ID] = client
	server.listOnlineRooms(UserJoinedAction)
}

func (server *WsServer) unregisterClientToServer(client *Client) {
	delete(server.clients, client.ID)
	server.listOnlineRooms(UserLeftAction)
}

func (server *WsServer) listOnlineRooms(action string) {
	rooms := make([]*Room, len(server.rooms))
	i := 0

	for _, k := range server.rooms {
		rooms[i] = k
		i++
	}

	clients := make([]*Client, len(server.clients))
	i = 0

	for _, k := range server.clients {
		clients[i] = k
		i++
	}
	for listener := range server.listeners {
		message := &Message{
			Online:     rooms,
			OnlineUser: clients,
		}
		listener.send <- message.encode()
	}
}

func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room

	foundRoom, ok := server.rooms[ID]
	if !ok {
		return nil
	}

	return foundRoom
}

func (server *WsServer) findUserByID(ID string) *Client {
	var foundclient *Client

	foundclient, ok := server.clients[ID]
	if !ok {
		return nil
	}

	return foundclient
}

func (server *WsServer) createRoom(id string) *Room {
	room := NewRoom(id)
	go room.RunRoom()
	server.rooms[id] = room

	return room
}

func SetupWebsocketRouter(dbCRM *gorm.DB, dbOmnichannel *gorm.DB, wsServer *WsServer) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(gin.LoggerWithWriter(logger.Logger.Writer()))

	router.Static("/files", "/var/www")

	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	interactionRepo := repository.NewInteractionRepository(dbOmnichannel)
	messageRepo := repository.NewMessageRepository(dbOmnichannel)
	userRepo := repository.NewUserRepository(dbCRM)
	reporterRepo := repository.NewReporterRepository(dbOmnichannel)

	gmailService := service.NewGmailService()
	emailRepo := repository.NewEmailRepository(dbOmnichannel, gmailService)
	threadRepo := repository.NewThreadRepository(dbOmnichannel)
	emailService := service.NewEmailService(interactionRepo, messageRepo, reporterRepo, emailRepo, *threadRepo)

	interactionService := service.NewInteractionService(interactionRepo, messageRepo, userRepo, reporterRepo, emailService, threadRepo)
	websocket := NewWebsocket(interactionService)

	router.GET("/ws/listen", websocket.WesocketListener(wsServer))
	router.GET("/ws", websocket.WesocketConnection(wsServer))

	return router
}
