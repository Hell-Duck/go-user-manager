package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/user/"))
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("invalid user ID")
	}
	return id, nil
}

func server() {
	log.Fatal(http.ListenAndServe(":4040", nil))
}

func getUser(w http.ResponseWriter, r *http.Request, store *UserStore) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := store.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request, store *UserStore) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := User{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}
	createdUser, err := store.AddUser(user.Name, user.Age)
	if err != nil {
		http.Error(w, "Invalid user info", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(createdUser)
	if err != nil {
		http.Error(w, "Invalid encode user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request, store *UserStore) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = store.DeleteUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func updateUser(w http.ResponseWriter, r *http.Request, store *UserStore) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := getID(r)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := User{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user info", http.StatusBadRequest)
		return
	}

	updatedUser, err := store.UpdateUser(id, user.Name, user.Age)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Invalid user info", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, "Invalid encode user", http.StatusInternalServerError)
		return
	}

	w.Header().Set(`Content-Type`, "application/json")

	_, err = w.Write(data)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
