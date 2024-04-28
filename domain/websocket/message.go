package ws

import (
	"Omnichannel-CRM/package/presentation"
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const RoomJoinedAction = "room-joined"
const ListOnlineUserAction = "online-users"

type Message struct {
	Action     string               `json:"action"`
	Message    presentation.Message `json:"message"`
	Sender     *Client              `json:"sender"`
	Room       *Room                `json:"room"`
	Online     []*Room              `json:"online"`
	OnlineUser []*Client            `json:"online_user"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
