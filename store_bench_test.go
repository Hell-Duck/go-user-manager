package main

import (
	"sync"
	"testing"
)

type UserStoreMutex struct {
	users map[int]User
	mu    sync.Mutex
}

func NewUserStoreMutex() *UserStoreMutex {
	userStoreMutex := UserStoreMutex{}
	userStoreMutex.users = make(map[int]User)
	userStoreMutex.users[1] = User{1, "Alex", 25}
	return &userStoreMutex
}

func (s *UserStoreMutex) GetUser(id int) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, ok := s.users[id]; ok {
		return value, nil
	}
	return User{}, ErrUserNotFound
}

func BenchmarkGetUserParallel(b *testing.B) {
	store := NewUserStore()
	store.AddUser("Alex", 25)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := store.GetUser(1)
			if err != nil {
				b.Errorf("expected nil, got %v", err)
			}
		}
	})
}

func BenchmarkGetUserMutexParallel(b *testing.B) {
	store := NewUserStoreMutex()

	b.RunParallel(func(rb *testing.PB) {
		for rb.Next() {
			_, err := store.GetUser(1)
			if err != nil {
				b.Errorf("expected nil, got %v", err)
				return
			}
		}
	})
}
