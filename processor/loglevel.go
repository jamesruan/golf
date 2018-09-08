package processor

import (
	"context"
	"github.com/jamesruan/golf/event"
	"strings"
)

type logLevelP struct {
	level      event.Level
	ch_started chan struct{}
	ch_either  chan P
	ch_or      chan P
	ch_event   chan *event.Event
	either     P
	or         P
	ctx        context.Context
	cancel     context.CancelFunc
}

func (l logLevelP) Name() string {
	var sb strings.Builder
	sb.WriteString("LogLevel")
	sb.WriteString(l.level.String())
	return sb.String()
}

func (l logLevelP) Context() context.Context {
	return l.ctx
}

func (l logLevelP) Process(e *event.Event) {
	<-l.ch_started
	l.ch_event <- e
}

func (l logLevelP) Judge(e *event.Event) bool {
	return e.Level >= l.level
}

func (l *logLevelP) Either(p P) EitherP {
	l.ch_either <- p
	return l
}

func (l *logLevelP) Or(p P) EitherP {
	l.ch_or <- p
	return l
}

func (l *logLevelP) Start(ctx context.Context) P {
	nctx, cancel := context.WithCancel(ctx)
	l.ctx = nctx
	l.cancel = cancel

	go func() {
		close(l.ch_started)
		for {
			select {
			case <-l.ctx.Done():
				return
			case e := <-l.ch_event:
				if l.Judge(e) {
					if l.either != nil {
						l.either.Process(e)
					}
				} else {
					if l.or != nil {
						l.or.Process(e)
					}
				}
			case p := <-l.ch_either:
				l.either = p
			case p := <-l.ch_or:
				l.or = p
			}
		}
	}()
	return l
}

func (l *logLevelP) End() {
	select {
	case <-l.ch_started:
		l.cancel()
	default:
	}
}

// NewLogLevelP returns a processor that process event.Level >= lvl in Either branch.
func NewLogLevelP(lvl event.Level) EitherP {
	return &logLevelP{
		ch_started: make(chan struct{}),
		ch_event:   make(chan *event.Event),
		ch_either:  make(chan P, 1),
		ch_or:      make(chan P, 1),
		level:      lvl,
	}
}
