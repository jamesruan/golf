package handlers

import (
	"github.com/jamesruan/golf"
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
func NewLevel(lvl golf.Level, qualified golf.Handler, below golf.Handler) *Level {
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

func (h *Level) Handle(e *golf.Event) {
	level := h.level.Load().(golf.Level)
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
func (h *Level) SetLevel(lvl golf.Level) {
	h.level.Store(lvl)
}

// SetQualifiedHandler atomically sets the handle for the event with Level equal or above.
func (h *Level) SetQualifiedHandler(qualified golf.Handler) {
	h.m.Store("_qualified", qualified)
}

// SetBelowHandler atomically sets the handle for the event with Level below.
func (h *Level) SetBelowHandler(below golf.Handler) {
	h.m.Store("_below", below)
}
