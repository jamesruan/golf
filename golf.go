package golf

import (
	"runtime"
	"time"
)

func newEvent(calldepth int, topic string, level Level, fmt string, args []interface{}, fields map[string]interface{}) *Event {
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

func Debugf(fmt string, args ...interface{}) {
	eh := getTopicHandler("")
	eh.in <- newEvent(2, "", DEBUG, fmt, args, nil)
}

func Infof(fmt string, args ...interface{}) {
	eh := getTopicHandler("")
	eh.in <- newEvent(2, "", INFO, fmt, args, nil)
}

func Logf(fmt string, args ...interface{}) {
	eh := getTopicHandler("")
	eh.in <- newEvent(2, "", LOG, fmt, args, nil)
}

func Warnf(fmt string, args ...interface{}) {
	eh := getTopicHandler("")
	eh.in <- newEvent(2, "", WARN, fmt, args, nil)
}

func Errorf(fmt string, args ...interface{}) {
	eh := getTopicHandler("")
	eh.in <- newEvent(2, "", ERROR, fmt, args, nil)
}
