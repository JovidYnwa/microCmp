package main

import "math/rand"

type Account struct {
	ID        int64
	FirstName string
	LastName  string
	Number    int64
	Balace    int64
}

func NewAccount(firstNmae, lastName string) *Account {
	return &Account{
		ID:        int64(rand.Intn(1000)),
		FirstName: firstNmae,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
	}
}
