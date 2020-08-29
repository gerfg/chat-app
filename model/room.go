package model

import (
	"log"
	"net/http"

	"github.com/gerfg/chat-app/trace"
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

	tracer trace.Tracer
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan []byte),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		tracer:  trace.Off(),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			// Joining
			r.Clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.Leave:
			// leaving
			delete(r.Clients, client)
			close(client.Send)
			r.tracer.Trace("Client left")
		case msg := <-r.Forward:
			r.tracer.Trace("Message received: ", string(msg))
			// forward message to all clients
			for client := range r.Clients {
				client.Send <- msg
				r.tracer.Trace(" -- sent to client")
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
