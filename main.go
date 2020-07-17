package main

import (
	"log"
	"net/http"

	"github.com/gerfg/chat-app/handler"
	"github.com/gerfg/chat-app/model"
)

func main() {
	r := model.NewRoom()
	http.Handle("/", &handler.Template{Filename: "chat.html"})
	http.Handle("/room", r)

	go r.Run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
