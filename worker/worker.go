package worker

import (
	"fmt"
	"time"

	"github.com/JovidYnwa/microCmp/db"
)

type CmpWoker struct {
	taskName string
	ticker   *time.Ticker
	dbPg     db.WorkerMethod
	dbDwh    db.DwhStore
}

// Updated constructor to accept db.WorkerMethod
func NewCmpWoker(t string, interval time.Duration, dbPg db.WorkerMethod, dbDwh db.DwhStore) *CmpWoker {
	return &CmpWoker{
		taskName: t,
		ticker:   time.NewTicker(interval),
		dbPg:     dbPg,
		dbDwh:    dbDwh,
	}
}

func (w *CmpWoker) Start() {
	fmt.Println("Starting Worker ... ")
	for range w.ticker.C {
		fmt.Println(w.RequestCmpID())
	}
}

func (c *CmpWoker) RequestCmpID() any {
	companies, err := c.dbPg.GetActiveCompanies()
	if err != nil {
		fmt.Println("receiving active companies err: ", err)
		return nil
	}

	for _, company := range companies {
		subs, err := c.dbDwh.GetCompanySubscribers(company.ID)
		if err != nil {
			fmt.Println("receiving active companies err: ", err)
			return nil
		}
		fmt.Printf("subscriber amount %d for companyId %d", len(subs), company.ID)
	}
	return len(companies)
}
