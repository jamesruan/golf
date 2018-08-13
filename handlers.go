package golf

import (
	"sync"
)

var topic_handlers = new(sync.Map) //string -> chan *Event

type eventHandler struct {
	in     chan *Event
	queue  []*Event
	logger EventLogger
	sync.Mutex
}

func (h *eventHandler) init() {
	h.Lock()
	// by default the stderrLogger is used
	h.logger = stderrLogger
	go func() {
		for e := range h.in {
			h.Lock()
			l := h.logger
			h.Unlock()
			// if no logger registered, the event is dropped
			if l != nil {
				l.Log(e)
			}
		}
	}()
	h.Unlock()
}

func (h *eventHandler) setLogger(logger EventLogger) {
	h.Lock()
	h.logger = logger
	h.Unlock()
}

func (h *eventHandler) setFilter(filter EventFilter) {
	h.Lock()
	h.logger.SetFilter(filter)
	h.Unlock()
}

// getTopicHandler return an eventHandler for topic, create and init one if not existed
func getTopicHandler(topic string) *eventHandler {
	h, ok := topic_handlers.Load(topic)
	if ok {
		return h.(*eventHandler)
	}

	hn := &eventHandler{
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
