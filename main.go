package main

import (
	"log"

	"github.com/gerfg/chat-app/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal("Listen and server:", err)
	}
}
