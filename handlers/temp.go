package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"fish/db"

	"github.com/gorilla/websocket"
)

// Temp represents the temperature data structure
type Temp struct {
	ID          int       `json:"id"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

// TempHandler handles the /temp WebSocket route
func TempHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Temp WebSocket client connected")

	for {
		_, _, err := conn.ReadMessage() // Assuming client sends something to trigger temperature send
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Temp WebSocket connection error: %v", err)
			}
			break
		}

		// Retrieve temperature data from the database or sensor
		temp, err := getTemperature()
		if err != nil {
			log.Printf("Failed to retrieve temperature: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error retrieving temperature"))
			continue
		}

		// Send temperature data to the client
		tempJSON, _ := json.Marshal(temp)
		err = conn.WriteMessage(websocket.TextMessage, tempJSON)
		if err != nil {
			log.Printf("Failed to send temperature data: %v", err)
			break
		}
	}

	log.Println("Temp WebSocket client disconnected")
}

func getTemperature() (Temp, error) {
	// Implement your temperature retrieval logic here
	// For example, fetch from the database or read from a sensor
	var temp Temp
	query := `SELECT id, temperature, timestamp FROM temperature ORDER BY timestamp DESC LIMIT 1`
	err := db.DB.QueryRow(query).Scan(&temp.ID, &temp.Temperature, &temp.Timestamp)
	return temp, err
}
