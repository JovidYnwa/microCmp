package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/JovidYnwa/microCmp/db"
	"github.com/JovidYnwa/microCmp/types"
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
		fmt.Println(w.CmpIterationData())
	}
}

func (c *CmpWoker) RequestCmpID() any {
	companies, err := c.dbPg.GetActiveCompanies()
	if err != nil {
		fmt.Println("receiving active companies err: ", err)
		return nil
	}

	for _, company := range companies {
		_, err := c.dbDwh.GetCompanySubscribers(company.ID)
		if err != nil {
			fmt.Println("receiving active companies err: ", err)
			return nil
		}
		fmt.Printf("subscriber amount %d for companyId %s", company.ID, company.SmsText)
	}
	return len(companies)
}

func (c *CmpWoker) CmpIterationData() (int, error) {
	companies, err := c.dbPg.GetActiveCompanies()
	if err != nil {
		return 0, fmt.Errorf("getting active companies: %w", err)
	}

	for _, company := range companies {
		cmp := types.CmpStatistic{
			ID:        company.ID,
			StartDate: time.Now(),
			Efficiency:       0, // or some default value
			SubscriberAmount: 0, // or some default value
		}

		_, err := c.dbPg.InsertCmpStatistic(cmp)
		if err != nil {
			return 0, fmt.Errorf("inserting company %d statistics: %w", company.ID, err)
		}

		log.Printf("company statistics created for ID: %d", company.ID)
	}

	return len(companies), nil
}
