package main

import "net/http"

func setupRoutes() {
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", registert)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(nil, w, r)
	})
}

func chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		status := http.StatusNotFound
		http.Error(w, "method not found", status)
		return
	}

	http.ServeFile(w, r, "templates/index.html")
}

func registert(w http.ResponseWriter, r *http.Request) {}

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func wsEndpoint(nil, w http.ResponseWriter, r *http.Request) {}
