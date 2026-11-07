package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestAddUser(t *testing.T) {
	store := NewUserStore()
	user, err := store.AddUser("Alex", 25)
	if err != nil {
		t.Fatalf("AddUser() unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("user.ID: expected %d, got %d", 1, user.ID)
	}
	if user.Name != "Alex" {
		t.Errorf("user.Name: expected %s, got %s", "Alex", user.Name)
	}
	if user.Age != 25 {
		t.Errorf("user.Age: expected %d, got %d", 25, user.Age)
	}
}

func TestAddUserInvalidData(t *testing.T) {
	store := NewUserStore()
	_, err := store.AddUser("", 25)
	if !errors.Is(err, ErrInvalidName) {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
	_, err = store.AddUser("Alex", -25)
	if !errors.Is(err, ErrInvalidAge) {
		t.Errorf("expected ErrInvalidAge, got %v", err)
	}
	user, err := store.AddUser("Alex", 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	store := NewUserStore()
	_, err := store.AddUser("Alex", 25)
	if err != nil {
		t.Fatalf("AddUser() id 1 unexpected error: %v", err)
	}

	_, err = store.AddUser("Bob", 2)
	if err != nil {
		t.Fatalf("AddUser() id 2 unexpected error: %v", err)
	}

	_, err = store.UpdateUser(1, "Xela", 52)
	if err != nil {
		t.Fatalf("UpdateUser() unexpected error: %v", err)
	}

	user, err := store.GetUser(1)
	if err != nil {
		t.Fatalf("GetUser() id %d unexpected error: %v", 1, err)
	}
	if user.ID != 1 {
		t.Errorf("user.ID: expected %d, got %d", 1, user.ID)
	}
	if user.Name != "Xela" {
		t.Errorf("user.Name: expected %s, got %s", "Xela", user.Name)
	}
	if user.Age != 52 {
		t.Errorf("user.Age: expected %d, got %d", 52, user.Age)
	}

	user2, err := store.GetUser(2)
	if err != nil {
		t.Fatalf("GetUser() id %d unexpected error: %v", 2, err)
	}
	if user2.ID != 2 {
		t.Errorf("user.ID: expected %d, got %d", 2, user2.ID)
	}
	if user2.Name != "Bob" {
		t.Errorf("user.Name: expected %s, got %s", "Bob", user2.Name)
	}
	if user2.Age != 2 {
		t.Errorf("user.Age: expected %d, got %d", 2, user2.Age)
	}
}

func TestGetUserHTTP(t *testing.T) {
	store := NewUserStore()
	_, err := store.AddUser("Alex", 30)
	if err != nil {
		t.Fatalf("user no add")
	}

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	rec := httptest.NewRecorder()

	getUser(rec, req, store)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestGetUserNotFoundHTTP(t *testing.T) {
	store := NewUserStore()

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	rec := httptest.NewRecorder()

	getUser(rec, req, store)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %v", rec.Code)
	}
}

func TestCreateUserEmptyBodyHTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	createUser(rec, req, NewUserStore())

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", rec.Code)
	}
}

func CheckUser(store *UserStore, t *testing.T, expected User) {
	t.Helper()

	savedUser, err := store.GetUser(expected.ID)
	if err != nil {
		t.Fatalf("GetUser(%d): unexpected error: %v", expected.ID, err)
	}
	if savedUser.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, savedUser.ID)
	}
	if savedUser.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, savedUser.Name)
	}
	if savedUser.Age != expected.Age {
		t.Errorf("expected age %d, got %d", expected.Age, savedUser.Age)
	}
}

func TestCreateUserHTTP(t *testing.T) {
	store := NewUserStore()

	body := strings.NewReader(`{"Name":"Alex", "Age":1}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	rec := httptest.NewRecorder()

	createUser(rec, req, store)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", rec.Code)
		return
	}
	var user User
	err := json.Unmarshal(rec.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("failed to decode POST response: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1, got %v", user.ID)
	}
	if user.Name != "Alex" {
		t.Errorf("expected name Alex, got %s", user.Name)
	}
	if user.Age != 1 {
		t.Errorf("expected age 1, got %d", user.Age)
	}

	CheckUser(store, t, User{ID: 1, Name: "Alex", Age: 1})
}

func TestUpdateUserHTTP(t *testing.T) {
	store := NewUserStore()
	_, err := store.AddUser("Alex", 30)
	if err != nil {
		t.Fatalf("user no add")
	}
	body := strings.NewReader(`{"Name":"Bob", "Age":25}`)
	req := httptest.NewRequest(http.MethodPut, "/user/1", body)
	rec := httptest.NewRecorder()
	updateUser(rec, req, store)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", rec.Code)
		return
	}

	var user User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("invalid decode rec.Body")
	}
	if user.Name != "Bob" {
		t.Errorf("expected Name Bob, got %v", user.Name)
		return
	}
	if user.Age != 25 {
		t.Errorf("expected Age 25, got %v", user.Age)
		return
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1, got %v", user.ID)
		return
	}

	CheckUser(store, t, User{ID: 1, Name: "Bob", Age: 25})
}

func TestDeleteUserHTTP(t *testing.T) {
	store := NewUserStore()
	_, err := store.AddUser("Alex", 25)
	if err != nil {
		t.Fatalf("AddUser() unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	rec := httptest.NewRecorder()
	deleteUser(rec, req, store)
	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %v", rec.Code)
		return
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected nil, got %v", rec.Body)
		return
	}

	_, err = store.GetUser(1)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected %v, got %v", ErrUserNotFound, err)
	}
}

func TestUpdateUserNotFoundHTTP(t *testing.T) {
	body := strings.NewReader(`{"Name":"Alex", "Age":15}`)
	req := httptest.NewRequest(http.MethodPut, "/user/1", body)
	rec := httptest.NewRecorder()
	updateUser(rec, req, NewUserStore())
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %v", rec.Code)
	}
}

func TestGetUserInvalidIDHTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/abc", nil)
	rec := httptest.NewRecorder()
	getUser(rec, req, NewUserStore())
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", rec.Code)
	}
}

func TestCreateUserInvalidJSONHTTP(t *testing.T) {
	store := NewUserStore()
	body := strings.NewReader(`{"Name":"Alex","Age":}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	rec := httptest.NewRecorder()
	createUser(rec, req, store)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", rec.Code)
		return
	}

	_, err := store.GetUser(1)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected %v, got %v", ErrUserNotFound, err)
	}
}

func TestCreateUserInvalidAgeHTTP(t *testing.T) {
	store := NewUserStore()
	body := strings.NewReader(`{"Name":"Alex","Age":-5}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	rec := httptest.NewRecorder()
	createUser(rec, req, store)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	_, err := store.GetUser(1)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected %v, got %v", ErrUserNotFound, err)
	}
}

func TestDeleteUserNotFoundHTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/user/999", nil)
	rec := httptest.NewRecorder()
	deleteUser(rec, req, NewUserStore())
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %v", rec.Code)
	}
}

func TestConcurrentAddUser(t *testing.T) {
	var wg sync.WaitGroup
	store := NewUserStore()
	ch := make(chan int, 1000)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			user, err := store.AddUser("Alex", 25)
			defer wg.Done()
			if err != nil {
				t.Errorf("expected nil, got %v", err)
				ch <- -1
				return
			}
			ch <- user.ID
		}()
	}
	wg.Wait()
	check := make(map[int]bool, 1000)
	for i := 0; i < 1000; i++ {
		index := <-ch
		if index == -1 {
			t.Errorf("invalid error in AddUser")
			return
		}
		if check[index] {
			t.Errorf("duplicate ID: %d", index)
			return
		}
		check[index] = true
	}
	count := len(store.users)
	if count != 1000 {
		t.Errorf("expected 1000 users, got %d", count)
	}
}
