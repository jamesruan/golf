package golf

type TopicHandler mapHandler

// NewTopicHandler creates a TopicHandler
//
// TopicHandler delivers to handler with a 'name' identical to the Topic of Event.
// If no suitable handler is configured, the Event is dropped.
func NewTopicHandler() *TopicHandler {
	return &TopicHandler{
		smap: new(handlerMap),
	}
}

func (h *TopicHandler) Handle(e *Event) {
	(*mapHandler)(h).smap.Range(func(k string, v Handler) bool {
		if k == e.Topic {
			v.Handle(e)
		}
		return true
	})
}

// AddHandler adds a handler 'v' to handle 'topic'
func (h *TopicHandler) AddHandler(topic string, v Handler) {
	(*mapHandler)(h).smap.Store(topic, v)
}

// DeleteHandler delete handler for 'topic'
func (h *TopicHandler) DeleteHandler(topic string) {
	(*mapHandler)(h).smap.Delete(topic)
}
