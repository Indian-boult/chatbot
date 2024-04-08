package main

import (
	"log"
	"net/http"
	"os"

	"chatbot/pkg/api"
	"chatbot/pkg/chat"
	"chatbot/pkg/config"
	"chatbot/pkg/db"
	"chatbot/pkg/websocket"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize database connection
	databaseConnection := db.InitDB(cfg)
	database := db.NewDatabase(databaseConnection)

	// Initialize OpenAI client
	chatAI := chat.NewOpenAI(cfg)

	// Setup HTTP router
	router := mux.NewRouter()

	// Create an instance of Handler with dependencies
	handler := api.NewHandler(database, chatAI)

	// Register HTTP handlers
	router.HandleFunc("/upload", handler.UploadImage).Methods("POST")
	router.HandleFunc("/retrieve", handler.RetrieveImage).Methods("GET")

	// The chat functionality will be demonstrated through WebSocket for real-time interaction
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.WebsocketHandler(chatAI, database, w, r)
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
