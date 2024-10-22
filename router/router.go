package router

import (
	"log"

	"fish/handlers"

	"github.com/gorilla/mux"
)

// InitializeRoutes sets up all the WebSocket routes
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	// Register WebSocket routes
	router.HandleFunc("/ws", handlers.WebSocketHandler).Methods("GET")
	router.HandleFunc("/data", handlers.DataHandler).Methods("GET")
	router.HandleFunc("/command", handlers.CommandHandler).Methods("POST")
	router.HandleFunc("/temp", handlers.TempHandler).Methods("GET")
	router.HandleFunc("/humi", handlers.HumiHandler).Methods("GET")

	log.Println("Routes have been initialized.")
	return router
}
