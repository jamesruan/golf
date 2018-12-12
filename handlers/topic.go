package handlers

import (
	"github.com/jamesruan/golf/event"
)

type Topic handlerMap

// NewTopic creates a handler to handle events by its Topic
//
// Topic delivers to handler with a 'name' identical to the Topic of Event.
// If no suitable handler is configured, the Event is dropped.
func NewTopic() *Topic {
	return new(Topic)
}

func (h *Topic) Handle(e *event.Event) {
	(*handlerMap)(h).Range(func(k string, v event.Handler) bool {
		if k == e.Topic {
			v.Handle(e)
		}
		return true
	})
}

// AddHandler adds a handler 'v' to handle 'topic'
func (h *Topic) AddHandler(topic string, v event.Handler) {
	(*handlerMap)(h).Store(topic, v)
}

// DeleteHandler delete handler for 'topic'
func (h *Topic) DeleteHandler(topic string) {
	(*handlerMap)(h).Delete(topic)
}
