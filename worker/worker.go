package worker

import (
	"fmt"
	"time"

	"github.com/JovidYnwa/microCmp/db"
)

type CmpWoker struct {
	taskName string
	ticker   *time.Ticker
	db       db.WorkerMethod
}

// Updated constructor to accept db.WorkerMethod
func NewCmpWoker(t string, interval time.Duration, db db.WorkerMethod) *CmpWoker {
	return &CmpWoker{
		taskName: t,
		ticker:   time.NewTicker(interval),
		db:       db,
	}
}

func (w *CmpWoker) Start() {
	fmt.Println("Starting Worker ... ")
	for range w.ticker.C {
		fmt.Println(w.RequestCmpID())
	}
}

func (c *CmpWoker) RequestCmpID() any {
	res, err := c.db.GetActiveCompanies()
	if err != nil {
		fmt.Println("Error while receiving: ", err)
		return nil
	}
	return res[0]
}