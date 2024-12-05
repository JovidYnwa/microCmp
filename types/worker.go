package types

import "time"

type ActiveCmp struct {
	ID      int
	SmsText any
}

type ActiveCmpIteration struct {
	ID           int
	ItarationDay time.Time
}

type CmpSubscriber struct {
	Msisdn string
	LangID int
}

type CmpStatistic struct {
	ID               int
	Efficiency       float64
	SubscriberAmount int
	StartDate        time.Time
}
