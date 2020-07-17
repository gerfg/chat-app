package model

import "github.com/gorilla/websocket"

// Client represent a single chatting user
type Client struct {
	// web socket for this client
	Socket *websocket.Conn
	// channel which messages are sent
	Send chan []byte
	// room on this client is chatting
	Room *Room
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		c.Room.Forward <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
