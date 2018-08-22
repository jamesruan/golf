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
)

type Event struct {
	Topic  string
	Level  Level
	Time   time.Time
	File   string
	Line   int
	Fmt    string
	Args   []interface{}
	Fields map[string]interface{}
}

func New(calldepth int, topic string, level Level, fmt string, args []interface{}, fields map[string]interface{}) *Event {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	return &Event{
		Topic:  topic,
		Level:  level,
		Time:   time.Now(),
		File:   file,
		Line:   line,
		Fmt:    fmt,
		Args:   args,
		Fields: fields,
	}
}

type Logger interface {
	Log(event *Event)
}
