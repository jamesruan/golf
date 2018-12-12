package handlers

import (
	"github.com/jamesruan/golf/event"
)

type Broadcast handlerMap

// NewBroadcast creates a handler to broadcast events.
//
// Broadcast handles event to all handlers registered.
// The order of calling registered handlers is random.
func NewBroadcast() *Broadcast {
	return new(Broadcast)
}

func (h *Broadcast) Handle(e *event.Event) {
	(*handlerMap)(h).Range(func(name string, handler event.Handler) bool {
		handler.Handle(e)
		return true
	})
}

//AddHandler adds handler 'v' identified as 'id'
func (h *Broadcast) AddHandler(id string, v event.Handler) {
	(*handlerMap)(h).Store(id, v)
}

//DeleteHandler deletes handler identified as 'id'
func (h *Broadcast) DeleteHandler(id string) {
	(*handlerMap)(h).Delete(id)
}
