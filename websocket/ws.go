package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	clientMu  sync.Mutex
	broadcast = make(chan []byte)
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	clientMu.Lock()
	clients[conn] = true
	clientMu.Unlock()

	defer func() {
		clientMu.Lock()
		delete(clients, conn)
		clientMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			return
		}
	}
}

func BroadcastMessage(msg []byte) {
	clientMu.Lock()
	defer clientMu.Unlock()
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing WebSocket message:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func StartWebSocketServer(addr string) {
	http.HandleFunc("/ws", HandleWebSocket)

	go func() {
		for {
			msg := <-broadcast
			BroadcastMessage(msg)
		}
	}()

	log.Println("WebSocket server running on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("WebSocket server error:", err)
	}
}
