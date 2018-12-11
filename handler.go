package golf

import (
	"sync/atomic"
)

type mapHandler struct {
	name string
	smap *handlerMap
}

func (h mapHandler) Name() string {
	return h.name
}

type TopicHandler mapHandler

// NewTopicHandler creates a TopicHandler
//
// TopicHandler delivers to handler with a 'name' identical to the Topic of Event.
// If no suitable handler is configured, the Event is dropped.
func NewTopicHandler(name string) *TopicHandler {
	return &TopicHandler{
		name: name,
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

func (h *TopicHandler) AddHandlerForTopic(topic string, v Handler) {
	(*mapHandler)(h).smap.Store(topic, v)
}

func (h *TopicHandler) DeleteHandlerForTopic(topic string) {
	(*mapHandler)(h).smap.Delete(topic)
}

type LevelHandler struct {
	mapHandler
	level atomic.Value
}

// NewLevelHandler creates a LevelHandler.
//
// LevelHandler delivers Event to 'qualified' handler when it's level greater than or equal to the 'lvl',
// to 'below' handler if not.
func NewLevelHandler(name string, lvl Level, qualified Handler, below Handler) *LevelHandler {
	m := new(handlerMap)
	if qualified != nil {
		m.Store("_qualified", qualified)
	}
	if below != nil {
		m.Store("_below", below)
	}

	var level atomic.Value
	level.Store(lvl)

	return &LevelHandler{
		mapHandler: mapHandler{
			name: name,
			smap: m,
		},
		level: level,
	}
}

func (h *LevelHandler) Handle(e *Event) {
	level := h.level.Load().(Level)
	if e.Level >= level {
		qualified := h.mapHandler.smap.Load("_qualified")
		if qualified != nil {
			qualified.Handle(e)
		}
	} else {
		below := h.mapHandler.smap.Load("_below")
		if below != nil {
			below.Handle(e)
		}
	}
}

// SetLevel atomically sets the level.
func (h *LevelHandler) SetLevel(lvl Level) {
	h.level.Store(lvl)
}

// SetQualifiedHandler atomically sets the handle for the event with Level equal or above.
func (h *LevelHandler) SetQualifiedHandler(qualified Handler) {
	h.mapHandler.smap.Store("_qualified", qualified)
}

// SetBelowHandler atomically sets the handle for the event with Level below.
func (h *LevelHandler) SetBelowHandler(below Handler) {
	h.mapHandler.smap.Store("_below", below)
}
