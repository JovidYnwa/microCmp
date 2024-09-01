package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
}

type Account struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstNmae, lastName string) *Account {
	return &Account{
		FirstName: firstNmae,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
		CreatedAt: time.Now().UTC(),
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
