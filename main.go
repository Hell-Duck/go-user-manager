package main

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidName  = errors.New("user name is empty")
	ErrInvalidAge   = errors.New("user age is negative")
)

type User struct {
	ID   int
	Name string
	Age  int
}

type UserStore struct {
	users  map[int]User
	nextID int
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
	user := User{ID: s.nextID, Name: name, Age: age}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

func (s *UserStore) GetUser(id int) (User, error) {
	if value, ok := s.users[id]; ok {
		return value, nil
	}
	return User{}, ErrUserNotFound
}

func (s *UserStore) DeleteUser(id int) error {
	if _, ok := s.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}

func (s *UserStore) UpdateUser(id int, name string, age int) (User, error) {
	if _, ok := s.users[id]; !ok {
		return User{}, ErrUserNotFound
	}
	if name == "" {
		return User{}, ErrInvalidName
	}
	if age < 0 {
		return User{}, ErrInvalidAge
	}
	user := User{ID: id, Name: name, Age: age}
	s.users[id] = user
	return user, nil
}

func main() {
	store := NewUserStore()
	_, err := store.AddUser("Alex", 25)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	_, err = store.AddUser("Bob", 30)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	user, err := store.GetUser(1)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else {
		fmt.Println(user)
	}
	user, err = store.GetUser(3)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(user)
	}
	err = store.DeleteUser(1)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	user, err = store.GetUser(1)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(user)
	}
	value, err := store.AddUser("Alex", 25)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else {
		fmt.Println(value)
	}
	value, err = store.UpdateUser(3, "xelA", 52)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else {
		fmt.Println(value)
	}
}
