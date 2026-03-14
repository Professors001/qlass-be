package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"qlass-be/middleware"
	"qlass-be/usecases"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients      ClientList            // Global list (for metrics/debug)
	rooms        map[string]ClientList // Key: GamePIN, Value: List of Clients in that room
	sync.RWMutex                       // Use RWMutex for better concurrency
	handlers     map[string]EventHandler
	gameUseCase  usecases.GameUseCase
}

func NewManager(gameUseCase usecases.GameUseCase) *Manager {
	m := &Manager{
		clients:     make(ClientList),
		rooms:       make(map[string]ClientList), // <--- Important!
		handlers:    make(map[string]EventHandler),
		gameUseCase: gameUseCase,
	}
	m.setEventHandlers()
	return m
}

func (m *Manager) setEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
	m.handlers[EventJoinRoom] = JoinRoomHandler
	m.handlers[EventStartGame] = StartGameHandler
	m.handlers[EventNext] = NextHandler
	m.handlers[EventStudentAnswer] = StudentAnswerHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	//Check
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func SendMessage(event Event, c *Client) error {
	if c.GamePIN == "" {
		return errors.New("client is not in a room")
	}
	c.manager.Broadcast(c.GamePIN, event)
	return nil
}

func JoinRoomHandler(event Event, c *Client) error {
	var payload JoinRoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload for join_room: %v", err)
	}

	if payload.GamePIN == "" {
		return errors.New("game_pin is required")
	}

	// 1. (Optional but recommended) Validate PIN with GameUseCase/Redis here
	// if !gameUseCase.Exists(payload.GamePIN) { return error }

	// 2. Add client to the specific room in Manager
	c.manager.AddToRoom(payload.GamePIN, c)

	log.Printf("User %d joined room %s", c.UserID, payload.GamePIN)
	return nil
}

func StartGameHandler(event Event, c *Client) error {
	if c.GamePIN == "" {
		return errors.New("client is not in a room")
	}

	if c.Role != "teacher" {
		return errors.New("unauthorized: only teacher can start the game")
	}

	if err := c.manager.gameUseCase.StartGame(context.Background(), c.GamePIN, c.UserID); err != nil {
		return err
	}

	// Broadcast to room that game has started
	c.manager.Broadcast(c.GamePIN, Event{
		Type:    "GAME_STARTED",
		Payload: []byte("{}"),
	})

	return nil
}

func NextHandler(event Event, c *Client) error {
	if c.GamePIN == "" {
		return errors.New("client is not in a room")
	}

	if c.Role != "teacher" {
		return errors.New("unauthorized: only teacher can control the game")
	}

	wsEvent, err := c.manager.gameUseCase.NextStep(context.Background(), c.GamePIN, c.UserID)
	if err != nil {
		return err
	}

	if wsEvent.Type == "SYNC_TRIGGER" {
		c.manager.BroadcastGameState(c.GamePIN)
	} else {
		payloadBytes, err := json.Marshal(wsEvent.Payload)
		if err != nil {
			return err
		}

		c.manager.Broadcast(c.GamePIN, Event{
			Type:    wsEvent.Type,
			Payload: payloadBytes,
		})
	}

	// Auto-timeout logic now derives from the synced game state.
	syncState, err := c.manager.gameUseCase.GetSyncState(context.Background(), c.GamePIN, c.UserID)
	if err == nil && syncState != nil && syncState.Status == "running" && syncState.QuestionState == "answering" && syncState.QuestionStateObject != nil {
		go func(pin string, duration int, questionIndex int) {
			time.Sleep(time.Duration(duration) * time.Second)

			timeoutEvent, err := c.manager.gameUseCase.TimeoutQuestion(context.Background(), pin, questionIndex)
			if err != nil {
				log.Println("Timeout error:", err)
				return
			}

			if timeoutEvent != nil {
				if timeoutEvent.Type == "SYNC_TRIGGER" {
					c.manager.BroadcastGameState(pin)
				} else {
					payloadBytes, _ := json.Marshal(timeoutEvent.Payload)
					c.manager.Broadcast(pin, Event{
						Type:    timeoutEvent.Type,
						Payload: payloadBytes,
					})
				}
			}
		}(c.GamePIN, syncState.QuestionStateObject.TimeLimitSeconds, syncState.QuestionStateObject.CurrentQuestion)
	}

	if syncState != nil && syncState.Status == "finished" {
		go func(pin string) {
			time.Sleep(2 * time.Second) // Give time for the GAME_OVER message to reach clients
			c.manager.CloseRoom(pin)
		}(c.GamePIN)
	}

	return nil
}

func StudentAnswerHandler(event Event, c *Client) error {
	if c.GamePIN == "" {
		return errors.New("client is not in a room")
	}

	var payload StudentAnswerPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("invalid payload for student_answer: %v", err)
	}

	resp, wsEvent, err := c.manager.gameUseCase.SubmitAnswer(context.Background(), c.GamePIN, c.UserID, payload.OptionID)
	if err != nil {
		return err
	}

	respBytes, _ := json.Marshal(resp)
	c.egress <- Event{
		Type:    "ANSWER_SUBMITTED",
		Payload: respBytes,
	}

	if wsEvent != nil {
		if wsEvent.Type == "SYNC_TRIGGER" {
			c.manager.BroadcastGameState(c.GamePIN)
		} else {
			payloadBytes, _ := json.Marshal(wsEvent.Payload)
			c.manager.Broadcast(c.GamePIN, Event{
				Type:    wsEvent.Type,
				Payload: payloadBytes,
			})
		}
	}

	return nil
}

func (m *Manager) setEventHandler(eventType string, handler EventHandler) {
	m.handlers[eventType] = handler
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request, claims *middleware.JWTCustomClaims, pin string) {
	// 1. Join Game Logic (Validate & Update Redis)
	gameInfo, joinEvent, err := m.gameUseCase.JoinGame(r.Context(), pin, claims.UserId)
	if err != nil {
		log.Println("JoinGame error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create client with the PIN
	client := NewClient(conn, m, claims.UserId, claims.Role, pin)

	m.addClient(client)

	// 2. Send Initial Game State to the connecting client ONLY
	gameInfoBytes, _ := json.Marshal(gameInfo)
	client.egress <- Event{
		Type:    EventGameState,
		Payload: gameInfoBytes,
	}

	// 3. Broadcast the join event (if any) to room members.
	if joinEvent != nil {
		payloadBytes, _ := json.Marshal(joinEvent.Payload)
		event := Event{
			Type:    joinEvent.Type,
			Payload: payloadBytes,
		}

		// Broadcast to room
		m.Broadcast(pin, event)
	}
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Add to global list
	m.clients[client] = true

	// Add to specific Room (GamePIN)
	if client.GamePIN != "" {
		// Create room map if it doesn't exist yet
		if _, ok := m.rooms[client.GamePIN]; !ok {
			m.rooms[client.GamePIN] = make(ClientList)
		}
		m.rooms[client.GamePIN][client] = true
	}

	// Start threads
	go client.readMessage()
	go client.writeMessage()
}

// Update removeClient to clean up the room as well
func (m *Manager) removeClient(client *Client) {
	m.Lock()

	if _, ok := m.clients[client]; !ok {
		m.Unlock()
		client.shutdown()
		return
	}

	delete(m.clients, client)

	// Remove from Room
	if client.GamePIN != "" {
		if room, ok := m.rooms[client.GamePIN]; ok {
			delete(room, client)

			// Optional: Delete the room if it's empty to save memory
			if len(room) == 0 {
				delete(m.rooms, client.GamePIN)
			}
		}
	}
	m.Unlock()

	client.shutdown()

	if client.GamePIN != "" && client.UserID != 0 {
		go func() {
			ctx := context.Background()
			lobbyUpdate, err := m.gameUseCase.LeaveGame(ctx, client.GamePIN, client.UserID)
			if err != nil {
				log.Println("LeaveGame error:", err)
				return
			}

			if lobbyUpdate == nil {
				return
			}

			payloadBytes, _ := json.Marshal(lobbyUpdate.Payload)
			event := Event{
				Type:    lobbyUpdate.Type,
				Payload: payloadBytes,
			}
			m.Broadcast(client.GamePIN, event)
		}()
	}
}

// Add a client to a specific room
func (m *Manager) AddToRoom(pin string, client *Client) {
	m.Lock()
	defer m.Unlock()

	// 1. If client was in another room, remove them first (optional logic)
	if client.GamePIN != "" && client.GamePIN != pin {
		// Logic to remove from old room could go here
	}

	// 2. Assign PIN to client
	client.GamePIN = pin

	// 3. Create room if it doesn't exist
	if _, ok := m.rooms[pin]; !ok {
		m.rooms[pin] = make(ClientList)
	}

	// 4. Add to room
	m.rooms[pin][client] = true
}

func (m *Manager) CloseRoom(pin string) {
	m.Lock()
	var clientsToClose []*Client
	if clients, ok := m.rooms[pin]; ok {
		for client := range clients {
			clientsToClose = append(clientsToClose, client)
			delete(m.clients, client)
		}
		delete(m.rooms, pin)
	}
	m.Unlock()

	for _, client := range clientsToClose {
		client.shutdown()
	}
}

func (m *Manager) Broadcast(pin string, event Event) {
	m.RLock()
	defer m.RUnlock()

	if clients, ok := m.rooms[pin]; ok {
		for client := range clients {
			select {
			case client.egress <- event:
			default:
				log.Println("egress channel full, dropping message")
			}
		}
	}
}

func (m *Manager) BroadcastGameState(pin string) {
	m.RLock()
	room, ok := m.rooms[pin]
	if !ok {
		m.RUnlock()
		return
	}

	clients := make([]*Client, 0, len(room))
	for client := range room {
		clients = append(clients, client)
	}
	m.RUnlock()

	for _, client := range clients {
		syncState, err := m.gameUseCase.GetSyncState(context.Background(), pin, client.UserID)
		if err != nil {
			log.Println("GetSyncState error:", err)
			continue
		}

		payloadBytes, err := json.Marshal(syncState)
		if err != nil {
			log.Println("failed to marshal game state:", err)
			continue
		}

		select {
		case client.egress <- Event{Type: EventGameState, Payload: payloadBytes}:
		default:
			log.Println("egress channel full, dropping game state")
		}
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:3000":
		return true
	case "http://localhost:8080":
		return true
	case "":
		return true
	default:
		return false
	}
}
