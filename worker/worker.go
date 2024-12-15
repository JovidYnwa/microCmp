package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JovidYnwa/microCmp/db"
	"github.com/JovidYnwa/microCmp/types"
)

type CmpWoker struct {
	taskName string
	ticker   *time.Ticker
	dbPg     db.WorkerMethod
	dbDwh    db.DwhStore
	taskFunc func() error // Function that holds the task
}

func NewCmpWoker(t string, interval time.Duration, dbPg db.WorkerMethod, dbDwh db.DwhStore, taskFunc func() error) *CmpWoker {
	return &CmpWoker{
		taskName: t,
		ticker:   time.NewTicker(interval),
		dbPg:     dbPg,
		dbDwh:    dbDwh,
		taskFunc: taskFunc,
	}
}

func (w *CmpWoker) Start() {
	fmt.Println("Starting Worker ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping...")
			return
		case <-w.ticker.C:
			// Call the task function directly (taskFunc)
			err := w.taskFunc()
			if err != nil {
				log.Printf("Error executing task %s: %v", w.taskName, err)
			}
		case <-shutdownSignal:
			log.Println("Shutdown signal received.")
			cancel()
		}
	}
}

// Task 1: SetCmpIteration function
func SetCmpIteration(dbPg db.WorkerMethod) func() error {
	return func() error {
		companies, err := dbPg.GetActiveCompanies()
		if err != nil {
			return fmt.Errorf("getting active companies: %w", err)
		}

		for _, company := range companies {
			cmp := types.CmpStatistic{
				BillingID:        company.BillingID,
				StartDate:        time.Now(),
				Efficiency:       0,
				SubscriberAmount: 0,
			}

			_, err := dbPg.InsertCmpStatistic(cmp)
			if err != nil {
				return fmt.Errorf("inserting company %d statistics: %w", company.BillingID, err)
			}

			log.Printf("Company statistics created for ID: %d", company.BillingID)
		}
		return nil
	}
}

// Task 2: CmpNotifier function
func CmpNotifier(dbPg db.WorkerMethod, dbDwh db.DwhStore) func() error {
	return func() error {
		companies, err := dbPg.GetActiveCompanies()
		if err != nil {
			return fmt.Errorf("getting active companies: %w", err)
		}

		for _, company := range companies {
			subs, err := dbDwh.GetCmpSubscribersNotify(company.BillingID)
			if err != nil {
				return fmt.Errorf("getting subscribers GetCmpSubscribersNotify: %w", err)
			}
			fmt.Printf("cmpID=%d subsciber amount = %d\n", company.BillingID, len(subs))
		}
		return nil
	}
}

// Task 3: CmpInerationStatistic function
func CmpStatisticUpdater(dbPg db.WorkerMethod, dbDwh db.DwhStore) func() error {
	return func() error {
		companies, err := dbPg.GetActiveCompanyItarations()
		if err != nil {
			return fmt.Errorf("getting active company iterations: %w", err)
		}

		for _, company := range companies {
			statisticData, err := dbDwh.GetCompanyStatistic(company.BillingID, company.ItarationDay)
			if err != nil {
				return fmt.Errorf("getting statistics for company %d: %w", company.BillingID, err)
			}

			if statisticData == nil {
				fmt.Printf("No statistics found for company ID %d on %s\n",
					company.BillingID, company.ItarationDay.Format("2006-01-02"))
				continue
			}

			fmt.Printf("Updating company ID %d - Subscribers: %d, date: %s,Efficiency: %.2f%%\n",
				company.ID, statisticData.SubscriberAmount, statisticData.StartDate, statisticData.Efficiency)

			if err = dbPg.UpdateIterationStatistic(company.ID, statisticData); err != nil {
				// Log the error and continue to the next company
				log.Printf("Error updating statistics for company ID %d (BillingID: %d) for date %s: %v\n",
				company.ID, 
				company.BillingID, 
				statisticData.StartDate.Format("2006-01-02"),
				err)
				continue
			}

			fmt.Printf("Successfully updated company ID %d statistics\n", company.ID)
		}
		return nil
	}
}


// Task 3: CmpStatistic function only fomr Pg Db (company_repetion) not using db
