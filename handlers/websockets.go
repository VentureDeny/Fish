package handlers

import (
	"fmt"
	"log"
	"net/http"

	"fish/db"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust origin check as needed
	},
}

// WebSocketHandler handles the /ws WebSocket route
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Generic WebSocket client connected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket connection error: %v", err)
			}
			break
		}

		log.Printf("Message received: %s", message)

		// Store the message in the database
		err = storeMessage(string(message))
		if err != nil {
			log.Printf("Failed to store message: %v", err)
			errorMsg := fmt.Sprintf("Error storing message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
			continue
		}

		// Echo the message back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}
	}

	log.Println("Generic WebSocket client disconnected")
}

func storeMessage(content string) error {
	query := `INSERT INTO messages (content) VALUES ($1)`
	_, err := db.DB.Exec(query, content)
	return err
}
