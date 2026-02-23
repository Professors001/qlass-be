package websocket

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventJoinRoom    = "join_room"
	EventStartGame   = "start_game"
)

type SendMessagePayload struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type JoinRoomPayload struct {
	GamePIN string `json:"game_pin"`
}
