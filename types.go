package main

import "math/rand"

type Account struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balace    int64  `json:"balance"`
}

func NewAccount(firstNmae, lastName string) *Account {
	return &Account{
		ID:        int64(rand.Intn(1000)),
		FirstName: firstNmae,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
	}
}

type SubscriberStatus struct {
	SubType int
	Desc    string
}

type TRPL struct {
	Name string
	ID   int
}

type Region struct {
	Name string
	ID   int
}
