package types

type BaseFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SubscriberGroupTypes struct {
	Notification  int
	ActionComited int
}

var SubscriberGroup SubscriberGroupTypes = SubscriberGroupTypes{
	Notification:  1, //Users for notification
	ActionComited: 2, //Users for statistics
}
