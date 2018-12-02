package event

import (
	"runtime"
	"time"
)

type Level int

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DBG"
	case INFO:
		return "INF"
	case WARN:
		return "WRN"
	case LOG:
		return "LOG"
	case ERROR:
		return "ERR"
	default:
		return ""
	}

}

const (
	DEBUG Level = iota
	INFO
	LOG
	WARN
	ERROR
	NOLEVEL
)

type Event struct {
	Topic  string
	Level  Level
	Time   time.Time
	Pc     []uintptr
	Fmt    string
	Args   []interface{}
	Fields map[string]interface{}
}

func Simple(fmt string, args []interface{}, fields map[string]interface{}) *Event {
	return &Event{
		Level:  NOLEVEL,
		Time:   time.Now(),
		Fmt:    fmt,
		Args:   args,
		Fields: fields,
	}
}

func Default(calldepth int, topic string, level Level, fmt string, args []interface{}, fields map[string]interface{}) *Event {
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

type Logger interface {
	Queue() <-chan *Event
	Enqueue(*Event)
	Log(*Event)
	Flush()
}
