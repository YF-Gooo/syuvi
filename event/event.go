package event

type Event struct {
	Target  string
	Source  string
	Content string
}

type EventReceiver interface {
	OnEvent(evt Event)
}
