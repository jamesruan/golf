package processor

import (
	"fmt"
	"github.com/jamesruan/golf/event"
	"strings"
)

type logLevelP struct {
	level       event.Level
	ch_started  chan struct{}
	ch_either   chan P
	ch_or       chan P
	ch_event    chan *event.Event
	ch_flushing chan struct{}
	ch_stopping chan struct{}
	ch_stopped  chan struct{}
	either      P
	or          P
}

func (l logLevelP) Name() string {
	var sb strings.Builder
	sb.WriteString("LogLevel")
	sb.WriteString(l.level.String())
	return sb.String()
}

func (l logLevelP) Process(e *event.Event) {
	<-l.ch_started
	select {
	case <-l.ch_stopping:
		return
	default:
		l.ch_event <- e
	}
}

func (l *logLevelP) Either(p P) EitherP {
	l.ch_either <- p
	return l
}

func (l *logLevelP) Or(p P) EitherP {
	l.ch_or <- p
	return l
}

func (l logLevelP) Judge(e *event.Event) bool {
	return e.Level >= l.level
}

func (l logLevelP) process(e *event.Event, flush bool) {
	if l.Judge(e) {
		if l.either != nil {
			l.either.Process(e)
			if flush {
				l.either.Flush()
			}
		}
	} else {
		if l.or != nil {
			l.or.Process(e)
			if flush {
				l.or.Flush()
			}
		}
	}
}

func (t logLevelP) Stopped() <-chan struct{} {
	return t.ch_stopped
}

func (l *logLevelP) Start(stop <-chan struct{}) P {
	select {
	case <-l.ch_started:
		panic(fmt.Sprintf("golf: Start an already started P: %s", l.Name()))
	default:
	}
	parent_stop := stop
	go func() {
		close(l.ch_started)
		flushing := false
		for {
			select {
			case p := <-l.ch_either:
				l.either = p
			case p := <-l.ch_or:
				l.or = p
			case <-l.ch_stopping:
				// flush pending event
				for {
					select {
					case e := <-l.ch_event:
						l.process(e, true)
					default:
						close(l.ch_stopped)
						return
					}
				}
			default:
			}
			select {
			case <-parent_stop:
				l.Stop()
				parent_stop = nil
			case e := <-l.ch_event:
				l.process(e, flushing)
				if len(l.ch_event) == 0 {
					flushing = false
				}
			case <-l.ch_flushing:
				flushing = true
			}
		}
	}()
	return l
}

func (l logLevelP) Stop() {
	select {
	case <-l.ch_stopped:
		panic(fmt.Sprintf("golf: Stop an already Stopped P: levll %s", l.Name()))
	default:
	}
	close(l.ch_stopping)
}

func (l logLevelP) Flush() {
	l.ch_flushing <- struct{}{}
}

// NewLogLevelP returns a processor that process event.Level >= lvl in Either branch.
func NewLogLevelP(lvl event.Level) EitherP {
	return &logLevelP{
		level:       lvl,
		ch_started:  make(chan struct{}),
		ch_event:    make(chan *event.Event),
		ch_either:   make(chan P, 1),
		ch_or:       make(chan P, 1),
		ch_flushing: make(chan struct{}),
		ch_stopping: make(chan struct{}),
		ch_stopped:  make(chan struct{}),
	}
}
