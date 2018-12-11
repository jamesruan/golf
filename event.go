package golf

import (
	"github.com/Workiva/go-datastructures/list"
	"runtime"
	"time"
)

// Event records something to be logged.
type Event struct {
	Topic  string
	Level  Level
	Time   time.Time
	Pc     []uintptr
	Fmt    string
	Args   []interface{}
	Fields list.PersistentList // list of EventField
}

type EventField struct {
	Name  string
	Value interface{}
}

// DefaultEvent records level, topic, callstack.
func DefaultEvent(calldepth int, topic string, level Level, fmt string, args []interface{}, fields list.PersistentList) *Event {
	pc := make([]uintptr, 32)
	n := runtime.Callers(calldepth, pc)
	if n > 0 {
		pc = pc[:n]
	} else {
		pc = nil
	}
	return &Event{
		Topic:  topic,
		Level:  level,
		Time:   time.Now(),
		Pc:     pc,
		Fmt:    fmt,
		Args:   args,
		Fields: fields,
	}
}

// SimpleEvent records without level, topic and callstack.
func SimpleEvent(fmt string, args []interface{}, fields list.PersistentList) *Event {
	return &Event{
		Level:  NOLEVEL,
		Time:   time.Now(),
		Fmt:    fmt,
		Args:   args,
		Fields: fields,
	}
}
