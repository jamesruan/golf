package golf

import (
	"time"
)

type Level int

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

type Filter interface {
	FilterEvent(event *Event) (pass bool)
}

type EventLogger interface {
	Log(event *Event)
}
