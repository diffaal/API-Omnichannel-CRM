package ws

import (
	"Omnichannel-CRM/domain/service"
	"Omnichannel-CRM/package/enum"
	"Omnichannel-CRM/package/logger"
	"Omnichannel-CRM/package/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Websocket struct {
	websocket service.IInteractionService
}

func NewWebsocket(websocket service.IInteractionService) *Websocket {
	interactionWebsocket := Websocket{
		websocket: websocket,
	}
	return &interactionWebsocket
}

func (ih *Websocket) WesocketListener(wsServer *WsServer) gin.HandlerFunc {
	errorMessage := make(map[string]string)

	fn := func(c *gin.Context) {
		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
			errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Websocket Connect] error upgrade to websocket: %+v", err))
			response.ResponseInvalidRequest(c, nil, errorMessage)
			return
		}

		listener := NewWsListener(conn, wsServer)
		wsServer.registerListener <- listener

		go listener.readPump()
		go listener.writePump()

	}
	return gin.HandlerFunc(fn)
}

func (ih *Websocket) WesocketConnection(wsServer *WsServer) gin.HandlerFunc {
	errorMessage := make(map[string]string)

	fn := func(c *gin.Context) {
		userId := c.Query("user_id")
		roomId := c.Query("room_id")
		// isAgent := c.Query("is_agent")

		if userId == "" || roomId == "" {
			errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
			errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Websocket Connect] user_id or room_id params is missing"))
			response.ResponseInvalidRequest(c, nil, errorMessage)
			return
		}

		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			errorMessage["errorMessage"] = enum.INVALID_QUERY_MESSAGE
			errorMessage["errorStatus"] = enum.INVALID_QUERY_STATUS
			logger.Info(fmt.Sprintf("[FAILED][Websocket Connect] error upgrade to websocket: %+v", err))
			response.ResponseInvalidRequest(c, nil, errorMessage)
			return
		}
		fmt.Println("INI ISI SERVER", wsServer)
		client := wsServer.findUserByID(userId)
		if client != nil {
			client.disconnect()
		}
		client = NewWsClient(conn, wsServer, userId, roomId)

		go client.writePump()
		go client.readPump()

		wsServer.registerClient <- client

	}
	return gin.HandlerFunc(fn)
}
