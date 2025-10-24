package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aliBazrkar/go-chatapp/auth"
	"github.com/aliBazrkar/go-chatapp/chat"
	"github.com/aliBazrkar/go-chatapp/db"
	"github.com/aliBazrkar/go-chatapp/handlers"
)

func main() {

	// —— database setup ——
	var bufferSize uint16 = 1000

	dbConn, err := db.Initializer("./db/database.db", bufferSize)
	if err != nil {
		log.Fatalf("Database Initializing Failed: %v", err)
	}
	defer func() {
		sqlDB, _ := dbConn.Gorm.DB()
		sqlDB.Close()
	}()

	go dbConn.MessageWriter()

	// —— session management ——
	sm := auth.NewSessionManager(dbConn)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			sm.CleanupExpiredSessions()
		}
	}()

	// —— chat hub ——
	initHub := InitializeHub("Main", "@Main", dbConn)
	hub := chat.NewHub(initHub.ID, initHub.Name, initHub.Address, 100)
	go hub.Run(dbConn)

	// —— http server ——
	handler := handlers.NewHandler(sm, hub, dbConn)
	mux := http.NewServeMux()
	handler.Setup(mux)
	log.Fatal(http.ListenAndServe(":3000", mux))
}

/*
InitializeHub() function is a temporary method, since current
state of project uses ONE SINGLE hub for communication
however i have designed entire code in a way it can easily
scale to multiple hub in a possible close future.

to not break the current scalable Database schema for future
changes, i am initializing a constant single hub
so the program can actually run without any problem.

this function returns a constant hub weather it exists already
in database or it should be added to the database.
*/
func InitializeHub(name string, address string, dbConn *db.Database) *db.Hub {
	var hub *db.Hub

	exists, err := dbConn.HubExists("address", address)
	if err != nil {
		log.Println("DB error:", err)
		return nil
	}

	if exists {
		hub, err = dbConn.GetHub("address", address)
		if err != nil {
			log.Println("Can't get existing hub:", err)
			return nil
		}
		log.Println("Hub already exists:", hub.Name)
	} else {
		hub, err = dbConn.CreateHub(name, address)
		if err != nil {
			log.Println("Error creating hub:", err)
			return nil
		}
		log.Println("Hub created:", hub.Name)
	}

	return hub
}
