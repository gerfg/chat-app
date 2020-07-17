package model

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type Room struct {
	// Forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	Forward chan []byte
	// join is a channel for clients whishing to join the room.
	Join chan *Client
	// Leave is a channel for clients whishing to leave the room.
	Leave chan *Client
	// clients holds all the current clients in this room.
	Clients map[*Client]bool
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan []byte),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			// Joining
			r.Clients[client] = true
		case client := <-r.Leave:
			// leaving
			delete(r.Clients, client)
			close(client.Send)
		case msg := <-r.Forward:
			// forward message to all clients
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &Client{
		Socket: socket,
		Send:   make(chan []byte, messageBufferSize),
		Room:   r,
	}
	r.Join <- client
	defer func() { r.Leave <- client }()
	go client.Write()
	client.Read()
}
