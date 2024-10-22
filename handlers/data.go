package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"fish/db"

	"github.com/gorilla/websocket"
)

// Data represents the structure of your data
type Data struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// DataHandler handles the /data WebSocket route
func DataHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Data WebSocket client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Data WebSocket connection error: %v", err)
			}
			break
		}

		log.Printf("Data received: %s", message)

		// Process the data (e.g., store or manipulate)
		data := Data{Content: string(message)}
		err = storeData(data)
		if err != nil {
			log.Printf("Failed to store data: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error storing data"))
			continue
		}

		// Optionally send a response back to the client
		response := map[string]string{"status": "Data received"}
		respJSON, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, respJSON)
	}

	log.Println("Data WebSocket client disconnected")
}

func storeData(data Data) error {
	query := `INSERT INTO data_table (content) VALUES ($1)`
	_, err := db.DB.Exec(query, data.Content)
	return err
}
