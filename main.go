package main

import (
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/aliBazrkar/go-chatapp/handlers"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	handlers.SetupRoutes()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
