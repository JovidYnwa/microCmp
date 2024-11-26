package types

import "time"

type ActiveCmp struct {
	ID      int
	SmsText any
}

type CmpSubscriber struct {
	Msisdn string
	LangID int
}

type CmpStatistic struct {
	ID               int
	Efficiency       int
	SubscriberAmount int
	StartDate        time.Time
}
