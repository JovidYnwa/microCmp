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
				ID:               company.ID,
				StartDate:        time.Now(),
				Efficiency:       0,
				SubscriberAmount: 0,
			}

			_, err := dbPg.InsertCmpStatistic(cmp)
			if err != nil {
				return fmt.Errorf("inserting company %d statistics: %w", company.ID, err)
			}

			log.Printf("Company statistics created for ID: %d", company.ID)
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
			subs, err := dbDwh.GetCmpSubscribersNotify(company.ID)
			if err != nil {
				return fmt.Errorf("getting subscribers GetCmpSubscribersNotify: %w", err)
			}
			fmt.Printf("cmpID=%d subsciber amount = %d\n", company.ID, len(subs))
		}
		return nil
	}
}
