package types

type ActiveCmp struct {
	ID      int
	SmsText any
}

type CmpSubscriber struct {
	Msisdn string
	LangID int
}
