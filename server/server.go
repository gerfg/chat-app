package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/gerfg/chat-app/handler"
	"github.com/gerfg/chat-app/model"
)

func Start() error {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	r := model.NewRoom()
	http.Handle("/chat", handler.MustAuth(&handler.Template{Filename: "chat.html"}))
	http.Handle("/login", &handler.Template{Filename: "login.html"})
	http.HandleFunc("/auth/", handler.LoginHandler)
	http.Handle("/room", r)

	go r.Run()

	log.Println("Starting web server on", *addr)
	return http.ListenAndServe(*addr, nil)
}
