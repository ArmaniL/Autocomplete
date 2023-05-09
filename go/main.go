package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type WebSocketMessage struct {
	Prefix string `json:"prefix"`
}

type Suggestion struct {
	Suggestions []string `json:"suggestions"`
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Define a WebSocket upgrader.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var completer AutoComplete

func main() {

	pathToDictionary := goDotEnvVariable("VOCAB")

	completer = NewAutoComplete()
	completer.AddWordsFromFile(pathToDictionary)

	gin.SetMode(gin.ReleaseMode)
	// Set up the Gin server.
	r := gin.Default()

	// Define the WebSocket handler.
	r.GET("/ws", func(c *gin.Context) {

		// Upgrade the HTTP connection to a WebSocket connection.
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		// Handle incoming WebSocket messages.
		for {
			// Read a WebSocket message.
			var message WebSocketMessage
			err := conn.ReadJSON(&message)

			if err != nil {
				fmt.Println(err)
				break
			}

			words, err := completer.GuessNWords(message.Prefix, 10)

			if err != nil {
				fmt.Println(err)
			}

			suggestion := Suggestion{
				Suggestions: words,
			}

			err = conn.WriteJSON(suggestion)

			if err != nil {
				fmt.Println(err)
				break
			}

		}
	})

	// Start the Gin server.
	r.Run()
}
