package handlers

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/chat", chat)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(nil, w, r)
	})
}

func chat(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat" {
		status := http.StatusNotFound
		http.Error(w, "not found", status)
		return
	}

	if r.Method != http.MethodGet {
		status := http.StatusNotFound
		http.Error(w, "Invalid Method", status)
		return
	}

	http.ServeFile(w, r, "templates/chat.html")
}

func register(w http.ResponseWriter, r *http.Request) {}

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func wsEndpoint(nil, w http.ResponseWriter, r *http.Request) {}
