package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"fish/db"

	"github.com/gorilla/websocket"
)

// Humi represents the humidity data structure
type Humi struct {
	ID        int       `json:"id"`
	Humidity  float64   `json:"humidity"`
	Timestamp time.Time `json:"timestamp"`
}

// HumiHandler handles the /humi WebSocket route
func HumiHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Humi WebSocket client connected")

	for {
		_, _, err := conn.ReadMessage() // Assuming client sends something to trigger humidity send
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Humi WebSocket connection error: %v", err)
			}
			break
		}

		// Retrieve humidity data from the database or sensor
		humi, err := getHumidity()
		if err != nil {
			log.Printf("Failed to retrieve humidity: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Error retrieving humidity"))
			continue
		}

		// Send humidity data to the client
		humiJSON, _ := json.Marshal(humi)
		err = conn.WriteMessage(websocket.TextMessage, humiJSON)
		if err != nil {
			log.Printf("Failed to send humidity data: %v", err)
			break
		}
	}

	log.Println("Humi WebSocket client disconnected")
}

func getHumidity() (Humi, error) {
	// Implement your humidity retrieval logic here
	// For example, fetch from the database or read from a sensor
	var humi Humi
	query := `SELECT id, humidity, timestamp FROM humidity ORDER BY timestamp DESC LIMIT 1`
	err := db.DB.QueryRow(query).Scan(&humi.ID, &humi.Humidity, &humi.Timestamp)
	return humi, err
}
