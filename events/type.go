package events

const (
	Unknown Type = iota
	Message
)

type Type int

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Fetcher interface {
	Fetcher(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}
