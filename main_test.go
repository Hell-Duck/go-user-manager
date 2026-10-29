package main

import (
	"errors"
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
