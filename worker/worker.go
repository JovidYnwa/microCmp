package worker

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type CmpWoker struct {
	taskName string
	ticker   *time.Ticker
}

func NewCmpWoker(t string, interval time.Duration) *CmpWoker {
	return &CmpWoker{
		taskName: t,
		ticker:   time.NewTicker(interval),
	}
}

func (w *CmpWoker) Start() {
	fmt.Println("Starting Worker ... ")
	for range w.ticker.C {
		fmt.Println(w.RequestCmpID())
	}
}

func (c *CmpWoker) RequestCmpID() string {
	cmpID := fmt.Sprintf("cmpId = %d", rand.IntN(100))
	return cmpID
}
