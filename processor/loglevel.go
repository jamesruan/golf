package processor

import (
	"github.com/jamesruan/golf/event"
)

type logLevelP struct {
	level  event.Level
	either P
	or     P
}

func (l logLevelP) Name() string {
	return "LogLevel" + l.level.String()
}

func (l logLevelP) Process(e *event.Event) {
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

func (l logLevelP) Judge(e *event.Event) bool {
	return e.Level >= l.level
}

func (l *logLevelP) Either(p P) EitherP {
	l.either = p
	return l
}

func (l *logLevelP) Or(p P) EitherP {
	l.or = p
	return l
}

// NewLogLevelP returns a processor that process event.Level >= lvl in Either branch.
func NewLogLevelP(lvl event.Level) EitherP {
	return &logLevelP{
		level: lvl,
	}
}
