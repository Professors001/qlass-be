package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	// Egress is used to avoid concurrent writes on the web socket connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessage() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("error setting read deadline:", err)
		return
	}

	// c.connection.SetReadLimit(1024) // We will be Back and Set it Later NOT RIGHT NOW

	c.connection.SetPongHandler(c.pongHandler)

	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message:", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("error marshalling message:", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("error handling message:", err)
		}

		if err := c.connection.WriteMessage(messageType, payload); err != nil {
			log.Println("error writing message:", err)
		}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		// cleanup connection
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case massage, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				return
			}

			data, err := json.Marshal(massage)
			if err != nil {
				log.Println("failed to marshal event", err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("failed to send message", err)
			}
			log.Println("message sent")

		case <-ticker.C:
			// log.Println("ping")

			//Send a pint to client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("write_msg_err : ", err)
				return
			}

		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	// log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
