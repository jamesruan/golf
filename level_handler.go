package golf

import (
	"sync/atomic"
)

type LevelHandler struct {
	mapHandler
	level atomic.Value
}

// NewLevelHandler creates a LevelHandler.
//
// LevelHandler delivers Event to 'qualified' handler when it's level greater than or equal to the 'lvl',
// to 'below' handler if not.
func NewLevelHandler(lvl Level, qualified Handler, below Handler) *LevelHandler {
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
			name: "_level_",
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

func (h *LevelHandler) Name() string {
	level := h.level.Load().(Level)
	return h.mapHandler.name + level.String()
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
