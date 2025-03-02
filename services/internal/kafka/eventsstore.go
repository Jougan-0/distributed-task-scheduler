package kafka

import (
	"sync"
)

var (
	mu     sync.Mutex
	events []KafkaEvent
)

func AddEvent(e KafkaEvent) {
	mu.Lock()
	defer mu.Unlock()
	events = append(events, e)
}

func GetEvents() []KafkaEvent {
	mu.Lock()
	defer mu.Unlock()
	copied := make([]KafkaEvent, len(events))
	copy(copied, events)
	return copied
}
