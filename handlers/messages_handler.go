package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-sms-gateway-api/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

type MessagesHandler struct {
	upgrader *websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mutex    *sync.Mutex
	db       *gorm.DB
}

func NewMessagesHandler(
	db *gorm.DB,
) *MessagesHandler {
	return &MessagesHandler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
		mutex:   &sync.Mutex{},
		db:      db,
	}
}

// HandleWebSocketConnection handles WebSocket connections
// @Summary Upgrade to WebSocket
// @Description Upgrade HTTP connection to WebSocket
// @Tags messages
// @Accept json
// @Produce json
// @Success 101 {string} string "WebSocket connection established"
// @Failure 500 {string} string "Internal Server Error"
// @Router /messages/ws [get]
// @Security BearerAuth
func (h *MessagesHandler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	log.Printf("WebSocket connection established %s\n", conn.RemoteAddr().String())

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("WebSocket connection closed: %v\n", err)
			h.mutex.Lock()
			delete(h.clients, conn)
			h.mutex.Unlock()
			break
		}
	}
}

// PublishMessage handles the message publishing endpoint
// @Summary Publish a message
// @Description Publish a message to the server
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} models.SMSMessage
// @Failure 500 {string} string "Internal Server Error"
// @Router /messages/publish [post]
// @Security BearerAuth
// @Param message body models.SMSMessageRequestDto true "Message to publish"
func (h *MessagesHandler) PublishMessage(w http.ResponseWriter, r *http.Request) {
	var request models.SMSMessageRequestDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid message format", http.StatusBadRequest)
		return
	}

	log.Printf("Received message: %s", request.Message)
	message := models.SMSMessage{
		Message:     request.Message,
		PhoneNumber: request.PhoneNumber,
	}

	clientErrors := make([]string, 0)

	h.mutex.Lock()
	for client := range h.clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			client.Close()
			delete(h.clients, client)
			clientErrors = append(clientErrors, fmt.Sprintf("Error sending message to %s: %v", client.RemoteAddr().String(), err))
		} else {
			log.Printf("Message sent to client: %s", client.RemoteAddr().String())
		}
	}
	h.mutex.Unlock()

	if len(clientErrors) > 0 {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send message to some clients: %v", clientErrors))
		return
	}

	respondWithJSON(w, http.StatusOK, message)
}
