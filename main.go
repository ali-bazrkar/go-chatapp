package main

import (
	"log"
	"net/http"
)

func main() {
	// hub := NewHub()
	// go hub.run()

	setupRoutes()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
