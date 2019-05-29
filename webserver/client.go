package main

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

// it is called when upgrading the HTTP connection to a websocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true // disabling CORS
	},
}

// Client struct defines types for handling client side operations
type Client struct {
	id   [16]byte
	hub  *Hub
	conn *websocket.Conn
	send chan []byte // outbound channel of a client
}

// reader reads messages from the websocket connection to hub.
// hub gets message from clients through reader method of each client
func (c *Client) reader() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Printf("error: %v", err)
			break
		}

		event, err := UnMarshalMessage(message)
		if err != nil {
			logger.Printf("unable to parse: %v", err)
		} else {
			logger.Println("message ", event)
		}

		message = bytes.TrimSpace(message)
		c.hub.inbound <- message
	}
}

// writer sends messages from the hub to all websocket listeners
// allowing hub to send gathered to all clients
func (c *Client) writer() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// adding queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func generateClientID() [16]byte {
	return uuid.NewV4()
}

// initializes websocket connection for a client by upgrading the http connection
// for each connected client writer and reader is created. these goroutines are responsible of handling receiving/sending messages.
func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}

	client := &Client{id: generateClientID(), hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writer()
	go client.reader()
}
