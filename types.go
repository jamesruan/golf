package golf

import (
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

type EventLogger interface {
	Log(event *Event)
	SetFilter(filter EventFilter)
}
