package handlers

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity; refine as needed.
	},
}

type webSocketHandler struct {
	grpcClient proto.TaskServiceClient
	clients    map[*websocket.Conn]struct{}
	broadcast  chan []byte
}

func newWebSocketHandler(grpcClient proto.TaskServiceClient) *webSocketHandler {
	return &webSocketHandler{
		grpcClient: grpcClient,
		clients:    make(map[*websocket.Conn]struct{}),
		broadcast:  make(chan []byte),
	}
}

func (h *webSocketHandler) handleConnections(c *gin.Context) {
	defer close(h.broadcast)
	go h.HandleMessages()
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	h.clients[ws] = struct {
	}{}

	for {
		//_, message, err := ws.ReadMessage()
		//if err != nil {
		//	log.Printf("Error reading message: %v", err)
		//	delete(h.clients, ws)
		//	break
		//}

		h.broadcast <- []byte("message")
		time.Sleep(1 * time.Second)
	}
}

func (h *webSocketHandler) HandleMessages() {
	for message := range h.broadcast {
		fmt.Println("Received ", string(message))
		for client := range h.clients {
			err := client.WriteJSON(response{Message: "response"})
			if err != nil {
				log.Printf("Error sending message: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
