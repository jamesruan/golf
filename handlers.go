package golf

import (
	"sync"
)

var topic_handlers *sync.Map //string -> chan *Event

type eventHandler struct {
	in chan *Event
	queue []*Event
	loggers []EventLogger
	*sync.Mutex
}

func (h *eventHandler) init() {
	h.Lock()
	h.loggers = make([]EventLogger,0, 1)
	go func() {
		for e := range h.in {
			h.Lock()
			for _, l := range h.loggers {
				l.Log(e)
			}
			// if no logger registered, the event is dropped
			h.Unlock()
		}
	}()
	h.Unlock()
}

func init() {
	topic_handlers = new(sync.Map)
}

func getTopic(topic string) *eventHandler {
	h, ok := topic_handlers.Load(topic)
	if ok {
		return h.(*eventHandler)
	}

	hn := &eventHandler {
		in: make(chan *Event, 1),
	}

	ha, loaded := topic_handlers.LoadOrStore(topic, hn)
	if !loaded {
		hn.init()
		return hn
	} else {
		close(hn.in)
		return ha.(*eventHandler)
	}
}
