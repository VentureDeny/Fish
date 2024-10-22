package handlers

import (
	"log"
	"net/http"

	"fish/db"

	"github.com/gorilla/websocket"
)

// CommandHandler handles the /command WebSocket route
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Command WebSocket client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Command WebSocket connection error: %v", err)
			}
			break
		}

		log.Printf("Command received: %s", message)

		// Process the command (e.g., execute system commands or control devices)
		err = executeCommand(string(message))
		if err != nil {
			log.Printf("Failed to execute command: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error executing command"))
			continue
		}

		// Acknowledge command execution
		conn.WriteMessage(websocket.TextMessage, []byte("Command executed successfully"))
	}

	log.Println("Command WebSocket client disconnected")
}

func executeCommand(cmd string) error {
	// Implement your command execution logic here
	// For example, store the command in the database or trigger an action
	query := `INSERT INTO commands (command) VALUES ($1)`
	_, err := db.DB.Exec(query, cmd)
	return err
}
