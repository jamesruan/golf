package processor

import (
	"github.com/jamesruan/golf/event"
)

type LogLevelP struct {
	level  event.Level
	either P
	or     P
}

func (l LogLevelP) Name() string {
	return "LogLevel" + l.level.String()
}

func (l LogLevelP) Process(e *event.Event) {
	if l.Judge(e) {
		if l.either != nil {
			l.either.Process(e)
		}
	} else {
		if l.or != nil {
			l.or.Process(e)
		}
	}
}

func (l LogLevelP) Judge(e *event.Event) bool {
	return e.Level >= l.level
}

func (l *LogLevelP) Either(p P) EitherP {
	l.either = p
	return l
}

func (l *LogLevelP) Or(p P) EitherP {
	l.or = p
	return l
}

func NewLogLevelP(lvl event.Level) *LogLevelP {
	return &LogLevelP{
		level: lvl,
	}
}
