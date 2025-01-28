package module

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/websocket/v2"
)

func SocketImpl(conn *websocket.Conn) {
	// Retrieve room ID and other context details
	roomID := conn.Params("id") // Unique room identifier
	// log.Printf("New WebSocket connection for room: %s", roomID)

	// Map to track connections by room
	var (
		mt  int    // Message type
		msg []byte // Message content
		err error  // Error handling
	)

	// Join room logic (could be further extended to a map of active rooms)
	joinRoom(roomID, conn)

	defer func() {
		// Ensure cleanup of resources upon disconnect
		leaveRoom(roomID, conn)
		// log.Printf("Connection closed for room: %s", roomID)
	}()

	for {
		// Read message from client
		if mt, msg, err = conn.ReadMessage(); err != nil {
			// log.Printf("Error reading message in room %s: %v", roomID, err)
			break
		}

		// log.Printf("Received message in room %s: %s", roomID, msg)

		// Broadcast message to other clients in the same room
		if err = broadcastToRoom(roomID, mt, msg); err != nil {
			// log.Printf("Error broadcasting message in room %s: %v", roomID, err)
			break
		}
	}
}

// Room management
var rooms = make(map[string]map[*websocket.Conn]bool)

// joinRoom adds a connection to a specified room.
func joinRoom(roomID string, conn *websocket.Conn) {
	if _, exists := rooms[roomID]; !exists {
		rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID][conn] = true
	// log.Printf("Connection added to room: %s", roomID)
}

// leaveRoom removes a connection from a specified room.
func leaveRoom(roomID string, conn *websocket.Conn) {
	if clients, exists := rooms[roomID]; exists {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(rooms, roomID) // Clean up empty rooms
			// log.Printf("Room %s is empty and removed", roomID)
		}
	}
}

// broadcastToRoom sends a message to all connections in a specified room.
func broadcastToRoom(roomID string, messageType int, message []byte) error {
	if clients, exists := rooms[roomID]; exists {
		for client := range clients {
			if err := client.WriteMessage(messageType, message); err != nil {
				log.Printf("Failed to send message to client in room %s: %v", roomID, err)
				client.Close()
				delete(clients, client) // Clean up broken connections
			}
		}
	}
	return nil
}

// Debug purpose and while on test only
func CheckEmptyClient() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		var count = 0
		for _, room := range rooms {
			fmt.Println(room)
			count++
		}
		if count > 0 {
			fmt.Println("Online rooms: " + strconv.Itoa(count))
		}
	}
}
