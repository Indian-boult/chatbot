package websocket

import (
	"chatbot/pkg/chat"
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"chatbot/pkg/db"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	conn    *websocket.Conn
	mu      sync.Mutex // Protects the history slice
	history []string   // Stores the last two messages sent and received
}

// WebsocketHandler handles incoming WebSocket connections and messages.
func WebsocketHandler(chatAI *chat.OpenAI, database *db.Database, w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	conn := &Connection{conn: wsConn}
	defer conn.conn.Close()

	for {
		_, msg, err := conn.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		intent, err := chatAI.ImageIntentChecker(context.Background(), string(msg))
		if err != nil {
			log.Println("Intent checking error:", err)
			break
		}

		var response string
		switch {
		case strings.HasPrefix(intent, "retrieve_image:"):
			identifier := strings.TrimPrefix(intent, "retrieve_image:")
			response = conn.retrieveImage(database, identifier)

		case strings.HasPrefix(intent, "discuss_image:"):
			identifier := strings.TrimPrefix(intent, "discuss_image:")
			response = conn.discussImage(chatAI, database, identifier)

		case intent == "upload_image":
			response = "To upload an image, please use the upload form below."

		case intent == "normal_conversation":
			response = conn.normalConversation(chatAI, string(msg))

		default:
			response = "I'm not sure how to respond to that."
		}

		if err := conn.conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func (c *Connection) addToHistory(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.history = append(c.history, message)
	// Ensure we only keep the last 4 messages (2 turns of conversation)
	if len(c.history) > 4 {
		c.history = c.history[len(c.history)-4:]
	}
}

func (c *Connection) retrieveImage(database *db.Database, identifier string) string {
	imageInfo, err := database.GetImage(identifier)
	if err != nil {
		return "Error retrieving image."
	}
	imageURL := "Here is your image URL: " + imageInfo.URL
	c.addToHistory(imageURL)
	return imageURL
}

func (c *Connection) discussImage(chatAI *chat.OpenAI, database *db.Database, identifier string) string {
	imageInfo, err := database.GetImage(identifier)
	if err != nil {
		return "Error fetching image for discussion."
	}
	response, err := chatAI.DiscussImage(context.Background(), imageInfo.URL)
	if err != nil {
		return "Error discussing the image."
	}
	c.addToHistory(response)
	return response
}

func (c *Connection) normalConversation(chatAI *chat.OpenAI, userInput string) string {
	ctx := context.Background()
	c.addToHistory(userInput) // Add user input to history
	prompt := chatAI.GeneratePromptWithHistory(c.history, userInput)
	response, err := chatAI.AskOpenAI(ctx, prompt)
	if err != nil {
		return "Error getting a response from OpenAI."
	}
	c.addToHistory(response) // Add AI response to history
	return response
}
