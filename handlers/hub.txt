package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Hub struct {
	// Registered clients.
	clients   map[string]*Client
	broadcast chan Message
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	mu         sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if existingClient, ok := h.clients[client.ID]; ok {
				fmt.Println(existingClient)
				delete(h.clients, client.ID)
				close(existingClient.Send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			targetClient, ok := h.clients[message.ReceiverID]
			if ok {
				msg, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling message for %s: %v", message.ReceiverID, err)
					h.mu.RUnlock() // Release lock before continuing
					continue
				}
				select {
				case targetClient.Send <- msg:
				default:
					h.unregister <- targetClient
				}
			}
			h.mu.RUnlock()
		}
	}
}
