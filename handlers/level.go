package handlers

import (
	"github.com/jamesruan/golf/event"
	"sync/atomic"
)

type Level struct {
	m     *handlerMap
	level atomic.Value
}

// NewLevel creates a handler that handles events by its Level.
//
// Level delivers Event to 'qualified' handler when it's level greater than or equal to the 'lvl',
// to 'below' handler if not.
func NewLevel(lvl event.Level, qualified event.Handler, below event.Handler) *Level {
	m := new(handlerMap)
	if qualified != nil {
		m.Store("_qualified", qualified)
	}
	if below != nil {
		m.Store("_below", below)
	}

	var level atomic.Value
	level.Store(lvl)

	return &Level{
		m:     m,
		level: level,
	}
}

func (h *Level) Handle(e *event.Event) {
	level := h.level.Load().(event.Level)
	if e.Level >= level {
		qualified := h.m.Load("_qualified")
		if qualified != nil {
			qualified.Handle(e)
		}
	} else {
		below := h.m.Load("_below")
		if below != nil {
			below.Handle(e)
		}
	}
}

// SetLevel atomically sets the level.
func (h *Level) SetLevel(lvl event.Level) {
	h.level.Store(lvl)
}

// SetQualifiedHandler atomically sets the handle for the event with Level equal or above.
func (h *Level) SetQualifiedHandler(qualified event.Handler) {
	h.m.Store("_qualified", qualified)
}

// SetBelowHandler atomically sets the handle for the event with Level below.
func (h *Level) SetBelowHandler(below event.Handler) {
	h.m.Store("_below", below)
}
