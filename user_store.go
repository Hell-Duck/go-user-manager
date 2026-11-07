package main

import (
	"errors"
	"sync"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidName  = errors.New("user name is empty")
	ErrInvalidAge   = errors.New("user age is negative")
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserStore struct {
	users  map[int]User
	nextID int
	mu     sync.RWMutex
}

func NewUserStore() *UserStore {
	userStore := UserStore{}
	userStore.users = make(map[int]User)
	userStore.nextID = 1
	return &userStore
}

func (s *UserStore) AddUser(name string, age int) (User, error) {
	if name == "" {
		return User{}, ErrInvalidName
	}
	if age < 0 {
		return User{}, ErrInvalidAge
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	user := User{ID: s.nextID, Name: name, Age: age}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

func (s *UserStore) GetUser(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if value, ok := s.users[id]; ok {
		return value, nil
	}
	return User{}, ErrUserNotFound
}

func (s *UserStore) DeleteUser(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}

func (s *UserStore) UpdateUser(id int, name string, age int) (User, error) {
	if name == "" {
		return User{}, ErrInvalidName
	}
	if age < 0 {
		return User{}, ErrInvalidAge
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return User{}, ErrUserNotFound
	}
	user := User{ID: id, Name: name, Age: age}
	s.users[id] = user
	return user, nil
}
