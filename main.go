package main

import (
	"net/http"
)

func main() {
	store := NewUserStore()
	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUser(w, r, store)
		case http.MethodPut:
			updateUser(w, r, store)
		case http.MethodDelete:
			deleteUser(w, r, store)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		createUser(w, r, store)
	})
	server()
}
