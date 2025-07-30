package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client represents a single chat user connected via WebSocket
type Client struct {
	Hub      *Hub
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte // allows me to send messages to the client concurrently without blocking the main loop
}

type Message struct {
	SenderID       string    `json:"senderId"`
	SenderUsername string    `json:"senderUsername"`
	ReceiverID     string    `json:"receiverId,omitempty"` // omitempty means it won't be marshaled if empty
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"` //when the message was sent/received by the server
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageBytes, err := c.Conn.ReadMessage() //awl value hia messageType li katkun int
		if err != nil {
			// Check for expected close errors (client closed connection normally)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client read error for %s (%s): %v", c.Username, c.ID, err)
			}
			break
		}

		messageBytes = bytes.TrimSpace(bytes.ReplaceAll(messageBytes, newline, space))
		var incomingMsg Message
		if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
			continue // Skip the message but keep connection open
		}

		incomingMsg.SenderID = c.ID
        incomingMsg.SenderUsername = c.Username
        incomingMsg.Timestamp = time.Now()

        c.Hub.broadcast <- incomingMsg
	}
}
