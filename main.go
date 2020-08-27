package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gerfg/chat-app/handler"
	"github.com/gerfg/chat-app/model"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	r := model.NewRoom()
	http.Handle("/", &handler.Template{Filename: "chat.html"})
	http.Handle("/room", r)

	go r.Run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen and server:", err)
	}
}
