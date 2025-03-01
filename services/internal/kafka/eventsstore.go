package kafka

import (
	"sync"
)

var (
	mu     sync.Mutex
	events []string
)

func AddEvent(e string) {
	mu.Lock()
	defer mu.Unlock()
	events = append(events, e)
}

func GetEvents() []string {
	mu.Lock()
	defer mu.Unlock()
	copied := make([]string, len(events))
	copy(copied, events)
	return copied
}
