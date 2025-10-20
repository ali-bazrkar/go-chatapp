package main

import (
	"log"
	"net/http"

	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/aliBazrkar/go-chatapp/db"
	"github.com/aliBazrkar/go-chatapp/handlers"
)

func main() {

	dbConn, err := db.Initializer("./db")
	if err != nil {
		log.Fatalf("Database Initializing Failed: %v\n", err)
	}
	defer func() {
		sqlDB, _ := dbConn.Gorm.DB()
		sqlDB.Close()
	}()

	hub := chat.NewHub()
	go hub.Run(dbConn.Gorm)

	mux := http.NewServeMux()
	handlers.Setup(dbConn, hub, mux)
	log.Fatal(http.ListenAndServe(":3000", mux))
}
