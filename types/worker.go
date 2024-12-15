package types

import "time"

type ActiveCmp struct {
	BillingID int
	SmsText   any
}

type ActiveCmpIteration struct {
	ID           int
	BillingID    int
	ItarationDay time.Time
}

type CmpSubscriber struct {
	Msisdn string
	LangID int
}

type CmpStatistic struct {
	ID               int
	BillingID        int
	Efficiency       float64
	SubscriberAmount int
	StartDate        time.Time
}
