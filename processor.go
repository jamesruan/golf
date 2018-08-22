package golf

import (
	"strconv"
)

type EventProcessor interface {
	Namer
	ProcessEvent(*Event)
}

type EitherEventProcessor interface {
	EventProcessor
	Judge(e *Event) bool
	// Set down stream processors for Judge() returning true
	Either(EventProcessor) EitherEventProcessor
	// Set down stream processors for Judge() returning false
	Or(EventProcessor) EitherEventProcessor
}

type LogLevelProcessor struct {
	level  Level
	either EventProcessor
	or     EventProcessor
}

func (l LogLevelProcessor) Name() string {
	return "logLevel " + strconv.Itoa(int(l.level))
}

func (l LogLevelProcessor) ProcessEvent(e *Event) {
	if l.Judge(e) {
		if l.either != nil {
			l.either.ProcessEvent(e)
		}
	} else {
		if l.or != nil {
			l.or.ProcessEvent(e)
		}
	}
}

func (l LogLevelProcessor) Judge(e *Event) bool {
	return e.Level >= l.level
}

func (l *LogLevelProcessor) Either(p EventProcessor) EitherEventProcessor {
	l.either = p
	return l
}

func (l *LogLevelProcessor) Or(p EventProcessor) EitherEventProcessor {
	l.or = p
	return l
}

// Create a new LogLevelProcessor with new lvl and all downstream processors
func (l *LogLevelProcessor) WithLevel(lvl Level) *LogLevelProcessor {
	return NewLogLevelProcessor(lvl).Either(l.either).Or(l.or).(*LogLevelProcessor)
}

func NewLogLevelProcessor(lvl Level) *LogLevelProcessor {
	return &LogLevelProcessor{
		level: lvl,
	}
}
