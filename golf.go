package golf

import (
	"sync"
)

type Formatter interface {
	Format(e *Event) []byte
}

var sinkWg sync.WaitGroup
var sinkCloseSignal chan struct{} = make(chan struct{})

func closeAllSink() {
	close(sinkCloseSignal)
	sinkWg.Wait()
}

// Level the logging Level.
type Level int

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case WARN:
		return "WARN "
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return ""
	}
}

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	NOLEVEL
)

// Handler handles Event.
type Handler interface {
	Handle(*Event)
}
