package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TaskRequest struct {
	TaskID string `json:"task_id"`
}

type TaskResponse struct {
	Message string `json:"message"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Read message from WebSocket
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var taskRequest TaskRequest
		if err := json.Unmarshal(message, &taskRequest); err != nil {
			fmt.Println("Invalid JSON format:", err)
			return
		}

		fmt.Printf("Received Task ID: %s\n", taskRequest.TaskID)

		// Create a response message
		response := TaskResponse{Message: "Task " + taskRequest.TaskID + " processed successfully"}
		responseMessage, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error creating response:", err)
			return
		}

		// Write response message back to WebSocket
		if err := conn.WriteMessage(websocket.TextMessage, responseMessage); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", echo) // Handle WebSocket requests

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
