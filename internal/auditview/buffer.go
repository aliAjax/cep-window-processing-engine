package auditview

type Event struct {
	ID     string
	Fields map[string]string
}

type Buffer struct {
	events []Event
}

func (b *Buffer) Append(event Event) {
	b.events = append(b.events, cloneEvent(event))
}

func (b *Buffer) Snapshot() []Event {
	return b.events
}

func cloneEvent(event Event) Event {
	cloned := Event{ID: event.ID, Fields: make(map[string]string, len(event.Fields))}
	for key, value := range event.Fields {
		cloned.Fields[key] = value
	}
	return cloned
}
